package models

import (
	"time"

	"github.com/axiomhq/hyperloglog"
	"github.com/bits-and-blooms/bitset"
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/nnurry/gopds/hyperbloom/config"
	"github.com/nnurry/gopds/hyperbloom/internal/database/postgres"
)

// HyperBlooms manages multiple HyperBloom instances
type HyperBlooms struct {
	blooms map[string]*HyperBloom
}

// HyperBloom represents a Bloom filter combined with a HyperLogLog sketch
type HyperBloom struct {
	bloom    *bloom.BloomFilter
	hyper    *hyperloglog.Sketch
	key      string
	decay    time.Duration
	lastUsed time.Time
}

// NewHyperBlooms creates a new HyperBlooms instance
func NewHyperBlooms() *HyperBlooms {
	return &HyperBlooms{
		blooms: make(map[string]*HyperBloom),
	}
}

// NewHyperBloom creates a new HyperBloom instance with given Bloom filter and HyperLogLog sketch
func NewHyperBloom(bf *bloom.BloomFilter, hll *hyperloglog.Sketch, key string) *HyperBloom {
	return &HyperBloom{
		bloom:    bf,
		hyper:    hll,
		key:      key,
		lastUsed: time.Now().UTC(),
		decay:    config.HyperBloomConfig.Decay,
	}
}

// NewHyperBloomFromParams creates a new HyperBloom with specified capacity and false positive rate
func NewHyperBloomFromParams(capacity uint, falsePositive float64, key string) *HyperBloom {
	bf := bloom.NewWithEstimates(capacity, falsePositive)
	hll := hyperloglog.New() // sparse HyperLogLog (heavier but better for low cardinality use cases)
	return NewHyperBloom(bf, hll, key)
}

// NewDefaultHyperBloom creates a new HyperBloom with default configuration
func NewDefaultHyperBloom(key string) *HyperBloom {
	capacity := config.HyperBloomConfig.Cardinality
	falsePositive := config.HyperBloomConfig.FalsePositive
	return NewHyperBloomFromParams(capacity, falsePositive, key)
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

// GetOrFetchHyperBloom retrieves a HyperBloom instance from the collection or fetches it from the database
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

// GetHyperBloom retrieves a HyperBloom instance by key from the collection
func (dbs *HyperBlooms) GetHyperBloom(key string) (*HyperBloom, bool) {
	db, ok := dbs.blooms[key]
	return db, ok
}

// Bloom returns the Bloom filter of a HyperBloom instance
func (db *HyperBloom) Bloom() *bloom.BloomFilter {
	return db.bloom
}

// Hyper returns the HyperLogLog sketch of a HyperBloom instance
func (db *HyperBloom) Hyper() *hyperloglog.Sketch {
	return db.hyper
}

// Key returns the key of a HyperBloom instance
func (db *HyperBloom) Key() string {
	return db.key
}

// Decay returns the decay duration of a HyperBloom instance
func (db *HyperBloom) Decay() time.Duration {
	return db.decay
}

// LastUsed returns the last used timestamp of a HyperBloom instance
func (db *HyperBloom) LastUsed() time.Time {
	return db.lastUsed
}

// BitSet returns the BitSet of a HyperBloom's Bloom filter
func (db *HyperBloom) BitSet() *bitset.BitSet {
	return db.bloom.BitSet()
}

// BloomCardinality returns the approximated size of a HyperBloom's Bloom filter
func (db *HyperBloom) BloomCardinality() uint32 {
	return db.bloom.ApproximatedSize()
}

// HyperCardinality returns the estimated cardinality of a HyperBloom's HyperLogLog sketch
func (db *HyperBloom) HyperCardinality() uint64 {
	return db.hyper.Estimate()
}

// SETTERS

// Remove deletes a HyperBloom instance from the HyperBlooms collection by key
func (dbs *HyperBlooms) Remove(key string) {
	delete(dbs.blooms, key)
}

// Set adds a HyperBloom instance to the HyperBlooms collection
func (dbs *HyperBlooms) Set(db *HyperBloom, key string) {
	// Refresh the last used timestamp of the HyperBloom instance
	db.Refresh()

	// Add the HyperBloom instance to the 'blooms' map in HyperBlooms
	dbs.blooms[key] = db
}

// Hash adds a value to both the Bloom filter and HyperLogLog sketch of a HyperBloom instance
func (db *HyperBloom) Hash(value string) {
	// Add the value to the Bloom filter of the HyperBloom instance
	db.bloom = db.bloom.AddString(value)

	// Insert the value into the HyperLogLog sketch of the HyperBloom instance
	db.hyper.Insert([]byte(value))
}

func (db *HyperBloom) Refresh() {
	// Update the lastUsed timestamp of the HyperBloom instance to the current time
	db.lastUsed = time.Now()
}

// MORE LOGICS

// CheckExists checks if a value exists in the Bloom filter of a HyperBloom instance
func (db *HyperBloom) CheckExists(value string) bool {
	return db.bloom.TestString(value)
}

// CheckDecayed checks if a HyperBloom instance has decayed based on the last used timestamp
func (dbs *HyperBlooms) CheckDecayed(key string, timemark time.Time) bool {
	// Retrieve the HyperBloom instance for the given key
	db, ok := dbs.GetHyperBloom(key)

	// If the HyperBloom instance was found, proceed to check decay
	return ok && db.CheckDecayed(timemark)
}

// CheckDecayed checks if a HyperBloom instance has decayed based on the last used timestamp
func (db *HyperBloom) CheckDecayed(timemark time.Time) bool {
	// Calculate the duration since the last used timestamp of the HyperBloom instance
	durationDiff := timemark.Sub(db.lastUsed)

	// Compare the duration difference with the decay threshold of the HyperBloom instance
	return durationDiff >= db.decay
}

// GetBloomFromDB fetches a HyperBloom instance from the database by key
func GetBloomFromDB(key string) (*HyperBloom, error) {
	var err error

	// Define a structure to hold the database query result
	record := &struct {
		Key       string // Key of the HyperBloom instance
		Bloombyte []byte // Serialized Bloom filter data
		Hyperbyte []byte // Serialized HyperLogLog data
		Decay     uint64 // Decay duration in seconds
	}{}

	// Query the database for the HyperBloom instance details
	err = postgres.DbClient.QueryRow(
		`SELECT 
			hb.key,
			decay_sec,
			bloombyte, 
			hyperbyte 
		FROM hyperblooms hb
		JOIN hyperbloom_metadata hb_meta
		ON hb.key = hb_meta.key
		WHERE hb.key = $1`, key).Scan(
		&record.Key,
		&record.Decay,
		&record.Bloombyte,
		&record.Hyperbyte,
	)

	// Handle any errors encountered during the database query
	if err != nil {
		return nil, err
	}

	// Create a new HyperBloom instance with retrieved data
	db := &HyperBloom{
		key:      key,
		hyper:    &hyperloglog.Sketch{},       // Initialize HyperLogLog sketch
		bloom:    &bloom.BloomFilter{},        // Initialize Bloom filter
		decay:    time.Duration(record.Decay), // Convert decay duration from seconds to time.Duration
		lastUsed: time.Now().UTC(),            // Set current time as last used timestamp in UTC
	}

	// Unmarshal the retrieved HyperLogLog data into db.hyper
	err = db.hyper.UnmarshalBinary(record.Hyperbyte)
	if err != nil {
		return nil, err
	}

	// Decode the retrieved Bloom filter data into db.bloom
	err = db.bloom.GobDecode(record.Bloombyte)
	if err != nil {
		return nil, err
	}

	// Return the created HyperBloom instance
	return db, nil
}

// JaccardSimBF calculates the Jaccard similarity between the Bloom filters of two HyperBloom instances
func JaccardSimBF(db1, db2 *HyperBloom) float32 {
	bs1 := db1.BitSet()
	bs2 := db2.BitSet()

	andCardinality := bs1.IntersectionCardinality(bs2)
	orCardinality := bs1.UnionCardinality(bs2)

	return float32(andCardinality) / float32(orCardinality)
}
