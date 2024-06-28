package postgres

import (
	_ "github.com/lib/pq"
)

func Bootstrap() {
	var err error

	tx, err := Client.Begin()

	if err != nil {
		panic(err)
	}

	_, err = Client.Exec(`
	CREATE TABLE IF NOT EXISTS filters (
		id SERIAL PRIMARY KEY,
		type VARCHAR NOT NULL,
		max_cardinality INTEGER NOT NULL,
		max_fp REAL NOT NULL,
		hash_func_num BIGINT,
		hash_func_type VARCHAR,
		UNIQUE (type, max_cardinality, max_fp, hash_func_type)
	)`)

	// Note: No need to create index yet

	if err != nil {
		tx.Rollback()
		panic(err)
	}

	_, err = Client.Exec(`
	CREATE TABLE IF NOT EXISTS cardinals (
		id SERIAL PRIMARY KEY,
		type VARCHAR NOT NULL,
		UNIQUE (type)
	)`)

	if err != nil {
		tx.Rollback()
		panic(err)
	}

	_, err = Client.Exec(`
	CREATE TABLE IF NOT EXISTS wrappers (
		filter_id SERIAL,
		cardinal_id SERIAL,
		key VARCHAR NOT NULL,
		filter_byte BYTEA,
		cardinal_byte BYTEA,
		FOREIGN KEY (filter_id) REFERENCES filters(id),
		FOREIGN KEY (cardinal_id) REFERENCES cardinals(id),
		UNIQUE (filter_id, cardinal_id)
	)`)

	if err != nil {
		tx.Rollback()
		panic(err)
	}

	tx.Commit()
}
