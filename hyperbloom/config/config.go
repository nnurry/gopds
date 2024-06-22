package config

import "time"

type DecayBloom struct {
	FalsePositive float64
	Cardinality   uint
	Decay         time.Duration
}

var DecayBloomConfig = DecayBloom{
	FalsePositive: 0.01,
	Cardinality:   10000,
	Decay:         time.Second * 60,
}
