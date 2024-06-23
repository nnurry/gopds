package models

import (
	"time"

	"github.com/axiomhq/hyperloglog"
	"github.com/bits-and-blooms/bitset"
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/nnurry/gopds/hyperbloom/config"
	"github.com/nnurry/gopds/hyperbloom/internal/database/postgres"
)

type HyperBlooms struct {
	blooms map[string](*HyperBloom)
}
type HyperBloom struct {
	bloom    *bloom.BloomFilter
	hyper    *hyperloglog.Sketch
	key      string
	decay    time.Duration
	lastUsed time.Time
}

type PGBloom struct {
	Key       string
	Bloombyte []byte
	Hyperbyte []byte
}

func NewHyperBlooms() *HyperBlooms {
	return &HyperBlooms{
		blooms: make(map[string]*HyperBloom),
	}
}

func NewHyperBloom(bf *bloom.BloomFilter, hll *hyperloglog.Sketch, key string) *HyperBloom {
	return &HyperBloom{
		bloom:    bf,
		hyper:    hll,
		key:      key,
		lastUsed: time.Now(),
		decay:    config.HyperBloomConfig.Decay,
	}
}

func NewHyperBloomFromParams(capacity uint, falsePositive float64, key string) *HyperBloom {
	bf := bloom.NewWithEstimates(capacity, falsePositive)
	hll := hyperloglog.New() // sparse Hyperloglog (heavier but better against low cardinality, as our usecase)
	return NewHyperBloom(bf, hll, key)
}

func NewDefaultHyperBloom(key string) *HyperBloom {
	capacity := config.HyperBloomConfig.Cardinality
	falsePositive := config.HyperBloomConfig.FalsePositive
	return NewHyperBloomFromParams(capacity, falsePositive, key)
}

// GETTERS

func (dbs *HyperBlooms) GetHyperBlooms() []*HyperBloom {
	output := []*HyperBloom{}
	for key := range dbs.blooms {
		db := dbs.blooms[key]
		output = append(output, db)
	}
	return output
}

func (dbs *HyperBlooms) GetHyperBloom(key string) (*HyperBloom, bool) {
	db, ok := dbs.blooms[key]
	return db, ok
}

func (db *HyperBloom) Bloom() *bloom.BloomFilter {
	return db.bloom
}

func (db *HyperBloom) Hyper() *hyperloglog.Sketch {
	return db.hyper
}

func (db *HyperBloom) Key() string {
	return db.key
}

func (db *HyperBloom) Decay() time.Duration {
	return db.decay
}

func (db *HyperBloom) LastUsed() time.Time {
	return db.lastUsed
}

func (db *HyperBloom) BitSet() *bitset.BitSet {
	return db.bloom.BitSet()
}

func (db *HyperBloom) BloomCardinality() uint32 {
	return db.bloom.ApproximatedSize()
}

func (db *HyperBloom) HyperCardinality() uint64 {
	return db.hyper.Estimate()
}

// SETTERS

func (dbs *HyperBlooms) Remove(key string) {
	delete(dbs.blooms, key)
}

func (dbs *HyperBlooms) Set(db *HyperBloom, key string) {
	dbs.blooms[key] = db
}

func (db *HyperBloom) Hash(value string) {
	db.bloom = db.bloom.AddString(value)
	db.hyper.Insert([]byte(value))
}

// MORE LOGICS

func (db *HyperBloom) CheckExists(value string) bool {
	return db.bloom.TestString(value)
}

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

func JaccardSimBF(db1, db2 *HyperBloom) float32 {
	bs1 := db1.BitSet()
	bs2 := db2.BitSet()

	andCardinality := bs1.IntersectionCardinality(bs2)
	orCardinality := bs1.UnionCardinality(bs2)

	return float32(andCardinality) / float32(orCardinality)
}
