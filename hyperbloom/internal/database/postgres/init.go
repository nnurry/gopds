package postgres

import (
	"database/sql"
	"fmt"

	//for connecting to db
	_ "github.com/lib/pq"
)

const (
	host     = "postgres" // replace to IP address, domain name depending on use case
	port     = 5432
	user     = "admin"    // replace with username
	password = "123"      // replace with password
	dbname   = "postgres" // replace with the database name
)

var (
	DbClient *sql.DB
)

func init() {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	var err error
	DbClient, err = sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}

	err = DbClient.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
}
