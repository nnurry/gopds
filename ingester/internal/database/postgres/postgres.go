package postgres

import (
	"database/sql"
	"log"
)

var Client *sql.DB

func init() {
	var err error

	Client, err = sql.Open(
		"postgres",
		"host=127.0.0.1 port=5432 user=admin password=123 dbname=postgres sslmode=disable",
	)

	if err != nil {
		panic(err)
	}

	err = Client.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("Successfully connected to PostgreSQL!")

}
