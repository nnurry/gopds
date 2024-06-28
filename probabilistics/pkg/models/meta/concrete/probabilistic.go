package concretemeta

import (
	"time"
)

type ProbabilisticMeta struct {
	id       uint
	key      string
	decay    time.Duration
	lastUsed time.Time
}

func (m *ProbabilisticMeta) Id() uint {
	return m.id
}

func (m *ProbabilisticMeta) SetId(id uint) {
	m.id = id
}

func (m *ProbabilisticMeta) Key() string {
	return m.key
}

func (m *ProbabilisticMeta) Decay() time.Duration {
	return m.decay
}

func (m *ProbabilisticMeta) LastUsed() time.Time {
	return m.lastUsed
}

func (m *ProbabilisticMeta) SetLastUsed(t time.Time) {
	m.lastUsed = t
}

func (m *ProbabilisticMeta) ResetLastUsed() {
	m.lastUsed = time.Now().UTC()
}

func (m *ProbabilisticMeta) IsDecayed(timemark time.Time) bool {
	return m.lastUsed.Add(m.decay).Compare(timemark) == -1
}

func NewProbabilisticMeta(key string, decay time.Duration) *ProbabilisticMeta {
	return &ProbabilisticMeta{
		key:      key,
		decay:    decay,
		lastUsed: time.Now().UTC(),
	}
}
