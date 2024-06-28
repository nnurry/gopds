package wrapper

import (
	"context"
	"fmt"
	"gopds/probabilistics/internal/config"
	"gopds/probabilistics/internal/database/postgres"
	"gopds/probabilistics/internal/service"
	"gopds/probabilistics/pkg/models/probabilistic"
	"net/http"
	"os"
	"sync"
	"time"
)

type WrapperKey service.ProbCreateBody

type Wrapper struct {
	probmap      map[WrapperKey]*probabilistic.Probabilistic
	counter      uint
	syncInterval time.Duration
}

var pw *Wrapper
var DecayTicker *time.Ticker
var DecayWg sync.WaitGroup
var DecayDone chan bool

func NewWrapper() *Wrapper {
	return &Wrapper{
		probmap:      make(map[WrapperKey]*probabilistic.Probabilistic),
		counter:      0,
		syncInterval: config.ProbabilisticCfg.SyncInterval,
	}
}

func (pw *Wrapper) Add(k WrapperKey, p *probabilistic.Probabilistic) {
	_, exists := pw.probmap[k]
	pw.probmap[k] = p
	if !exists {
		pw.counter++
	}
}

func (pw *Wrapper) Delete(k WrapperKey) {
	_, exists := pw.probmap[k]
	delete(pw.probmap, k)
	if exists {
		pw.counter--
	}
}

func GetWrapper() *Wrapper {
	return pw
}

func Synchronize(ticker *time.Ticker, mainWg *sync.WaitGroup, mainDone chan bool) {
	decayedProbKeys := []WrapperKey{}
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
				decayedProbKeys = []WrapperKey{}
				startTime := time.Now().UTC()
				mu.Lock()
				tx, _ := postgres.Client.Begin()
				for key, prob := range pw.probmap {
					subWg.Add(1)
					go func(k WrapperKey, p *probabilistic.Probabilistic) {
						defer subWg.Done()
						if p.Meta().IsDecayed(startTime) {
							decayedProbKeys = append(decayedProbKeys, k)
						}
						// Logic to synchronize with database here
						service.SaveProbabilistic(p, false, false, tx)
					}(key, prob)
				}
				subWg.Wait()
				if len(decayedProbKeys) > 0 {
					for _, key := range decayedProbKeys {
						fmt.Println("Decayed", key)
						pw.Delete(key)
					}
				}
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
