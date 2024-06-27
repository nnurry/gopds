package meta

import (
	"time"
)

type CardinalFilterMeta struct {
	key      string
	decay    time.Duration
	lastUsed time.Time
}

func (cfm *CardinalFilterMeta) Key() string {
	return cfm.key
}

func (cfm *CardinalFilterMeta) Decay() time.Duration {
	return cfm.decay
}

func (cfm *CardinalFilterMeta) LastUsed() time.Time {
	return cfm.lastUsed
}

func NewCardinalFilterMeta(key string, decay time.Duration) *CardinalFilterMeta {
	return &CardinalFilterMeta{
		key:      key,
		decay:    decay,
		lastUsed: time.Now().UTC(),
	}
}
