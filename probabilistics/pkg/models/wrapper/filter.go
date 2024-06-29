package wrapper

import (
	"gopds/probabilistics/pkg/models/decayable"
)

type FilterKey struct {
	Type           string
	Key            string
	MaxCardinality uint
	ErrorRate      float64
}

type FilterWrapper struct {
	core    map[FilterKey]*decayable.Filter
	counter uint
}

func NewFilterWrapper() *FilterWrapper {
	return &FilterWrapper{
		core:    make(map[FilterKey]*decayable.Filter),
		counter: 0,
	}
}

func (pw *FilterWrapper) Core() map[FilterKey]*decayable.Filter {
	return pw.core
}

func (pw *FilterWrapper) Counter() uint {
	return pw.counter
}

func (pw *FilterWrapper) Add(k FilterKey, v *decayable.Filter) {
	_, exists := pw.core[k]
	pw.core[k] = v
	if !exists {
		pw.counter++
	}
}

func (pw *FilterWrapper) Delete(k FilterKey) {
	_, exists := pw.core[k]
	delete(pw.core, k)
	if exists {
		pw.counter--
	}
}
