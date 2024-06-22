package models

import (
	"time"

	"github.com/bits-and-blooms/bitset"
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/nnurry/gopds/hyperbloom/config"
	"github.com/nnurry/gopds/hyperbloom/internal/database/postgres"
)

type DecayBlooms struct {
	blooms map[string](*DecayBloom)
}
type DecayBloom struct {
	bloom    *bloom.BloomFilter
	key      string
	decay    time.Duration
	lastUsed time.Time
}

type PGBloom struct {
	Key       string
	Bloombyte []byte
}

func NewDecayBlooms() *DecayBlooms {
	return &DecayBlooms{
		blooms: make(map[string]*DecayBloom),
	}
}

func NewDecayBloom(bf *bloom.BloomFilter, key string) *DecayBloom {
	return &DecayBloom{
		bloom:    bf,
		key:      key,
		lastUsed: time.Now(),
		decay:    config.DecayBloomConfig.Decay,
	}
}

func NewDecayBloomFromParams(capacity uint, falsePositive float64, key string) *DecayBloom {
	bf := bloom.NewWithEstimates(capacity, falsePositive)
	return NewDecayBloom(bf, key)
}

func NewDefaultDecayBloom(key string) *DecayBloom {
	capacity := config.DecayBloomConfig.Cardinality
	falsePositive := config.DecayBloomConfig.FalsePositive
	return NewDecayBloomFromParams(capacity, falsePositive, key)
}

// GETTERS

func (dbs *DecayBlooms) GetBlooms() []*DecayBloom {
	output := []*DecayBloom{}
	for key := range dbs.blooms {
		db := dbs.blooms[key]
		output = append(output, db)
	}
	return output
}

func (dbs *DecayBlooms) GetBloom(key string) (*DecayBloom, bool) {
	db, ok := dbs.blooms[key]
	return db, ok
}

func (db *DecayBloom) Bloom() *bloom.BloomFilter {
	return db.bloom
}

func (db *DecayBloom) Key() string {
	return db.key
}

func (db *DecayBloom) Decay() time.Duration {
	return db.decay
}

func (db *DecayBloom) LastUsed() time.Time {
	return db.lastUsed
}

func (db *DecayBloom) BitSet() *bitset.BitSet {
	return db.bloom.BitSet()
}

func (db *DecayBloom) ApproximatedSize() uint32 {
	return db.bloom.ApproximatedSize()
}

// SETTERS

func (dbs *DecayBlooms) Remove(key string) {
	delete(dbs.blooms, key)
}

func (dbs *DecayBlooms) Set(db *DecayBloom, key string) {
	dbs.blooms[key] = db
}

func (db *DecayBloom) Hash(value string) {
	db.bloom = db.bloom.AddString(value)
}

// MORE LOGICS

func (db *DecayBloom) CheckExists(value string) bool {
	return db.bloom.TestString(value)
}

func (dbs *DecayBlooms) CheckDecayed(key string) bool {
	currentTime := time.Now().UTC()
	db, ok := dbs.GetBloom(key)
	if !ok {
		return false
	}
	if currentTime.Sub(db.lastUsed) >= db.decay {
		return true
	}
	return false
}

func GetBloomFromDB(key string) (*DecayBloom, error) {
	var err error
	record := &PGBloom{}
	query := "SELECT key, bloombyte FROM bloom_filters WHERE key = $1"
	err = postgres.DbClient.QueryRow(query, key).Scan(&record.Key, &record.Bloombyte)
	if err != nil {
		return nil, err
	}

	db := &DecayBloom{
		key: key,
		bloom: bloom.NewWithEstimates(
			config.DecayBloomConfig.Cardinality,
			config.DecayBloomConfig.FalsePositive,
		),
	}

	err = db.bloom.GobDecode(record.Bloombyte)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func JaccardSimBF(db1, db2 *DecayBloom) float32 {
	bs1 := db1.BitSet()
	bs2 := db2.BitSet()

	andCardinality := bs1.IntersectionCardinality(bs2)
	orCardinality := bs1.UnionCardinality(bs2)

	return float32(andCardinality) / float32(orCardinality)
}
