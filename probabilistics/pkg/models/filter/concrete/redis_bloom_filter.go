package concretefilter

import (
	"bytes"
	"fmt"
	myredis "gopds/probabilistics/internal/database/redis"
	abstractmeta "gopds/probabilistics/pkg/models/meta/abstract"
	concretemeta "gopds/probabilistics/pkg/models/meta/concrete"

	"github.com/redis/go-redis/v9"
)

var SigFlag = []byte("Redis deserialize flag")

type RedisBloomFilter struct {
	core *redis.Client
	meta *concretemeta.RedisBloomFilterMeta
}

func (f *RedisBloomFilter) Serialize() []byte {
	return []byte{}
}

func (f *RedisBloomFilter) Deserialize(byterepr []byte) error {
	if !bytes.Equal(byterepr, SigFlag) {
		panic("Wrong signal")
	}
	f.core = myredis.Client
	return nil
}

func (f *RedisBloomFilter) Meta() abstractmeta.FilterMeta {
	return f.meta
}

func (f *RedisBloomFilter) addAny(value interface{}) (bool, error) {
	val, err := f.core.BFAdd(myredis.Ctx, f.getKey(), value).Result()
	return val, err
}

func (f *RedisBloomFilter) Add(value []byte) error {
	_, err := f.addAny(value)
	return err
}

func (f *RedisBloomFilter) AddString(value string) error {
	_, err := f.addAny(value)
	return err
}

func (f *RedisBloomFilter) Exists(value []byte) bool {
	val, _ := f.core.BFExists(myredis.Ctx, f.getKey(), value).Result()
	return val
}

func (f *RedisBloomFilter) getKey() string {
	return fmt.Sprintf(
		"bloom:key=%s:capacity=%d:error_rate=%f:expansion=%d:scaling=%t",
		f.meta.BFKey(),
		f.meta.MaxCard(),
		f.meta.MaxFp(),
		f.meta.ExpansionFactor(),
		f.meta.NonScaling(),
	)
}

func NewRedisBF(
	maxCard uint, maxFp float64, expansionFactor uint,
	nonScaling bool, fKey string) *RedisBloomFilter {
	myredis.Initialize()
	f := &RedisBloomFilter{}
	f.core = myredis.Client
	f.meta = concretemeta.NewRedisBFMeta(maxCard, maxFp, expansionFactor, nonScaling, fKey)
	err := f.core.BFReserveWithArgs(
		myredis.Ctx,
		f.getKey(),
		&redis.BFReserveOptions{
			Capacity:   int64(maxCard),
			Error:      maxFp,
			Expansion:  int64(expansionFactor),
			NonScaling: nonScaling,
		},
	).Err()

	if err != nil && err.Error() != "ERR item exists" {
		panic(err)
	}
	return f
}
