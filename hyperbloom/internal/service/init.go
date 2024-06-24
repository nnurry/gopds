package service

import (
	"fmt"
	"log"
	"time"

	"github.com/nnurry/gopds/hyperbloom/config"
	"github.com/nnurry/gopds/hyperbloom/internal/database/postgres"
)

func init() {
	// Initialize database client
	client := postgres.DbClient

	// Begin a new transaction
	tx, _ := client.Begin()

	// Execute SQL query to create 'bloom_filters' table if not exists
	_, err := client.Exec(`
	CREATE TABLE IF NOT EXISTS bloom_filters (
		key VARCHAR PRIMARY KEY, 
		bloombyte BYTEA, 
		hyperbyte BYTEA
	)`)

	// Rollback transaction and log fatal error if table creation fails
	if err != nil {
		tx.Rollback()
		log.Fatal("Can't create table bloom_filters", err)
	}

	// Create a new ticker that ticks at the specified interval in milliseconds
	ticker := time.NewTicker(config.HyperBloomConfig.UpdateRate)

	// Start asynchronous process to update bloom filters using the ticker
	AsyncBloomUpdate(ticker, StopAsyncBloomUpdate)

	// Print a message indicating successful initialization
	fmt.Println("Init service")
}
