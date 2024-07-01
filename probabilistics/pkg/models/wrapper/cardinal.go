package wrapper

import (
	"database/sql"
	"log"

	"github.com/nnurry/gopds/probabilistics/internal/config"
	"github.com/nnurry/gopds/probabilistics/internal/database/postgres"
	concretecardinal "github.com/nnurry/gopds/probabilistics/pkg/models/cardinal/concrete"
	"github.com/nnurry/gopds/probabilistics/pkg/models/decayable"
	concretemeta "github.com/nnurry/gopds/probabilistics/pkg/models/meta/concrete"
)

type CardinalKey struct {
	Type string
	Key  string
}

type CardinalWrapper struct {
	core    map[CardinalKey]*decayable.Cardinal
	counter uint
}

func NewCardinalWrapper() *CardinalWrapper {
	return &CardinalWrapper{
		core:    make(map[CardinalKey]*decayable.Cardinal),
		counter: 0,
	}
}

func (pw *CardinalWrapper) GetCardinal(k CardinalKey, inMemOnly bool) *decayable.Cardinal {
	var v *decayable.Cardinal
	v, exists := pw.core[k]
	if inMemOnly || exists {
		return v
	}
	v = pw.FetchCardinal(k)
	if v != nil {
		pw.Add(k, v)
	}
	return v
}

func (pw *CardinalWrapper) FetchCardinal(k CardinalKey) *decayable.Cardinal {
	cardinal := &decayable.Cardinal{}

	switch k.Type {
	case "STANDARD_HLL":
		cardinal.SetCore(concretecardinal.NewStandardHLL(false, 14, k.Key))
	case "REDIS_HLL":
		cardinal.SetCore(concretecardinal.NewRedisHLL(k.Key))
	}

	cardinal.SetMeta(concretemeta.NewDecayableMeta(config.ProbabilisticCfg.DecayDuration))
	queryStruct := struct {
		Id   uint
		Blob []byte
	}{}

	err := postgres.Client.QueryRow(`
	SELECT id, cardinal_byte
	FROM (SELECT id, key, type FROM cardinals WHERE key = $1 AND type = $2) cardinals
	JOIN cardinal_blob
	ON cardinals.id = cardinal_blob.cardinal_id
	LIMIT 1;
	`, k.Key, k.Type).Scan(&queryStruct.Id, &queryStruct.Blob)

	if err == sql.ErrNoRows {
		return nil
	}

	if err != nil {
		log.Println("Can't fetch cardinals from db: ", err)
		panic(err)
	}

	cardinal.Core().Meta().SetId(queryStruct.Id)
	cardinal.Core().Deserialize(queryStruct.Blob)

	log.Println(
		"Fetched cardinal:",
		cardinal.Core().Meta().CardinalType(),
		cardinal.Core().Meta().Key(),
		cardinal.Core().Meta().Id(),
		cardinal.Core().Cardinality(),
	)
	return cardinal

}

func (pw *CardinalWrapper) Core() map[CardinalKey]*decayable.Cardinal {
	return pw.core
}

func (pw *CardinalWrapper) Counter() uint {
	return pw.counter
}

func (pw *CardinalWrapper) Add(k CardinalKey, v *decayable.Cardinal) {
	_, exists := pw.core[k]
	pw.core[k] = v
	if !exists {
		pw.counter++
	}
}

func (pw *CardinalWrapper) Delete(k CardinalKey) {
	_, exists := pw.core[k]
	delete(pw.core, k)
	if exists {
		pw.counter--
	}
}
