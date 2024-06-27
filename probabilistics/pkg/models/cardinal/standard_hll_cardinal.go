package cardinal

import (
	"errors"
	"gopds/probabilistics/pkg/models/meta"

	"github.com/axiomhq/hyperloglog"
)

type StandardHyperLogLog struct {
	core *hyperloglog.Sketch
	meta *meta.CardinalMeta
}

func (hll *StandardHyperLogLog) Meta() *meta.CardinalMeta {
	return hll.meta
}

func (hll *StandardHyperLogLog) Add(value []byte) error {
	if ok := hll.core.Insert(value); ok {
		return nil
	}
	return errors.New("Unable to insert " + string(value) + " for some reason ...")
}

func (hll *StandardHyperLogLog) AddString(value string) error {
	if ok := hll.core.Insert([]byte(value)); ok {
		return nil
	}
	return errors.New("Unable to insert " + value + " for some reason ...")
}

func (hll *StandardHyperLogLog) Cardinality() uint64 {
	return hll.core.Estimate()
}

func NewStandardHLL(sparse bool, precision uint8) *StandardHyperLogLog {
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

	hll.meta = meta.NewCardinalMeta("standard_hll")
	return hll
}
