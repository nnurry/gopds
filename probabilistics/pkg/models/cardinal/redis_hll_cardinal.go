package cardinal

import (
	myredis "gopds/probabilistics/internal/database/redis"
	"gopds/probabilistics/pkg/models/meta"

	"github.com/redis/go-redis/v9"
)

type RedisHyperLogLog struct {
	core *redis.Client
	meta *meta.RedisHyperLogLogMeta
}

func (hll *RedisHyperLogLog) Meta() *meta.RedisHyperLogLogMeta {
	return hll.meta
}

func (hll *RedisHyperLogLog) addAny(value interface{}) (int64, error) {
	val, err := hll.core.PFAdd(myredis.Ctx, hll.Meta().PFKey(), value).Result()
	return val, err
}

func (hll *RedisHyperLogLog) Add(value []byte) error {
	_, err := hll.addAny(value)
	return err
}

func (hll *RedisHyperLogLog) AddString(value string) error {
	_, err := hll.addAny(value)
	return err
}

func (hll *RedisHyperLogLog) Cardinality() uint64 {
	val, err := hll.core.PFCount(myredis.Ctx, hll.Meta().PFKey()).Result()
	if err != nil {
		return 0
	}
	return uint64(val)
}

func NewRedisHLL(pfKey string) *RedisHyperLogLog {
	hll := &RedisHyperLogLog{}
	myredis.Initialize()
	hll.core = myredis.Client
	hll.meta = meta.NewRedisHLLMeta("redis_hll", pfKey)
	return hll
}
