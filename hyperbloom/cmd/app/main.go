package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"gopds/hyperbloom/internal/api"
	"gopds/hyperbloom/internal/config"
	"gopds/hyperbloom/internal/service"
	"gopds/hyperbloom/pkg/utils"
)

// main sets up the HTTP server and routes, including handling OS interrupts for graceful shutdown.
func main() {
	var err error

	// Load application configuration from environment variables or configuration files
	config.LoadConfigApplication()

	// Create a new ServeMux instance to handle HTTP requests
	mux := http.NewServeMux()

	// Set up a channel to receive OS interrupt signals
	osChan := make(chan os.Signal, 1)
	signal.Notify(osChan, syscall.SIGTERM, syscall.SIGINT)

	// Goroutine to handle OS interrupt signals and perform cleanup tasks
	service.WG.Add(1)
	go utils.Cleanup(osChan, &service.WG)

	// Register HTTP request handlers for specific API endpoints
	api.Serve(mux)

	// Start the HTTP server on port 5000
	err = http.ListenAndServe(config.ApplicationCfg.Addr, mux)
	if err != nil {
		log.Println("Can't start server:", err) // Log error if the server fails to start
		osChan <- syscall.SIGTERM               // Signal to initiate graceful shutdown
	}

	service.WG.Wait() // Wait for all cleanup tasks to finish before exiting
}
