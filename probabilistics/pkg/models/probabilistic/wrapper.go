package probabilistic

import (
	"fmt"
	"gopds/probabilistics/internal/config"
	"os"
	"sync"
	"time"
)

type WrapperKey string
type Wrapper struct {
	probmap      map[WrapperKey]*Probabilistic
	counter      uint
	syncInterval time.Duration
}

var pw *Wrapper
var DecayTicker *time.Ticker
var DecayWg sync.WaitGroup
var DecayDone chan bool

func NewWrapper() *Wrapper {
	return &Wrapper{
		probmap:      make(map[WrapperKey]*Probabilistic),
		counter:      0,
		syncInterval: config.ProbabilisticCfg.SyncInterval,
	}
}

func (pw *Wrapper) Add(k WrapperKey, p *Probabilistic) {
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
		defer mainWg.Done()
		subWg := &sync.WaitGroup{}
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
				for key, prob := range pw.probmap {
					subWg.Add(1)
					go func(k WrapperKey, p *Probabilistic) {
						defer subWg.Done()
						// Logic to synchronize with database here
						if p.Meta().IsDecayed(startTime) {
							decayedProbKeys = append(decayedProbKeys, k)
						}
					}(key, prob)
				}
				subWg.Wait()
				if len(decayedProbKeys) > 0 {
					for _, key := range decayedProbKeys {
						pw.Delete(key)
					}
				}
				mu.Unlock()
			}
		}
	}()

}

func Cleanup(osChan chan os.Signal) {
	defer DecayWg.Done()

	sig := <-osChan
	fmt.Println("Encountered signal:", sig.String())
	DecayTicker.Stop()
	DecayDone <- true

	time.Sleep(500 * time.Millisecond)
	close(DecayDone)
	close(osChan)
}

func init() {
	pw = NewWrapper()
	DecayTicker = time.NewTicker(pw.syncInterval)
	DecayWg = sync.WaitGroup{}
	DecayDone = make(chan bool, 1)
	Synchronize(DecayTicker, &DecayWg, DecayDone)
}
