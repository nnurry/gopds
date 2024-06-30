package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/nnurry/gopds/probabilistics/internal/api"
	"github.com/nnurry/gopds/probabilistics/internal/database/postgres"
	"github.com/nnurry/gopds/probabilistics/pkg/models/wrapper"
)

func main() {
	postgres.Bootstrap()
	osChan := make(chan os.Signal, 1)
	signal.Notify(osChan, syscall.SIGTERM, syscall.SIGINT)

	wrapper.DecayWg.Add(1)

	mux := api.SetupMux()
	srv := http.Server{
		Addr:    ":5000",
		Handler: mux,
	}

	go wrapper.Cleanup(osChan, &srv)
	err := srv.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		fmt.Println("Can't start server:", err)
		osChan <- syscall.SIGTERM
	}
	wrapper.DecayWg.Wait()
}
