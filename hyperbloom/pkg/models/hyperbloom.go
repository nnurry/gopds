// Package models defines data structures and methods related to HyperBloom,
// which combines a Bloom filter and a HyperLogLog sketch for efficient data
// membership testing and cardinality estimation.
package models

import (
	"time"

	"gopds/hyperbloom/internal/config"
	"gopds/hyperbloom/internal/database/postgres"

	"github.com/axiomhq/hyperloglog"
	"github.com/bits-and-blooms/bitset"
	"github.com/bits-and-blooms/bloom/v3"
)

// HyperBloom represents a data structure combining a Bloom filter and a HyperLogLog sketch.
// It supports operations for hashing values, checking existence, cardinality estimation,
// and database serialization.
type HyperBloom struct {
	bloom    *bloom.BloomFilter  // Bloom filter for membership testing
	hyper    *hyperloglog.Sketch // HyperLogLog sketch for cardinality estimation
	key      string              // Unique identifier for the HyperBloom instance
	decay    time.Duration       // Time duration after which the instance is considered decayed
	lastUsed time.Time           // Timestamp of the last operation on the instance
}

// NewHyperBloom creates a new HyperBloom instance initialized with given Bloom filter,
// HyperLogLog sketch, and metadata.
func NewHyperBloom(bf *bloom.BloomFilter, hll *hyperloglog.Sketch, key string) *HyperBloom {
	return &HyperBloom{
		bloom:    bf,
		hyper:    hll,
		key:      key,
		lastUsed: time.Now().UTC(),
		decay:    config.HyperBloomCfg.Decay,
	}
}

// NewHyperBloomFromParams creates a new HyperBloom instance with specified capacity,
// false positive rate, and metadata.
func NewHyperBloomFromParams(capacity uint, falsePositive float64, key string) *HyperBloom {
	bf := bloom.NewWithEstimates(capacity, falsePositive)
	hll := hyperloglog.New()
	return NewHyperBloom(bf, hll, key)
}

// NewDefaultHyperBloom creates a new HyperBloom instance with default configuration
// specified in the application's configuration.
func NewDefaultHyperBloom(key string) *HyperBloom {
	capacity := config.HyperBloomCfg.Cardinality
	falsePositive := config.HyperBloomCfg.FalsePositive
	return NewHyperBloomFromParams(capacity, falsePositive, key)
}

// GETTERS

// Bloom returns the Bloom filter instance of the HyperBloom.
func (db *HyperBloom) Bloom() *bloom.BloomFilter {
	return db.bloom
}

// Hyper returns the HyperLogLog sketch instance of the HyperBloom.
func (db *HyperBloom) Hyper() *hyperloglog.Sketch {
	return db.hyper
}

// Key returns the unique identifier (key) of the HyperBloom.
func (db *HyperBloom) Key() string {
	return db.key
}

// Decay returns the decay duration after which the HyperBloom instance is considered decayed.
func (db *HyperBloom) Decay() time.Duration {
	return db.decay
}

// LastUsed returns the timestamp of the last operation on the HyperBloom instance.
func (db *HyperBloom) LastUsed() time.Time {
	return db.lastUsed
}

// BitSet returns the underlying BitSet of the Bloom filter in the HyperBloom instance.
func (db *HyperBloom) BitSet() *bitset.BitSet {
	return db.bloom.BitSet()
}

// BloomCardinality returns the estimated cardinality of the Bloom filter in the HyperBloom instance.
func (db *HyperBloom) BloomCardinality() uint32 {
	return db.bloom.ApproximatedSize()
}

// HyperCardinality returns the estimated cardinality of the HyperLogLog sketch in the HyperBloom instance.
func (db *HyperBloom) HyperCardinality() uint64 {
	return db.hyper.Estimate()
}

// SETTERS

// Hash adds a value to both the Bloom filter and HyperLogLog sketch of the HyperBloom instance.
func (db *HyperBloom) Hash(value string) {
	db.bloom.AddString(value)
	db.hyper.Insert([]byte(value))
}

// Refresh updates the last used timestamp of the HyperBloom instance to the current time.
func (db *HyperBloom) Refresh() {
	db.lastUsed = time.Now()
}

// MORE LOGICS

// CheckExists checks if a value exists in the Bloom filter of the HyperBloom instance.
func (db *HyperBloom) CheckExists(value string) bool {
	return db.bloom.TestString(value)
}

// CheckDecayed checks if the HyperBloom instance has decayed based on the last used timestamp.
func (db *HyperBloom) CheckDecayed(timemark time.Time) bool {
	durationDiff := timemark.Sub(db.lastUsed)
	return durationDiff >= db.decay
}

// GetBloomFromDB fetches a HyperBloom instance from the database by its unique key.
func GetBloomFromDB(key string) (*HyperBloom, error) {
	var err error

	// Query the database for the serialized data of the HyperBloom instance.
	record := &struct {
		Key       string // Unique key of the HyperBloom instance
		Bloombyte []byte // Serialized data of the Bloom filter
		Hyperbyte []byte // Serialized data of the HyperLogLog sketch
		Decay     uint64 // Decay duration in seconds
	}{}

	err = postgres.DbClient.QueryRow(
		`SELECT 
			hb.key,
			decay_sec,
			bloombyte, 
			hyperbyte 
		FROM hyperblooms hb
		JOIN hyperblooms_metadata hb_meta
		ON hb.key = hb_meta.key
		WHERE hb.key = $1`, key).Scan(
		&record.Key,
		&record.Decay,
		&record.Bloombyte,
		&record.Hyperbyte,
	)

	if err != nil {
		return nil, err
	}

	// Create a new HyperBloom instance and populate it with the deserialized data.
	db := &HyperBloom{
		key:      key,
		hyper:    &hyperloglog.Sketch{},
		bloom:    &bloom.BloomFilter{},
		decay:    time.Duration(record.Decay),
		lastUsed: time.Now().UTC(),
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

// JaccardSimBF calculates the Jaccard similarity between the Bloom filters of two HyperBloom instances.
func JaccardSimBF(db1, db2 *HyperBloom) float32 {
	bs1 := db1.BitSet()
	bs2 := db2.BitSet()

	andCardinality := bs1.IntersectionCardinality(bs2)
	orCardinality := bs1.UnionCardinality(bs2)

	return float32(andCardinality) / float32(orCardinality)
}
