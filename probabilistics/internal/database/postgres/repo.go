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
		key VARCHAR NOT NULL,
		type VARCHAR NOT NULL,
		max_cardinality INTEGER NOT NULL,
		max_fp REAL NOT NULL,
		hash_func_num BIGINT,
		hash_func_type VARCHAR,
		UNIQUE (type, key, max_cardinality, max_fp, hash_func_type)
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
		key VARCHAR NOT NULL,
		UNIQUE (type, key)
	)`)

	if err != nil {
		tx.Rollback()
		panic(err)
	}

	_, err = Client.Exec(`
	CREATE TABLE IF NOT EXISTS filter_blob (
		filter_id SERIAL,
		filter_byte BYTEA,
		FOREIGN KEY (filter_id) REFERENCES filters(id),
		UNIQUE (filter_id)
	)`)

	if err != nil {
		tx.Rollback()
		panic(err)
	}

	_, err = Client.Exec(`
	CREATE TABLE IF NOT EXISTS cardinal_blob (
		cardinal_id SERIAL,
		cardinal_byte BYTEA,
		FOREIGN KEY (cardinal_id) REFERENCES cardinals(id),
		UNIQUE (cardinal_id)
	)`)

	if err != nil {
		tx.Rollback()
		panic(err)
	}

	tx.Commit()
}
