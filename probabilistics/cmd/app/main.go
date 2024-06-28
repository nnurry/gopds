package main

import (
	"gopds/probabilistics/pkg/models/probabilistic"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	osChan := make(chan os.Signal, 1)
	signal.Notify(osChan, syscall.SIGTERM, syscall.SIGINT)

	go probabilistic.Cleanup(osChan)

	probabilistic.DecayWg.Wait()
}
