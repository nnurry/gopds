// Package models defines data structures and methods for managing multiple HyperBloom instances,
// including fetching from database, caching in memory, and providing access methods.
package models

import (
	"time"
)

// HyperBlooms manages a collection of HyperBloom instances.
type HyperBlooms struct {
	blooms map[string]*HyperBloom // Map to store HyperBloom instances by key
}

// NewHyperBlooms creates a new HyperBlooms instance initialized with an empty map.
func NewHyperBlooms() *HyperBlooms {
	return &HyperBlooms{
		blooms: make(map[string]*HyperBloom),
	}
}

// GETTERS

// GetInMemoryHyperBlooms retrieves all HyperBloom instances that are currently in memory.
func (dbs *HyperBlooms) GetInMemoryHyperBlooms() []*HyperBloom {
	output := []*HyperBloom{} // Initialize an empty slice to store output

	// Iterate over each key in the 'blooms' map
	for key := range dbs.blooms {
		// Attempt to fetch the HyperBloom for the current 'key'
		db, ok := dbs.GetHyperBloom(key)

		// If fetching the HyperBloom was successful (no error)
		if ok {
			output = append(output, db) // Append the fetched HyperBloom to the output slice
		}
	}

	return output // Return the slice containing all HyperBloom instances in memory
}

// GetHyperBlooms retrieves all HyperBloom instances from the HyperBlooms collection,
// fetching them if necessary.
func (dbs *HyperBlooms) GetHyperBlooms() []*HyperBloom {
	output := []*HyperBloom{} // Initialize an empty slice to store output

	// Iterate over each key in the 'blooms' map
	for key := range dbs.blooms {
		// Attempt to fetch or retrieve the HyperBloom for the current 'key'
		db, err := dbs.GetOrFetchHyperBloom(key)

		// If fetching the HyperBloom was successful (no error)
		if err == nil {
			output = append(output, db) // Append the fetched HyperBloom to the output slice
		}
	}

	return output // Return the slice containing all fetched HyperBloom instances
}

// GetHyperBloomKeys retrieves keys of HyperBloom instances from the 'blooms' map,
// attempting to fetch each HyperBloom and returning keys of successfully fetched ones.
func (dbs *HyperBlooms) GetHyperBloomKeys() []string {
	// Initialize an empty slice to store output
	output := []string{}

	// Iterate over each key in the 'blooms' map
	for key := range dbs.blooms {
		// Attempt to fetch the hyperbloom for the current 'key'
		db, err := dbs.GetOrFetchHyperBloom(key)

		// If fetching the hyperbloom was successful (no error)
		if err == nil {
			// Set the fetched hyperbloom in the 'blooms' map
			dbs.Set(db, key)
			// Append the key of the fetched hyperbloom to the output slice
			output = append(output, db.Key())
		}
	}

	// Return the slice containing keys of successfully fetched hyperblooms
	return output
}

// GetInMemoryHyperBloomKeys retrieves keys of HyperBloom instances from the 'blooms' map,
// assuming they are already in memory, and returning keys of successfully fetched ones.
func (dbs *HyperBlooms) GetInMemoryHyperBloomKeys() []string {
	// Initialize an empty slice to store output
	output := []string{}

	// Iterate over each key in the 'blooms' map
	for key := range dbs.blooms {
		// Attempt to fetch the hyperbloom for the current 'key'
		db, ok := dbs.GetHyperBloom(key)

		// If fetching the hyperbloom was successful (found in memory)
		if ok {
			// Set the fetched hyperbloom in the 'blooms' map
			dbs.Set(db, key)
			// Append the key of the fetched hyperbloom to the output slice
			output = append(output, db.Key())
		}
	}

	// Return the slice containing keys of successfully fetched hyperblooms
	return output
}

// GetOrFetchHyperBloom retrieves a HyperBloom instance from the collection or fetches it from the database.
func (dbs *HyperBlooms) GetOrFetchHyperBloom(key string) (*HyperBloom, error) {
	var db *HyperBloom
	var ok bool
	var err error

	db, ok = dbs.GetHyperBloom(key)
	if !ok {
		db, err = GetBloomFromDB(key)
	}
	if err != nil {
		return nil, err
	}
	dbs.Set(db, key)
	return db, nil
}

// GetHyperBloom retrieves a HyperBloom instance by key from the collection.
func (dbs *HyperBlooms) GetHyperBloom(key string) (*HyperBloom, bool) {
	db, ok := dbs.blooms[key]
	return db, ok
}

// SETTERS

// Remove deletes a HyperBloom instance from the HyperBlooms collection by key.
func (dbs *HyperBlooms) Remove(key string) {
	delete(dbs.blooms, key)
}

// Set adds a HyperBloom instance to the HyperBlooms collection.
func (dbs *HyperBlooms) Set(db *HyperBloom, key string) {
	// Refresh the last used timestamp of the HyperBloom instance
	db.Refresh()

	// Add the HyperBloom instance to the 'blooms' map in HyperBlooms
	dbs.blooms[key] = db
}

// CheckDecayed checks if a HyperBloom instance has decayed based on the last used timestamp.
func (dbs *HyperBlooms) CheckDecayed(key string, timemark time.Time) bool {
	// Retrieve the HyperBloom instance for the given key
	db, ok := dbs.GetHyperBloom(key)

	// If the HyperBloom instance was found, proceed to check decay
	return ok && db.CheckDecayed(timemark)
}
