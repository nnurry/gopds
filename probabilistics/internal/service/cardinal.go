package service

import (
	"database/sql"
	"errors"
	request_schema "gopds/probabilistics/internal/api/schemas/request"
	"gopds/probabilistics/internal/config"
	"gopds/probabilistics/internal/database/postgres"
	concretecardinal "gopds/probabilistics/pkg/models/cardinal/concrete"
	"gopds/probabilistics/pkg/models/decayable"
	concretemeta "gopds/probabilistics/pkg/models/meta/concrete"
)

func setCardinal(body *request_schema.CardinalCreateBody, pw *decayable.Cardinal) {
	switch body.Cardinal.Type {
	case "standard_hll":
		core := concretecardinal.NewStandardHLL(false, 14, body.Meta.Key)
		pw.SetCore(core)
	case "redis_hll":
		core := concretecardinal.NewRedisHLL(body.Meta.Key)
		pw.SetCore(core)
	default:
		panic(errors.New("Not implemented this kind of cardinal: " + body.Cardinal.Type))
	}
}

func setCardinalMeta(pw *decayable.Cardinal) {
	pw.SetMeta(concretemeta.NewDecayableMeta(
		config.ProbabilisticCfg.DecayDuration,
	))
}

func CreateCardinal(body *request_schema.CardinalCreateBody) *decayable.Cardinal {
	// create cardinal
	prob := &decayable.Cardinal{}
	setCardinal(body, prob)

	// create meta
	setCardinalMeta(prob)
	return prob
}

func SaveCardinal(
	pw *decayable.Cardinal,
	isCreate bool,
	doCommit bool,
	refreshLastUsed bool,
	tx *sql.Tx) error {
	var err error
	var cardinalId uint

	if isCreate {
		err = postgres.Client.QueryRow(`
		INSERT INTO cardinals (type, key)
		VALUES ($1, $2)
		ON CONFLICT (type, key) DO NOTHING
		RETURNING id
		`,
			pw.Core().Meta().CardinalType(),
			pw.Core().Meta().Key(),
		).Scan(&cardinalId)

		if err == sql.ErrNoRows {
			err = postgres.Client.QueryRow(`
			SELECT id 
			FROM cardinals
			WHERE type = $1
			AND key = $2
			LIMIT 1
			`,
				pw.Core().Meta().CardinalType(),
				pw.Core().Meta().Key(),
			).Scan(&cardinalId)
		}

		if err != nil {
			tx.Rollback()
			return err
		}

		pw.Core().Meta().SetId(cardinalId)
	}

	_, err = postgres.Client.Exec(`
	INSERT INTO cardinal_blob (cardinal_id, cardinal_byte)
	VALUES ($1, $2)
	ON CONFLICT (cardinal_id) DO UPDATE
	SET cardinal_byte = EXCLUDED.cardinal_byte;
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
