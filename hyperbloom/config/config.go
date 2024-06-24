package config

import "time"

type HyperBloom struct {
	FalsePositive float64
	Cardinality   uint
	Decay         time.Duration
	UpdateRate    time.Duration
}

var HyperBloomConfig = HyperBloom{
	FalsePositive: 0.01,
	Cardinality:   10000,
	Decay:         time.Second * 60,
	UpdateRate:    time.Second * 15,
}
