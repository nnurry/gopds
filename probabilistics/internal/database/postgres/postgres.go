package postgres

import (
	"database/sql"
	"log"

	"github.com/nnurry/gopds/probabilistics/internal/config"
)

var Client *sql.DB

func init() {
	var err error
	config.LoadPostgresConfig()

	Client, err = sql.Open("postgres", config.PostgresCfg().GetDataSourceName())
	if err != nil {
		panic(err)
	}

	err = Client.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("Successfully connected to PostgreSQL!")

}
