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
	Decay:         time.Second * 300,
	UpdateRate:    time.Second * 20,
}
