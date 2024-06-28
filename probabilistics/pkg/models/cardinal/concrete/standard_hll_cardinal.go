package concretecardinal

import (
	abstractmeta "gopds/probabilistics/pkg/models/meta/abstract"
	concretemeta "gopds/probabilistics/pkg/models/meta/concrete"

	"github.com/axiomhq/hyperloglog"
)

type StandardHyperLogLog struct {
	core *hyperloglog.Sketch
	meta *concretemeta.StandardHyperLogLogMeta
}

func (f *StandardHyperLogLog) Serialize() []byte {
	byterepr, err := f.core.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return byterepr
}

func (f *StandardHyperLogLog) Deserialize(byterepr []byte) error {
	err := f.core.UnmarshalBinary(byterepr)
	if err != nil {
		panic(err)
	}
	return nil
}

func (hll *StandardHyperLogLog) Meta() abstractmeta.CardinalMeta {
	return hll.meta
}

func (hll *StandardHyperLogLog) Add(value []byte) error {
	hll.core.Insert(value)
	return nil
}

func (hll *StandardHyperLogLog) AddString(value string) error {
	hll.core.Insert([]byte(value))
	return nil
}

func (hll *StandardHyperLogLog) Cardinality() uint64 {
	return hll.core.Estimate()
}

func NewStandardHLL(sparse bool, precision uint8, pfKey string) *StandardHyperLogLog {
	hll := &StandardHyperLogLog{}
	if precision == 14 && sparse {
		hll.core = hyperloglog.New14()
	} else if precision == 14 && !sparse {
		hll.core = hyperloglog.NewNoSparse()
	} else if precision == 16 && sparse {
		hll.core = hyperloglog.New16()
	} else if precision == 16 && !sparse {
		hll.core = hyperloglog.New16NoSparse()
	}

	hll.meta = concretemeta.NewStandardHLLMeta("standard_hll", pfKey)
	return hll
}
