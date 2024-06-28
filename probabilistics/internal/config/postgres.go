package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type PostgresConfig struct {
	Host     string `env:"DB_HOST" envDefault:"127.0.0.1"`
	Port     int    `env:"DB_PORT" envDefault:"5432"`
	Database string `env:"DB_NAME" envDefault:"postgres"`
	Username string `env:"DB_USER" envDefault:"admin"`
	Password string `env:"DB_PASS" envDefault:"123"`
	SSLMode  string `env:"DB_SSL" envDefault:"disable"`
}

var postgresCfg = PostgresConfig{}

func LoadPostgresConfig() {
	if err := env.Parse(&postgresCfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
}

func PostgresCfg() *PostgresConfig {
	return &postgresCfg
}

func (cfg *PostgresConfig) GetDataSourceName() string {
	baseStr := "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s"
	return fmt.Sprintf(
		baseStr,
		cfg.Host, cfg.Port, cfg.Username,
		cfg.Password, cfg.Database, cfg.SSLMode,
	)
}
