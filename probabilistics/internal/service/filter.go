package service

import (
	"database/sql"
	"errors"

	request_schema "github.com/nnurry/gopds/probabilistics/internal/api/rest/schemas/request"
	"github.com/nnurry/gopds/probabilistics/internal/config"
	"github.com/nnurry/gopds/probabilistics/internal/database/postgres"
	"github.com/nnurry/gopds/probabilistics/pkg/models/decayable"
	concretefilter "github.com/nnurry/gopds/probabilistics/pkg/models/filter/concrete"
	concretemeta "github.com/nnurry/gopds/probabilistics/pkg/models/meta/concrete"
)

func setFilter(body *request_schema.FilterCreateBody, pw *decayable.Filter) {
	switch body.Filter.Type {
	case "STANDARD_BLOOM":
		core := concretefilter.NewStandardBF(
			body.Filter.MaxCardinality,
			body.Filter.ErrorRate,
			body.Meta.Key,
		)
		pw.SetCore(core)
	case "REDIS_BLOOM":
		core := concretefilter.NewRedisBF(
			body.Filter.MaxCardinality,
			body.Filter.ErrorRate,
			2,
			false,
			body.Meta.Key,
		)
		pw.SetCore(core)
	default:
		panic(errors.New("Not implemented this kind of filter: " + body.Filter.Type))
	}
}

func setFilterMeta(pw *decayable.Filter) {
	pw.SetMeta(concretemeta.NewDecayableMeta(
		config.ProbabilisticCfg.DecayDuration,
	))
}

func CreateFilter(body *request_schema.FilterCreateBody) *decayable.Filter {
	// create filter
	prob := &decayable.Filter{}
	setFilter(body, prob)

	// create meta
	setFilterMeta(prob)
	return prob
}

func SaveFilter(
	pw *decayable.Filter,
	isCreate bool,
	doCommit bool,
	refreshLastUsed bool,
	tx *sql.Tx) error {
	var err error
	var filterId uint

	if isCreate {
		err = postgres.Client.QueryRow(`
		INSERT INTO filters (type, key, max_cardinality, max_fp, hash_func_num, hash_func_type)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (type, key, max_cardinality, max_fp, hash_func_type) DO NOTHING
		RETURNING id
		`,
			pw.Core().Meta().FilterType(),
			pw.Core().Meta().Key(),
			pw.Core().Meta().MaxCard(),
			pw.Core().Meta().MaxFp(),
			pw.Core().Meta().HashFuncNum(),
			pw.Core().Meta().HashFuncType(),
		).Scan(&filterId)

		if err == sql.ErrNoRows {
			err = postgres.Client.QueryRow(`
			SELECT id 
			FROM filters
			WHERE type = $1
			AND key = $2
			AND max_cardinality = $3
			AND max_fp = $4
			AND hash_func_type = $5
			LIMIT 1
			`,
				pw.Core().Meta().FilterType(),
				pw.Core().Meta().Key(),
				pw.Core().Meta().MaxCard(),
				pw.Core().Meta().MaxFp(),
				pw.Core().Meta().HashFuncType(),
			).Scan(&filterId)
		}

		if err != nil {
			tx.Rollback()
			return err
		}

		pw.Core().Meta().SetId(filterId)
	}

	_, err = postgres.Client.Exec(`
	INSERT INTO filter_blob (filter_id, filter_byte)
	VALUES ($1, $2)
	ON CONFLICT (filter_id) DO UPDATE
	SET filter_byte = EXCLUDED.filter_byte;
	`,
		pw.Core().Meta().Id(),
		pw.Core().Serialize(),
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	if doCommit {
		tx.Commit()
	}

	if refreshLastUsed {
		pw.Meta().ResetLastUsed()
	}

	return nil

}
