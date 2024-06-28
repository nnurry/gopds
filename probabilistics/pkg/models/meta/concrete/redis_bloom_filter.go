package concretemeta

import "math"

type RedisBloomFilterMeta struct {
	id              uint
	filterType      string
	maxCard         uint
	maxFp           float64
	hashFuncNum     uint
	hashFuncType    string
	bfKey           string
	expansionFactor uint
	nonScaling      bool
}

func (m *RedisBloomFilterMeta) Id() uint {
	return m.id
}

func (m *RedisBloomFilterMeta) SetId(id uint) {
	m.id = id
}

func (m *RedisBloomFilterMeta) FilterType() string {
	return m.filterType
}

func (m *RedisBloomFilterMeta) MaxCard() uint {
	return m.maxCard
}

func (m *RedisBloomFilterMeta) MaxFp() float64 {
	return m.maxFp
}

func (m *RedisBloomFilterMeta) HashFuncNum() uint {
	return m.hashFuncNum
}

func (m *RedisBloomFilterMeta) HashFuncType() string {
	return m.hashFuncType
}

func (m *RedisBloomFilterMeta) BFKey() string {
	return m.bfKey
}

func (m *RedisBloomFilterMeta) ExpansionFactor() uint {
	return m.expansionFactor
}

func (m *RedisBloomFilterMeta) NonScaling() bool {
	return m.nonScaling
}

func NewRedisBFMeta(
	maxCard uint, maxFp float64, expansionFactor uint,
	nonScaling bool, bfKey string) *RedisBloomFilterMeta {
	k := math.Ceil(-(math.Log(maxFp) / math.Log(2)))
	return &RedisBloomFilterMeta{
		filterType:      "redis_bloom",
		maxCard:         maxCard,
		maxFp:           maxFp,
		hashFuncNum:     uint(k),
		hashFuncType:    "murmur64",
		bfKey:           bfKey,
		expansionFactor: expansionFactor,
		nonScaling:      nonScaling,
	}
}
