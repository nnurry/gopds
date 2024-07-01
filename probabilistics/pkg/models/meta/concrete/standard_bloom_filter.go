package concretemeta

import (
	"github.com/bits-and-blooms/bloom/v3"
)

type StandardBloomFilterMeta struct {
	id           uint
	key          string
	filterType   string
	maxCard      uint
	maxFp        float64
	hashFuncNum  uint
	hashFuncType string
}

func (m *StandardBloomFilterMeta) Id() uint {
	return m.id
}

func (m *StandardBloomFilterMeta) SetId(id uint) {
	m.id = id
}

func (m *StandardBloomFilterMeta) FilterType() string {
	return m.filterType
}

func (m *StandardBloomFilterMeta) MaxCard() uint {
	return m.maxCard
}

func (m *StandardBloomFilterMeta) MaxFp() float64 {
	return m.maxFp
}

func (m *StandardBloomFilterMeta) HashFuncNum() uint {
	return m.hashFuncNum
}

func (m *StandardBloomFilterMeta) HashFuncType() string {
	return m.hashFuncType
}

func (m *StandardBloomFilterMeta) Key() string {
	return m.key
}

func NewStandardBFMeta(
	maxCard uint, maxFp float64, hashFuncType string, key string) *StandardBloomFilterMeta {
	_, hashFuncNum := bloom.EstimateParameters(maxCard, maxFp)
	return &StandardBloomFilterMeta{
		key:          key,
		filterType:   "STANDARD_BLOOM",
		maxCard:      maxCard,
		maxFp:        maxFp,
		hashFuncNum:  hashFuncNum,
		hashFuncType: hashFuncType,
	}
}
