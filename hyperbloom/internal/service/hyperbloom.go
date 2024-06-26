package service

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"gopds/hyperbloom/internal/config"
	"gopds/hyperbloom/internal/database/postgres"
	"gopds/hyperbloom/pkg/models"

	"github.com/bits-and-blooms/bloom/v3"
)

// AsyncBloomUpdate starts a goroutine that periodically updates all HyperBloom instances in memory
// at the specified interval (in milliseconds). The updates are performed asynchronously.
func AsyncBloomUpdate(ticker *time.Ticker, done chan bool) {
	fmt.Println("AsyncBloomUpdate") // Print a message indicating the function has started
	mutex := &sync.Mutex{}          // Initialize a new mutex for thread-safe operations
	WG.Add(1)
	// Start a new goroutine to handle the periodic updates
	go func() {
		defer WG.Done()
		for {
			// Loop indefinitely, executing at each tick of the ticker or done signal
			select {
			case <-done:
				fmt.Println("Received signal to stop AsyncBloomUpdate")
				return // Exit the goroutine when done signal is received

			case <-ticker.C:
				// Execute when the ticker ticks

				keysToPrune := []string{} // Initialize an empty slice to store keys that need pruning

				// Lock the mutex for writing to ensure exclusive access to the dbs resource
				mutex.Lock()

				tx, _ := postgres.DbClient.Begin() // Begin a database transaction

				// Iterate over all HyperBloom instances and update each one
				currentTime := time.Now().UTC() // Get the current time in UTC
				for _, db := range dbs.GetInMemoryHyperBlooms() {
					fmt.Println("Sync Hyperbloom object with database", db.Key()) // Print a synchronization message
					BloomUpdate(db, false)                                        // Update the HyperBloom instance

					// Check if the HyperBloom instance has decayed
					if db.CheckDecayed(currentTime) {
						keysToPrune = append(keysToPrune, db.Key()) // Add the key to prune list if decayed
					}
				}

				tx.Commit() // Commit the database transaction

				// Remove decayed HyperBloom instances from memory
				for _, key := range keysToPrune {
					fmt.Println("Decay", key, "from in-memory Hyperblooms") // Print a message for each decayed instance
					dbs.Remove(key)                                         // Remove the decayed instance from memory
				}

				// Unlock the mutex after updates are done
				mutex.Unlock()
			}
		}
	}()
}

// BloomList retrieves all HyperBloom instances stored in memory
func BloomList() []*models.HyperBloom {
	return dbs.GetHyperBlooms()
}

// BloomGet retrieves a HyperBloom instance by key.
// It first attempts to get the HyperBloom from memory,
// and if not found, it fetches it from the database.
func BloomGet(key string) *models.HyperBloom {
	// Attempt to get the HyperBloom from memory or fetch it from the database
	db, err := dbs.GetOrFetchHyperBloom(key)

	// If an error occurred (e.g., the HyperBloom couldn't be fetched from the database), return nil
	if err != nil {
		return nil
	}

	// Refresh the HyperBloom instance to update its last used timestamp or any other necessary fields
	db.Refresh()

	// Return the retrieved HyperBloom instance
	return db
}

// BloomHash adds a value to the Bloom filter and HyperLogLog sketch of the HyperBloom identified by key.
// If the HyperBloom does not exist, it creates a new one.
func BloomHash(key, value string) {
	var err error
	var db *models.HyperBloom

	// Attempt to fetch or retrieve the HyperBloom for the given key
	db, err = dbs.GetOrFetchHyperBloom(key)

	// If there's an error (HyperBloom not found in memory or database)
	if err != nil {
		// Create a new HyperBloom instance using default configuration
		db = BloomCreate(
			config.HyperBloomCfg.Cardinality,
			config.HyperBloomCfg.FalsePositive,
			key,
		)
	}

	// Hash the value using Bloom filter and HyperLogLog
	db.Hash(value)

	// Add the HyperBloom instance into memory (or update if already exists)
	dbs.Set(db, key)
}

// BloomUpdate synchronizes the HyperBloom instance in memory with the database.
func BloomUpdate(db *models.HyperBloom, doCommit bool) {
	// Define the SQL query to insert or update the bloom_filters table
	query := `
		INSERT INTO hyperblooms (key, bloombyte, hyperbyte)
		VALUES ($1, $2, $3)
		ON CONFLICT (key) DO UPDATE
		SET bloombyte = EXCLUDED.bloombyte,
			hyperbyte = EXCLUDED.hyperbyte;
	`

	// Initialize a transaction if doCommit is true
	var tx *sql.Tx
	if doCommit {
		tx, _ = postgres.DbClient.Begin() // Begin a transaction
	}

	// Encode the Bloom filter and HyperLogLog data into byte representations
	bloomByterepr, _ := db.Bloom().GobEncode()     // Encode Bloom filter data
	hyperByterepr, _ := db.Hyper().MarshalBinary() // Marshal HyperLogLog data

	// Execute the SQL query to insert or update the record
	postgres.DbClient.Exec(query, db.Key(), bloomByterepr, hyperByterepr)

	// Commit the transaction if doCommit is true
	if doCommit {
		tx.Commit() // Commit the transaction
	}
}

// BloomDecay removes a HyperBloom instance from memory if it has decayed (i.e., last used timestamp exceeds decay duration).
func BloomDecay(key string) {
	// Check if the HyperBloom instance with the given key has decayed based on the current UTC time
	if dbs.CheckDecayed(key, time.Now().UTC()) {
		// If decayed, remove the HyperBloom instance from memory
		dbs.Remove(key)
	}
}

// BloomExists checks if a value exists in the Bloom filter of the HyperBloom identified by key.
func BloomExists(key, value string) bool {
	db := BloomGet(key)
	if db != nil {
		return db.CheckExists(value)
	}
	return false
}

// AllBoolList checks if all elements in boolList are equal.
func AllBoolList(boolList []bool) bool {
	// Iterate through the boolList slice
	for i := 0; i < len(boolList); i++ {
		// If any element is false, return false immediately
		if !boolList[i] {
			return false
		}
	}
	// If all elements are true, return true
	return true
}

// AnyBoolList checks if any element in boolList is equal to true.
func AnyBoolList(boolList []bool) bool {
	// Iterate through the boolList slice
	for i := 0; i < len(boolList); i++ {
		// If any element is true, return true immediately
		if boolList[i] {
			return true
		}
	}
	// If no element is true, return false
	return false
}

// BloomChainingExists checks existence of a value in Bloom filters associated with given keys.
func BloomChainingExists(keys []string, value string, operator string) bool {
	// Initialize an empty boolean slice to store results for each key
	boolList := []bool{}

	// Iterate through each key
	for _, key := range keys {
		// Retrieve Bloom filter for the key
		db := BloomGet(key)
		_bool := false

		// If Bloom filter exists for the key, check if value exists in it
		if db != nil {
			_bool = db.CheckExists(value)
		}

		// Append the result (true/false) to boolList
		boolList = append(boolList, _bool)
	}

	// Determine the final result based on the specified operator
	if operator == "AND" {
		// Return true if all elements in boolList are true
		return AllBoolList(boolList)
	} else if operator == "OR" {
		// Return true if any element in boolList is true
		return AnyBoolList(boolList)
	} else {
		// Default case: return false if operator is neither "AND" nor "OR"
		return false
	}
}

// BloomBitwiseExists checks the existence of a value in Bloom filters associated with given keys using bitwise operations.
func BloomBitwiseExists(keys []string, value string, operator string) bool {
	// Check if there are keys provided
	if len(keys) < 1 {
		return false
	}

	// Get the Bloom filter for the first key
	db := BloomGet(keys[0])
	if db == nil {
		return false
	}

	// Get the BitSet of the Bloom filter for the first key
	bs := db.Bloom().BitSet()

	// Iterate through the rest of the keys
	for i := 1; i < len(keys); i++ {
		// Get the Bloom filter for the current key
		db = BloomGet(keys[i])
		switch operator {
		case "AND":
			// If the Bloom filter does not exist, return false for AND operation
			if db == nil {
				return false
			}
			// Perform bitwise AND operation with the BitSet of the current Bloom filter
			bs = bs.Intersection(db.BitSet())
		case "OR":
			// If the Bloom filter does not exist, skip to the next key for OR operation
			if db == nil {
				continue
			}
			// Perform bitwise OR operation with the BitSet of the current Bloom filter
			bs = bs.Union(db.BitSet())
		}
	}

	// Create a new Bloom filter using the resulting BitSet
	b := bloom.FromWithM(
		bs.Bytes(),
		db.Bloom().Cap(),
		db.Bloom().K(),
	)

	// Test if the value exists in the new Bloom filter
	return b.TestString(value)
}

// BloomCardinality returns the cardinality of the Bloom filter and HyperLogLog sketch of the HyperBloom identified by key.
func BloomCardinality(key string) (uint32, uint64) {
	db := BloomGet(key)
	if db != nil {
		bCard := db.BloomCardinality()
		hCard := db.HyperCardinality()

		return bCard, hCard
	}
	return 0, 0
}

// BloomSimilarity calculates the Jaccard similarity between two Bloom filters identified by key1 and key2.
// It returns a float32 value representing the similarity score.
func BloomSimilarity(key1, key2 string) float32 {
	// Retrieve Bloom filter for key1
	db1 := BloomGet(key1)
	// If Bloom filter for key1 is not found, return similarity score of 0.0
	if db1 == nil {
		return 0.0
	}

	// Retrieve Bloom filter for key2
	db2 := BloomGet(key2)
	// If Bloom filter for key2 is not found, return similarity score of 0.0
	if db2 == nil {
		return 0.0
	}

	// Calculate Jaccard similarity between db1 and db2 using models.JaccardSimBF function
	return models.JaccardSimBF(db1, db2)
}

// BloomCreate creates a new HyperBloom instance with specified parameters and stores it in the database.
func BloomCreate(capacity uint, falsePositive float64, key string) *models.HyperBloom {
	// Create a new HyperBloom instance using provided parameters
	db := models.NewHyperBloomFromParams(capacity, falsePositive, key)

	// Serialize the Bloom filter and HyperLogLog data structures to bytes
	bloomByterepr, _ := db.Bloom().GobEncode()     // Serialize Bloom filter
	hyperByterepr, _ := db.Hyper().MarshalBinary() // Serialize HyperLogLog

	// Begin a database transaction
	tx, _ := postgres.DbClient.Begin()

	// Insert the serialized data into the hyperblooms table
	postgres.DbClient.Exec(
		`INSERT INTO hyperblooms (
			key, 
			bloombyte, 
			hyperbyte
		) 
		VALUES ($1, $2, $3)`,
		key,
		bloomByterepr,
		hyperByterepr,
	)

	// Insert metadata about the HyperBloom instance into the hyperblooms_metadata table
	postgres.DbClient.Exec(
		`INSERT INTO hyperblooms_metadata (
			key, 
			max_cardinality, 
			false_positive, 
			bit_capacity, 
			no_hash_func, 
			decay_sec
		) 
		VALUES ($1, $2, $3, $4, $5, $6)`,
		key,
		capacity,
		falsePositive,
		db.Bloom().Cap(),
		db.Bloom().K(),
		db.Decay(),
	)

	// Commit the database transaction
	tx.Commit()

	// Return the created HyperBloom instance
	return db
}
