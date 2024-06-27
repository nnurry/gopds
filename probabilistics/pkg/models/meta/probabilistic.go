package meta

import (
	"time"
)

type ProbabilisticMeta struct {
	key      string
	decay    time.Duration
	lastUsed time.Time
}

func (cfm *ProbabilisticMeta) Key() string {
	return cfm.key
}

func (cfm *ProbabilisticMeta) Decay() time.Duration {
	return cfm.decay
}

func (cfm *ProbabilisticMeta) LastUsed() time.Time {
	return cfm.lastUsed
}

func (cfm *ProbabilisticMeta) IsDecayed(timemark time.Time) bool {
	return cfm.lastUsed.Add(cfm.decay).Compare(timemark) == -1
}

func NewProbabilisticMeta(key string, decay time.Duration) *ProbabilisticMeta {
	return &ProbabilisticMeta{
		key:      key,
		decay:    decay,
		lastUsed: time.Now().UTC(),
	}
}
