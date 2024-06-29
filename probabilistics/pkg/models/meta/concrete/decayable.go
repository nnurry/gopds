package concretemeta

import (
	"time"
)

type DecayableMeta struct {
	decay    time.Duration
	lastUsed time.Time
}

func (m *DecayableMeta) Decay() time.Duration {
	return m.decay
}

func (m *DecayableMeta) LastUsed() time.Time {
	return m.lastUsed
}

func (m *DecayableMeta) SetLastUsed(t time.Time) {
	m.lastUsed = t
}

func (m *DecayableMeta) ResetLastUsed() {
	m.lastUsed = time.Now().UTC()
}

func (m *DecayableMeta) IsDecayed(timemark time.Time) bool {
	return m.lastUsed.Add(m.decay).Compare(timemark) == -1
}

func NewProbabilisticMeta(decay time.Duration) *DecayableMeta {
	return &DecayableMeta{
		decay:    decay,
		lastUsed: time.Now().UTC(),
	}
}
