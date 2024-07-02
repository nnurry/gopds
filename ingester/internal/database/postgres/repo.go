package postgres

import (
	_ "github.com/lib/pq"
)

func Bootstrap() {
	var err error

	tx, _ := Client.Begin()

	_, err = Client.Exec(`
	CREATE TABLE IF NOT EXISTS raw_data (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT timezone('UTC', now()),
		raw_json JSONB
	);`)

	if err != nil {
		tx.Rollback()
		panic(err)
	}

	_, err = Client.Exec(`
		CREATE TABLE IF NOT EXISTS json_config (
		id SERIAL PRIMARY KEY,
		created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT timezone('UTC', now()),
		key_path VARCHAR NOT NULL,
		value_path VARCHAR NOT NULL,
		cardinal_config JSONB NOT NULL,
		filter_config JSONB NOT NULL
	);`)

	if err != nil {
		tx.Rollback()
		panic(err)
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		panic(err)
	}

}
