package service

import (
	"database/sql"
	"errors"
	"gopds/probabilistics/internal/config"
	"gopds/probabilistics/internal/database/postgres"
	"gopds/probabilistics/pkg/models/decayable"
	concretefilter "gopds/probabilistics/pkg/models/filter/concrete"
	concretemeta "gopds/probabilistics/pkg/models/meta/concrete"
)

type FilterBody struct {
	Type           string  `json:"type"`
	MaxCardinality uint    `json:"max_cardinality"`
	ErrorRate      float64 `json:"error_rate"`
}

type FilterCreateBody struct {
	Meta   MetaBody   `json:"meta"`
	Filter FilterBody `json:"filter"`
}

type FilterExistsBody struct {
	Meta   PairMetaBody `json:"meta"`
	Filter FilterBody   `json:"filter"`
}

type FilterAddBody struct {
	Meta   PairMetaBody `json:"meta"`
	Filter FilterBody   `json:"filter"`
}

func setFilter(body *FilterCreateBody, pw *decayable.Filter) {
	switch body.Filter.Type {
	case "standard_bloom":
		core := concretefilter.NewStandardBF(
			body.Filter.MaxCardinality,
			body.Filter.ErrorRate,
			body.Meta.Key,
		)
		pw.SetCore(core)
	case "redis_bloom":
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
	pw.SetMeta(concretemeta.NewProbabilisticMeta(
		config.ProbabilisticCfg.DecayDuration,
	))
}

func CreateFilter(body *FilterCreateBody) *decayable.Filter {
	// create filter
	prob := &decayable.Filter{}
	setFilter(body, prob)

	// create meta
	setFilterMeta(prob)
	return prob
}

func SaveFilter(pw *decayable.Filter, isCreate bool, doCommit bool, tx *sql.Tx) error {
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

	pw.Meta().ResetLastUsed()

	return nil

}
