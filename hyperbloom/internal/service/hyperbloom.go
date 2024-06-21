package service

import "github.com/nnurry/gopds/hyperbloom/pkg/models"

var dbs = models.NewDecayBlooms()

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
		db = models.NewDefaultDecayBloom(key)
	}

	db.Hash(value)
}

// After 2s, spawn a goroutine to sync data in memory with database
func BloomUpdate() {
}

// If time.Now() - time.LastUpdate() >= decayDuration
// Remove it from memory
func BloomDecay(key string) {
	if dbs.CheckDecayed(key) {
		dbs.Remove(key)
	}
}
