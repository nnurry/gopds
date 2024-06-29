package concretefilter

import (
	abstractmeta "gopds/probabilistics/pkg/models/meta/abstract"
	concretemeta "gopds/probabilistics/pkg/models/meta/concrete"

	"github.com/bits-and-blooms/bloom/v3"
)

type StandardBloomFilter struct {
	core *bloom.BloomFilter
	meta *concretemeta.StandardBloomFilterMeta
}

func (f *StandardBloomFilter) Serialize() []byte {
	byterepr, err := f.core.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return byterepr
}

func (f *StandardBloomFilter) Deserialize(byterepr []byte) error {
	err := f.core.UnmarshalBinary(byterepr)
	if err != nil {
		panic(err)
	}
	return nil
}

func (f *StandardBloomFilter) Meta() abstractmeta.FilterMeta {
	return f.meta
}

func (f *StandardBloomFilter) Add(value []byte) error {
	f.core = f.core.Add(value)
	return nil
}

func (f *StandardBloomFilter) AddString(value string) error {
	f.core = f.core.AddString(value)
	return nil
}

func (f *StandardBloomFilter) Exists(value []byte) bool {
	return f.core.Test(value)
}

func NewStandardBF(maxCard uint, maxFp float64, key string) *StandardBloomFilter {
	f := &StandardBloomFilter{}
	f.core = bloom.NewWithEstimates(uint(maxCard), maxFp)
	f.meta = concretemeta.NewStandardBFMeta(maxCard, maxFp, "murmur128", key)
	return f
}
