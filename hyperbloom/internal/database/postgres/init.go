package postgres

import (
	"database/sql"
	"fmt"
	"gopds/hyperbloom/internal/config"

	// Importing "github.com/lib/pq" for PostgreSQL driver
	_ "github.com/lib/pq"
)

// DbClient is a global variable holding the connection pool to the PostgreSQL database.
var DbClient *sql.DB

// init initializes the PostgreSQL database connection upon package initialization.
func init() {
	// Load PostgreSQL configuration from environment or configuration files
	config.LoadConfigPostgres()

	// Retrieve the connection string from the loaded configuration
	connStr := config.PostgresCfg.GetDataSourceName()

	// Open a connection to the PostgreSQL database using the retrieved connection string
	var err error
	DbClient, err = sql.Open("postgres", connStr)
	if err != nil {
		// Panic if there's an error establishing the database connection
		panic(err)
	}

	// Ping the database to verify connectivity
	err = DbClient.Ping()
	if err != nil {
		// Panic if there's an error pinging the database
		panic(err)
	}

	// Print a success message to indicate a successful database connection
	fmt.Println("Successfully connected to PostgreSQL!")
}
