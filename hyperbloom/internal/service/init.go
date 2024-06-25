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
	var err error
	client := postgres.DbClient

	// Begin a new transaction
	tx, _ := client.Begin()

	// Execute SQL query to create 'hyperblooms' table if it does not exist
	_, err = client.Exec(`
	CREATE TABLE IF NOT EXISTS hyperblooms (
		key VARCHAR PRIMARY KEY,
		bloombyte BYTEA, 
		hyperbyte BYTEA
	)`)

	// Rollback transaction and log fatal error if table creation fails
	if err != nil {
		tx.Rollback()
		log.Fatal("Can't create table hyperblooms", err)
	}

	// Execute SQL query to create 'hyperbloom_metadata' table if it does not exist
	_, err = client.Exec(`
	CREATE TABLE IF NOT EXISTS hyperbloom_metadata (
		key VARCHAR,
		max_cardinality INTEGER,
		false_positive REAL,
		bit_capacity INTEGER,
		no_hash_func INTEGER,
		decay_sec BIGINT,
		FOREIGN KEY (key) REFERENCES hyperblooms(key)
	)`)

	// Rollback transaction and log fatal error if table creation fails
	if err != nil {
		tx.Rollback()
		log.Fatal("Can't create table hyperbloom_metadata", err)
	}

	// Commit the transaction after successful table creations
	tx.Commit()

	// Create a new ticker that ticks at the specified interval in milliseconds
	ticker := time.NewTicker(config.HyperBloomConfig.UpdateRate)

	// Start asynchronous process to update bloom filters using the ticker
	AsyncBloomUpdate(ticker, StopAsyncBloomUpdate)

	// Print a message indicating successful initialization
	fmt.Println("Init service")
}
