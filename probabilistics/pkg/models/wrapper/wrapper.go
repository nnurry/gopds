package wrapper

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/nnurry/gopds/probabilistics/internal/config"
	"github.com/nnurry/gopds/probabilistics/internal/database/postgres"
	"github.com/nnurry/gopds/probabilistics/internal/service"
	"github.com/nnurry/gopds/probabilistics/pkg/models/decayable"
)

type Wrapper struct {
	filters      *FilterWrapper
	cardinals    *CardinalWrapper
	syncInterval time.Duration
}

var pw *Wrapper
var DecayTicker *time.Ticker
var DecayWg sync.WaitGroup
var DecayDone chan bool

func NewWrapper() *Wrapper {
	return &Wrapper{
		filters:      NewFilterWrapper(),
		cardinals:    NewCardinalWrapper(),
		syncInterval: config.ProbabilisticCfg.SyncInterval,
	}
}

func (pw *Wrapper) FilterWrapper() *FilterWrapper {
	return pw.filters
}

func (pw *Wrapper) CardinalWrapper() *CardinalWrapper {
	return pw.cardinals
}

func (pw *Wrapper) AddFilter(k FilterKey, v *decayable.Filter) {
	pw.filters.Add(k, v)
}

func (pw *Wrapper) AddCardinal(k CardinalKey, v *decayable.Cardinal) {
	pw.cardinals.Add(k, v)
}

func (pw *Wrapper) DeleteFilter(k FilterKey) {
	pw.filters.Delete(k)
}

func (pw *Wrapper) DeleteCardinal(k CardinalKey) {
	pw.cardinals.Delete(k)
}

func GetWrapper() *Wrapper {
	return pw
}

func Synchronize(ticker *time.Ticker, mainWg *sync.WaitGroup, mainDone chan bool) {
	decayedFilterKeys := []FilterKey{}
	decayedCardinalKeys := []CardinalKey{}
	mainWg.Add(1)
	var mu sync.Mutex
	go func() {
		defer func() {
			mainWg.Done()
			fmt.Println("Stopped Synchronize goroutines")
		}()
		subWg := sync.WaitGroup{}
		for {
			select {
			case <-mainDone:
				subWg.Wait()
				fmt.Println("Synchronize() cleanly stopped")
				return
			case <-ticker.C:
				decayedFilterKeys = []FilterKey{}
				decayedCardinalKeys = []CardinalKey{}
				startTime := time.Now().UTC()
				mu.Lock()

				tx, _ := postgres.Client.Begin()
				for key, filter := range pw.filters.Core() {
					subWg.Add(1)
					go func(k FilterKey, v *decayable.Filter) {
						defer subWg.Done()
						if v.Meta().IsDecayed(startTime) {
							decayedFilterKeys = append(decayedFilterKeys, k)
						}
						// Logic to synchronize with database here
						if err := service.SaveFilter(v, false, false, false, tx); err != nil {
							fmt.Println("Got error while syncing filter", key, err)
						}
					}(key, filter)
				}

				for key, cardinal := range pw.cardinals.Core() {
					subWg.Add(1)
					go func(k CardinalKey, v *decayable.Cardinal) {
						defer subWg.Done()
						if v.Meta().IsDecayed(startTime) {
							decayedCardinalKeys = append(decayedCardinalKeys, k)
						}
						// Logic to synchronize with database here
						if err := service.SaveCardinal(v, false, false, false, tx); err != nil {
							fmt.Println("Got error while syncing cardinal", key, err)
						}
					}(key, cardinal)
				}
				subWg.Wait()
				if len(decayedFilterKeys) > 0 {
					for _, key := range decayedFilterKeys {
						fmt.Println("Decayed filter:", key)
						pw.DeleteFilter(key)
					}
				}

				if len(decayedCardinalKeys) > 0 {
					for _, key := range decayedCardinalKeys {
						fmt.Println("Decayed cardinal:", key)
						pw.DeleteCardinal(key)
					}
				}

				tx.Commit()

				mu.Unlock()
			}
		}
	}()

}

func Cleanup(osChan chan os.Signal, srv *http.Server) {
	defer DecayWg.Done()

	sig := <-osChan
	fmt.Println("Encountered signal:", sig.String())
	DecayTicker.Stop()
	DecayDone <- true

	time.Sleep(500 * time.Millisecond)
	close(osChan)

	if err := srv.Shutdown(context.Background()); err != nil {
		fmt.Printf("HTTP server Shutdown: %v", srv.Shutdown(context.Background()))
	}
	close(DecayDone)
}

func init() {
	pw = NewWrapper()
	DecayTicker = time.NewTicker(pw.syncInterval)
	DecayWg = sync.WaitGroup{}
	DecayDone = make(chan bool, 1)
	Synchronize(DecayTicker, &DecayWg, DecayDone)
}
