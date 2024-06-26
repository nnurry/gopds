package config

import (
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env"
)

// ApplicationConfig holds configuration related to the application's HTTP server.
type ApplicationConfig struct {
	Addr        string      `env:"MUX_ADDR" envDefault:":5000"` // Addr is the address the HTTP server listens on.
	InfoLogger  *log.Logger // InfoLogger is the logger for informational messages.
	ErrorLogger *log.Logger // ErrorLogger is the logger for error messages.
}

// PostgresConfig holds configuration related to PostgreSQL database connection.
type PostgresConfig struct {
	Host     string `env:"DB_HOST" envDefault:"127.0.0.1"` // Host is the PostgreSQL server host.
	Port     int    `env:"DB_PORT" envDefault:"5432"`      // Port is the PostgreSQL server port.
	Database string `env:"DB_NAME" envDefault:"postgres"`  // Database is the name of the PostgreSQL database.
	Username string `env:"DB_USER" envDefault:"admin"`     // Username is the username for PostgreSQL authentication.
	Password string `env:"DB_PASS" envDefault:"123"`       // Password is the password for PostgreSQL authentication.
	SSLMode  string `env:"DB_SSL" envDefault:"disable"`    // SSLMode specifies whether to use SSL for PostgreSQL connection.
}

// HyperBloomConfig holds configuration specific to HyperBloom.
type HyperBloomConfig struct {
	FalsePositive float64       `env:"HB_FP" envDefault:"0.0081"`       // FalsePositive is the desired false positive rate for HyperBloom.
	Cardinality   uint          `env:"HB_CARD" envDefault:"10000"`      // Cardinality is the expected number of elements to be stored in HyperBloom.
	Decay         time.Duration `env:"HB_DECAY" envDefault:"120s"`      // Decay is the decay period for HyperBloom data.
	UpdateRate    time.Duration `env:"HB_UPDATE_RATE" envDefault:"20s"` // UpdateRate is the rate at which HyperBloom should be updated.
}

// Global variables holding the loaded configurations.
var (
	HyperBloomCfg  HyperBloomConfig  // HyperBloomCfg holds the loaded HyperBloom configuration.
	PostgresCfg    PostgresConfig    // PostgresCfg holds the loaded PostgreSQL configuration.
	ApplicationCfg ApplicationConfig // ApplicationCfg holds the loaded application configuration.
)

// LoadConfigPostgres loads PostgreSQL configuration from environment variables.
func LoadConfigPostgres() {
	if err := env.Parse(&PostgresCfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
}

// LoadConfigHyperBloom loads HyperBloom configuration from environment variables.
func LoadConfigHyperBloom() {
	if err := env.Parse(&HyperBloomCfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
}

// LoadConfigApplication loads application configuration from environment variables.
func LoadConfigApplication() {
	if err := env.Parse(&ApplicationCfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
}

// GetDataSourceName constructs and returns the data source name for PostgreSQL connection.
func (cfg PostgresConfig) GetDataSourceName() string {
	baseStr := "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s"
	return fmt.Sprintf(
		baseStr,
		cfg.Host, cfg.Port, cfg.Username,
		cfg.Password, cfg.Database, cfg.SSLMode,
	)
}
