package concretecardinal

import (
	"bytes"
	"fmt"
	myredis "gopds/probabilistics/internal/database/redis"
	abstractmeta "gopds/probabilistics/pkg/models/meta/abstract"
	concretemeta "gopds/probabilistics/pkg/models/meta/concrete"

	"github.com/redis/go-redis/v9"
)

var SigFlag = []byte{}

type RedisHyperLogLog struct {
	core *redis.Client
	meta *concretemeta.RedisHyperLogLogMeta
}

func (f *RedisHyperLogLog) Serialize() []byte {
	return []byte{}
}

func (f *RedisHyperLogLog) Deserialize(byterepr []byte) error {
	if !bytes.Equal(byterepr, SigFlag) {
		panic("Wrong signal")
	}
	f.core = myredis.Client
	return nil
}

func (hll *RedisHyperLogLog) Meta() abstractmeta.CardinalMeta {
	return hll.meta
}

func (hll *RedisHyperLogLog) addAny(value []byte) (int64, error) {
	val, err := hll.core.PFAdd(myredis.Ctx, hll.getKey(), value).Result()
	return val, err
}

func (hll *RedisHyperLogLog) Add(value []byte) error {
	_, err := hll.addAny(value)
	return err
}

func (hll *RedisHyperLogLog) AddString(value string) error {
	_, err := hll.addAny([]byte(value))
	return err
}

func (hll *RedisHyperLogLog) Cardinality() uint64 {
	val, err := hll.core.PFCount(myredis.Ctx, hll.getKey()).Result()
	if err != nil {
		return 0
	}
	return uint64(val)
}

func (f *RedisHyperLogLog) getKey() string {
	return fmt.Sprintf("hll:key=%s", f.meta.Key())
}

func NewRedisHLL(key string) *RedisHyperLogLog {
	hll := &RedisHyperLogLog{}
	myredis.Initialize()
	hll.core = myredis.Client
	hll.meta = concretemeta.NewRedisHLLMeta("redis_hll", key)
	return hll
}
