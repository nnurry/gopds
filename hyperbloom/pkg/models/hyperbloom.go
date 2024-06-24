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

// PGBloom is a database model for storing Bloom filters and HyperLogLog sketches
type PGBloom struct {
	Key       string
	Bloombyte []byte
	Hyperbyte []byte
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
		lastUsed: time.Now(),
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

// GetHyperBlooms retrieves all HyperBloom instances from the HyperBlooms collection
func (dbs *HyperBlooms) GetHyperBlooms() []*HyperBloom {
	output := []*HyperBloom{}
	for key := range dbs.blooms {
		db, err := dbs.GetOrFetchHyperBloom(key)
		if err == nil {
			output = append(output, db)
		}
	}
	return output
}

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
func (dbs *HyperBlooms) CheckDecayed(key string) bool {
	currentTime := time.Now().UTC()
	db, ok := dbs.GetHyperBloom(key)
	if !ok {
		return false
	}
	if currentTime.Sub(db.lastUsed) >= db.decay {
		return true
	}
	return false
}

// GetBloomFromDB fetches a HyperBloom instance from the database by key
func GetBloomFromDB(key string) (*HyperBloom, error) {
	var err error
	record := &PGBloom{}
	query := "SELECT key, bloombyte, hyperbyte FROM bloom_filters WHERE key = $1"
	err = postgres.DbClient.QueryRow(query, key).Scan(&record.Key, &record.Bloombyte, &record.Hyperbyte)
	if err != nil {
		return nil, err
	}

	db := &HyperBloom{
		key:   key,
		hyper: &hyperloglog.Sketch{},
		bloom: &bloom.BloomFilter{},
	}

	err = db.hyper.UnmarshalBinary(record.Hyperbyte)
	if err != nil {
		return nil, err
	}

	err = db.bloom.GobDecode(record.Bloombyte)
	if err != nil {
		return nil, err
	}

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
