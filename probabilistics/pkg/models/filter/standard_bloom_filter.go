package filter

import (
	"gopds/probabilistics/pkg/models/meta"

	"github.com/bits-and-blooms/bloom/v3"
)

type StandardBloomFilter struct {
	core *bloom.BloomFilter
	meta *meta.FilterMeta
}

func (f *StandardBloomFilter) Meta() *meta.FilterMeta {
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

func NewStandardBF(maxCard uint, maxFp float64) *StandardBloomFilter {
	f := &StandardBloomFilter{}
	f.core = bloom.NewWithEstimates(uint(maxCard), maxFp)
	f.meta = meta.NewFilterMeta("standard_bloom", maxCard, maxFp, f.core.K(), "murmur128")
	return f
}
