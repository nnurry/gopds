package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type ProbabilisticConfig struct {
	DecayDuration time.Duration `env:"PROB_DECAY" envDefault:"120s"`
	SyncInterval  time.Duration `env:"PROB_SYNC" envDefault:"20s"`
}

var ProbabilisticCfg = ProbabilisticConfig{}

func init() {
	if err := env.Parse(&ProbabilisticCfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
}
