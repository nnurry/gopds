package config

import "time"

// App specific config for Hyperbloom.
// We control the cardinality, FP and other factors here
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
