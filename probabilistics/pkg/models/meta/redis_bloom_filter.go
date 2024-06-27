package meta

import "math"

type RedisBloomFilterMeta struct {
	FilterMeta
	bfKey           string
	expansionFactor uint
	nonScaling      bool
}

func (cm *RedisBloomFilterMeta) BFKey() string {
	return cm.bfKey
}

func (cm *RedisBloomFilterMeta) ExpansionFactor() uint {
	return cm.expansionFactor
}

func (cm *RedisBloomFilterMeta) NonScaling() bool {
	return cm.nonScaling
}

func NewRedisBFMeta(
	maxCard uint, maxFp float64, expansionFactor uint,
	nonScaling bool, bfKey string) *RedisBloomFilterMeta {
	k := math.Ceil(-(math.Log(maxFp) / math.Log(2)))
	return &RedisBloomFilterMeta{
		FilterMeta: FilterMeta{
			filterType:   "redis_bloom",
			maxCard:      maxCard,
			maxFp:        maxFp,
			hashFuncNum:  uint(k),
			hashFuncType: "murmur64",
		},
		bfKey:           bfKey,
		expansionFactor: expansionFactor,
		nonScaling:      nonScaling,
	}
}
