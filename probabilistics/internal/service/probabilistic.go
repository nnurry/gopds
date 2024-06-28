package service

import (
	"database/sql"
	"errors"
	"gopds/probabilistics/internal/config"
	"gopds/probabilistics/internal/database/postgres"
	concretecardinal "gopds/probabilistics/pkg/models/cardinal/concrete"
	concretefilter "gopds/probabilistics/pkg/models/filter/concrete"
	concretemeta "gopds/probabilistics/pkg/models/meta/concrete"
	"gopds/probabilistics/pkg/models/probabilistic"
)

type FilterBody struct {
	Type           string  `json:"type"`
	MaxCardinality uint    `json:"max_cardinality"`
	ErrorRate      float64 `json:"error_rate"`
}

type CardinalBody struct {
	Type string `json:"type"`
}

type MetaBody struct {
	Key string `json:"key"`
}

type PairMetaBody struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ProbCreateBody struct {
	Meta     MetaBody     `json:"meta"`
	Filter   FilterBody   `json:"filter"`
	Cardinal CardinalBody `json:"cardinal"`
}

type ProbExistsBody struct {
	Meta   PairMetaBody `json:"meta"`
	Filter FilterBody   `json:"filter"`
}

type ProbCardBody struct {
	Meta     MetaBody     `json:"meta"`
	Cardinal CardinalBody `json:"cardinal"`
}

type ProbAddBody struct {
	Meta     PairMetaBody `json:"meta"`
	Cardinal CardinalBody `json:"cardinal"`
	Filter   FilterBody   `json:"filter"`
}

func setFilter(body *ProbCreateBody, prob *probabilistic.Probabilistic) {
	switch body.Filter.Type {
	case "standard_bloom":
		prob.SetFilter(concretefilter.NewStandardBF(
			body.Filter.MaxCardinality,
			body.Filter.ErrorRate,
		))
	case "redis_bloom":
		prob.SetFilter(concretefilter.NewRedisBF(
			10,
			0.01,
			2,
			false,
			body.Meta.Key,
		))
	default:
		panic(errors.New("Not implemented this kind of filter: " + body.Filter.Type))
	}
}

func setCardinal(body *ProbCreateBody, prob *probabilistic.Probabilistic) {
	switch body.Cardinal.Type {
	case "standard_hll":
		prob.SetCardinal(concretecardinal.NewStandardHLL(false, 14, body.Meta.Key))
	case "redis_hll":
		prob.SetCardinal(concretecardinal.NewRedisHLL(body.Meta.Key))
	default:
		panic(errors.New("Not implemented this kind of cardinal: " + body.Filter.Type))
	}
}

func setMeta(body *ProbCreateBody, prob *probabilistic.Probabilistic) {
	prob.SetMeta(concretemeta.NewProbabilisticMeta(
		body.Meta.Key,
		config.ProbabilisticCfg.DecayDuration,
	))
}

func CreateProbabilistic(body *ProbCreateBody) *probabilistic.Probabilistic {
	// create filter
	prob := &probabilistic.Probabilistic{}
	setFilter(body, prob)

	// create cardinal
	setCardinal(body, prob)

	// create meta
	setMeta(body, prob)
	return prob
}

func SaveProbabilistic(
	p *probabilistic.Probabilistic,
	isCreate bool,
	doCommit bool,
	tx *sql.Tx) error {

	var err error
	client := postgres.Client

	var filterId uint
	var cardinalId uint
	if isCreate {
		err = client.QueryRow(`
		INSERT INTO filters (type, max_cardinality, max_fp, hash_func_num, hash_func_type)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (type, max_cardinality, max_fp, hash_func_type) DO NOTHING
		RETURNING id
		`,
			p.Filter().Meta().FilterType(),
			p.Filter().Meta().MaxCard(),
			p.Filter().Meta().MaxFp(),
			p.Filter().Meta().HashFuncNum(),
			p.Filter().Meta().HashFuncType(),
		).Scan(&filterId)

		if err != nil {
			tx.Rollback()
			return err
		}

		p.Filter().Meta().SetId(filterId)

		err = client.QueryRow(`
		INSERT INTO cardinals (type)
		VALUES ($1)
		ON CONFLICT (type) DO NOTHING
		RETURNING id
		`,
			p.Cardinal().Meta().CardinalType(),
		).Scan(&cardinalId)

		if err != nil {
			tx.Rollback()
			return err
		}

		p.Cardinal().Meta().SetId(filterId)
	}

	_, err = client.Exec(`
	INSERT INTO wrappers (filter_id, cardinal_id, key, filter_byte, cardinal_byte)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (filter_id, cardinal_id) DO UPDATE
	SET 
		filter_byte = EXCLUDED.filter_byte,
		cardinal_byte = EXCLUDED.cardinal_byte;
	`,
		p.Filter().Meta().Id(),
		p.Cardinal().Meta().Id(),
		p.Meta().Key(),
		p.Filter().Serialize(),
		p.Cardinal().Serialize(),
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	if doCommit {
		tx.Commit()
	}

	p.Meta().ResetLastUsed()

	return nil
}
