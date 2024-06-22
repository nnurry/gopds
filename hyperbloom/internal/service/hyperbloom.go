package service

import (
	"github.com/nnurry/gopds/hyperbloom/internal/database/postgres"
	"github.com/nnurry/gopds/hyperbloom/pkg/models"
)

var dbs = models.NewDecayBlooms()

// Get each bloom stored in memory
func BloomList() []*models.DecayBloom {
	return dbs.GetBlooms()
}

// Get each bloom stored in memory
func BloomGet(key string) *models.DecayBloom {
	db, ok := dbs.GetBloom(key)
	if ok {
		return db
	}
	return nil
}

// Check in memory for the bloom filter
// Not exist -> check in database
// Not exist -> create a new one
func BloomHash(key, value string) {
	var err error
	var db *models.DecayBloom

	// Check in memory
	db, ok := dbs.GetBloom(key)
	if !ok {
		// Can't find in memory
		// Fetch from database
		db, err = models.GetBloomFromDB(key)
	}

	if err != nil {
		// Can't find in database
		// Create new DecayBloom
		db = BloomCreate(key)
	}
	// Hash the value
	db.Hash(value)

	// Add into memory
	dbs.Set(db, key)

	// Update database
	BloomUpdate(db)
}

// Sync data in memory with database (prototype)
func BloomUpdate(db *models.DecayBloom) {
	query := "UPDATE bloom_filters SET bloombyte = $2 WHERE key = $1;"
	tx, _ := postgres.DbClient.Begin()
	byterepr, _ := db.Bloom().GobEncode()
	postgres.DbClient.Exec(query, db.Key(), byterepr)
	tx.Commit()
}

// If time.Now() - time.LastUpdate() >= decayDuration
// Remove it from memory
func BloomDecay(key string) {
	if dbs.CheckDecayed(key) {
		dbs.Remove(key)
	}
}

// Create new bloom filter in database and return that object
func BloomCreate(key string) *models.DecayBloom {
	db := models.NewDefaultDecayBloom(key)
	query := "INSERT INTO bloom_filters (key, bloombyte) VALUES ($1, $2)"
	// NOTE: Haven't handled the error for serializing BloomFilter object
	byterepr, _ := db.Bloom().GobEncode()
	// NOTE: Haven't handled the error for inserting into bloom_filters
	tx, _ := postgres.DbClient.Begin()
	postgres.DbClient.Exec(query, key, byterepr)
	tx.Commit()
	return db
}
