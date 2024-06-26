package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type RedisConfig struct {
	Addr string `env:"REDIS_ADDR" envDefault:"127.0.0.1:6379"`
}

var redisCfg = RedisConfig{}

func LoadRedisConfig() {
	if err := env.Parse(&redisCfg); err != nil {
		log.Printf("%+v\n", err)
	}
}

func RedisCfg() *RedisConfig {
	return &redisCfg
}
