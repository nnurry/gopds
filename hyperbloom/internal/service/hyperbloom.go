package service

import (
	"github.com/nnurry/gopds/hyperbloom/internal/database/postgres"
	"github.com/nnurry/gopds/hyperbloom/pkg/models"
)

var dbs = models.NewHyperBlooms()

// Get each bloom stored in memory
func BloomList() []*models.HyperBloom {
	return dbs.GetHyperBlooms()
}

// Get each bloom stored in memory
func BloomGet(key string) *models.HyperBloom {
	db, ok := dbs.GetHyperBloom(key)
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
	var db *models.HyperBloom

	// Check in memory
	db, ok := dbs.GetHyperBloom(key)
	if !ok {
		// Can't find in memory
		// Fetch from database
		db, err = models.GetBloomFromDB(key)
	}

	if err != nil {
		// Can't find in database
		// Create new HyperBloom
		db = BloomCreate(key)
	}
	// Hash the value
	db.Hash(value)

	// Add into memory
	dbs.Set(db, key)

	// Update database (haven't decoupled updating from hashing)
	BloomUpdate(db)
}

// Sync data in memory with database.
// Will implement queue to periodically update database
func BloomUpdate(db *models.HyperBloom) {
	query := "UPDATE bloom_filters SET bloombyte = $2, hyperbyte = $3 WHERE key = $1;"
	tx, _ := postgres.DbClient.Begin()
	bloomByterepr, _ := db.Bloom().GobEncode()
	hyperByterepr, _ := db.Hyper().MarshalBinary()
	postgres.DbClient.Exec(query, db.Key(), bloomByterepr, hyperByterepr)
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
