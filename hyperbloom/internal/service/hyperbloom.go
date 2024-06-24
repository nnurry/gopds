package service

import (
	"github.com/nnurry/gopds/hyperbloom/internal/database/postgres"
	"github.com/nnurry/gopds/hyperbloom/pkg/models"
)

var dbs = models.NewHyperBlooms()

// BloomList retrieves all HyperBloom instances stored in memory
func BloomList() []*models.HyperBloom {
	return dbs.GetHyperBlooms()
}

// BloomGet retrieves a HyperBloom instance by key, fetching from the database if not found in memory
func BloomGet(key string) *models.HyperBloom {
	db, err := dbs.GetOrFetchHyperBloom(key)
	if err != nil {
		return nil
	}
	return db
}

// BloomHash adds a value to the Bloom filter and HyperLogLog sketch of the HyperBloom identified by key.
// If the HyperBloom does not exist, it creates a new one.
func BloomHash(key, value string) {
	var err error
	var db *models.HyperBloom

	db, err = dbs.GetOrFetchHyperBloom(key)

	if err != nil {
		// Can't find in both memory and database, create a new HyperBloom
		db = BloomCreate(key)
	}
	// Hash the value
	db.Hash(value)

	// Add into memory
	dbs.Set(db, key)

	// Update database (haven't decoupled updating from hashing)
	BloomUpdate(db)
}

// BloomUpdate synchronizes the HyperBloom instance in memory with the database.
func BloomUpdate(db *models.HyperBloom) {
	query := "UPDATE bloom_filters SET bloombyte = $2, hyperbyte = $3 WHERE key = $1;"
	tx, _ := postgres.DbClient.Begin()
	bloomByterepr, _ := db.Bloom().GobEncode()
	hyperByterepr, _ := db.Hyper().MarshalBinary()
	postgres.DbClient.Exec(query, db.Key(), bloomByterepr, hyperByterepr)
	tx.Commit()
}

// BloomDecay removes a HyperBloom instance from memory if it has decayed (i.e., last used timestamp exceeds decay duration).
func BloomDecay(key string) {
	if dbs.CheckDecayed(key) {
		dbs.Remove(key)
	}
}

// BloomExists checks if a value exists in the Bloom filter of the HyperBloom identified by key.
func BloomExists(key, value string) bool {
	db := BloomGet(key)
	return db.CheckExists(value)
}

// BloomCardinality returns the cardinality of the Bloom filter and HyperLogLog sketch of the HyperBloom identified by key.
func BloomCardinality(key string) (uint32, uint64) {
	db, _ := dbs.GetHyperBloom(key)
	bCard := db.BloomCardinality()
	hCard := db.HyperCardinality()

	return bCard, hCard
}

// BloomCreate creates a new HyperBloom instance and stores it in the database.
func BloomCreate(key string) *models.HyperBloom {
	db := models.NewDefaultHyperBloom(key)
	query := "INSERT INTO bloom_filters (key, bloombyte, hyperbyte) VALUES ($1, $2, $3)"
	// NOTE: Haven't handled the error for serializing HyperBloom object
	bloomByterepr, _ := db.Bloom().GobEncode()
	hyperByterepr, _ := db.Hyper().MarshalBinary()
	// NOTE: Haven't handled the error for inserting into bloom_filters
	tx, _ := postgres.DbClient.Begin()
	postgres.DbClient.Exec(query, key, bloomByterepr, hyperByterepr)
	tx.Commit()
	return db
}
