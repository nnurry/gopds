package postgres

import (
	"database/sql"
	"fmt"
	"gopds/probabilistics/internal/config"
	"sync"
)

var Client *sql.DB

var Initialize = sync.OnceFunc(func() {

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

	fmt.Println("Successfully connected to PostgreSQL!")

})
