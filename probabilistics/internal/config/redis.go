package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type RedisConfig struct {
	Addr string `env:"REDIS_ADDR" envDefault:"redis:6379"`
}

var redisCfg = RedisConfig{}

func LoadRedisConfig() {
	if err := env.Parse(&redisCfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
}

func RedisCfg() *RedisConfig {
	return &redisCfg
}
