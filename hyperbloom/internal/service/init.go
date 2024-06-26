package service

import (
	"fmt"
	"log"
	"sync"
	"time"

	"gopds/hyperbloom/internal/config"
	"gopds/hyperbloom/internal/database/postgres"
	"gopds/hyperbloom/pkg/models"
)

// dbs is a global instance of models.HyperBlooms, representing the collection of Bloom filters.
var dbs = models.NewHyperBlooms()

// StopAsyncBloomUpdate is a channel used to signal stopping the asynchronous Bloom filter update process.
var StopAsyncBloomUpdate = make(chan bool, 1)

// WG is a WaitGroup used to synchronize concurrent operations.
var WG sync.WaitGroup

func init() {
	// Initialize database client
	var err error
	config.LoadConfigHyperBloom()
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
		log.Fatal("Can't create table hyperblooms", err)
		tx.Rollback()
	}

	// Execute SQL query to create 'hyperblooms_metadata' table if it does not exist
	_, err = client.Exec(`
	CREATE TABLE IF NOT EXISTS hyperblooms_metadata (
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
		log.Fatal("Can't create table hyperblooms_metadata", err)
		tx.Rollback()
	}

	// Commit the transaction after successful table creations
	tx.Commit()

	// Create a new ticker that ticks at the specified interval in milliseconds
	ticker := time.NewTicker(config.HyperBloomCfg.UpdateRate)

	// Start asynchronous process to update bloom filters using the ticker
	AsyncBloomUpdate(ticker, StopAsyncBloomUpdate)

	// Print a message indicating successful initialization
	fmt.Println("Init service")
}
