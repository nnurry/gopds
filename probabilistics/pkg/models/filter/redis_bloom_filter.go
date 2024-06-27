package filter

import (
	myredis "gopds/probabilistics/internal/database/redis"
	"gopds/probabilistics/pkg/models/meta"

	"github.com/redis/go-redis/v9"
)

type RedisBloomFilter struct {
	core *redis.Client
	meta *meta.RedisBloomFilterMeta
}

func (hll *RedisBloomFilter) Meta() *meta.RedisBloomFilterMeta {
	return hll.meta
}

func (hll *RedisBloomFilter) addAny(value interface{}) (int64, error) {
	val, err := hll.core.PFAdd(myredis.Ctx, hll.Meta().BFKey(), value).Result()
	return val, err
}

func (hll *RedisBloomFilter) Add(value []byte) error {
	_, err := hll.addAny(value)
	return err
}

func (hll *RedisBloomFilter) AddString(value string) error {
	_, err := hll.addAny(value)
	return err
}

func (hll *RedisBloomFilter) Cardinality() uint64 {
	val, err := hll.core.PFCount(myredis.Ctx, hll.Meta().BFKey()).Result()
	if err != nil {
		return 0
	}
	return uint64(val)
}

func NewRedisBF(
	maxCard uint, maxFp float64, expansionFactor uint,
	nonScaling bool, bfKey string) *RedisBloomFilter {
	hll := &RedisBloomFilter{}
	hll.core = myredis.Client
	hll.meta = meta.NewRedisBFMeta(maxCard, maxFp, expansionFactor, nonScaling, bfKey)
	err := hll.core.BFReserveWithArgs(
		myredis.Ctx,
		bfKey,
		&redis.BFReserveOptions{
			Capacity:   int64(maxCard),
			Error:      maxFp,
			Expansion:  int64(expansionFactor),
			NonScaling: nonScaling,
		},
	)
	if err != nil {
		return nil
	}
	return hll
}
