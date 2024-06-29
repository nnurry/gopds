package wrapper

import (
	"gopds/probabilistics/pkg/models/decayable"
)

type CardinalKey struct {
	Type string
	Key  string
}

type CardinalWrapper struct {
	core    map[CardinalKey]*decayable.Cardinal
	counter uint
}

func NewCardinalWrapper() *CardinalWrapper {
	return &CardinalWrapper{
		core:    make(map[CardinalKey]*decayable.Cardinal),
		counter: 0,
	}
}

func (pw *CardinalWrapper) Core() map[CardinalKey]*decayable.Cardinal {
	return pw.core
}

func (pw *CardinalWrapper) Counter() uint {
	return pw.counter
}

func (pw *CardinalWrapper) Add(k CardinalKey, v *decayable.Cardinal) {
	_, exists := pw.core[k]
	pw.core[k] = v
	if !exists {
		pw.counter++
	}
}

func (pw *CardinalWrapper) Delete(k CardinalKey) {
	_, exists := pw.core[k]
	delete(pw.core, k)
	if exists {
		pw.counter--
	}
}
