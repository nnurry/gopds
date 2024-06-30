package wrapper

import (
	"database/sql"
	"fmt"

	"github.com/nnurry/gopds/probabilistics/internal/config"
	"github.com/nnurry/gopds/probabilistics/internal/database/postgres"
	"github.com/nnurry/gopds/probabilistics/pkg/models/decayable"
	concretefilter "github.com/nnurry/gopds/probabilistics/pkg/models/filter/concrete"
	concretemeta "github.com/nnurry/gopds/probabilistics/pkg/models/meta/concrete"
)

type FilterKey struct {
	Type           string
	Key            string
	MaxCardinality uint
	ErrorRate      float64
}

type FilterWrapper struct {
	core    map[FilterKey]*decayable.Filter
	counter uint
}

func NewFilterWrapper() *FilterWrapper {
	return &FilterWrapper{
		core:    make(map[FilterKey]*decayable.Filter),
		counter: 0,
	}
}

func (pw *FilterWrapper) GetFilter(k FilterKey, inMemOnly bool) *decayable.Filter {
	var v *decayable.Filter
	v, exists := pw.core[k]
	if inMemOnly || exists {
		return v
	}
	v = pw.FetchFilter(k)
	if v != nil {
		pw.Add(k, v)
	}
	return v
}

func (pw *FilterWrapper) FetchFilter(k FilterKey) *decayable.Filter {
	filter := &decayable.Filter{}

	switch k.Type {
	case "standard_bloom":
		filter.SetCore(concretefilter.NewStandardBF(k.MaxCardinality, k.ErrorRate, k.Key))
	case "redis_bloom":
		filter.SetCore(concretefilter.NewRedisBF(k.MaxCardinality, k.ErrorRate, 2, false, k.Key))
	}

	filter.SetMeta(concretemeta.NewDecayableMeta(config.ProbabilisticCfg.DecayDuration))
	queryStruct := struct {
		Id   uint
		Blob []byte
	}{}

	err := postgres.Client.QueryRow(`
	SELECT id, filter_byte
	FROM (
		SELECT id, max_cardinality, max_fp, key, type
		FROM filters 
		WHERE max_cardinality = $1
		AND max_fp = $2
		AND key = $3
		AND type = $4
	) filters
	JOIN filter_blob
	ON filters.id = filter_blob.filter_id
	LIMIT 1;
	`,
		k.MaxCardinality, k.ErrorRate, k.Key, k.Type,
	).Scan(&queryStruct.Id, &queryStruct.Blob)

	if err == sql.ErrNoRows {
		return nil
	}

	if err != nil {
		fmt.Println("Can't fetch filters from db: ", err)
		panic(err)
	}

	filter.Core().Meta().SetId(queryStruct.Id)
	filter.Core().Deserialize(queryStruct.Blob)

	fmt.Println(
		"Fetched filter:",
		filter.Core().Meta().FilterType(),
		filter.Core().Meta().Key(),
		filter.Core().Meta().MaxCard(),
		filter.Core().Meta().MaxFp(),
	)
	return filter

}

func (pw *FilterWrapper) Core() map[FilterKey]*decayable.Filter {
	return pw.core
}

func (pw *FilterWrapper) Counter() uint {
	return pw.counter
}

func (pw *FilterWrapper) Add(k FilterKey, v *decayable.Filter) {
	_, exists := pw.core[k]
	pw.core[k] = v
	if !exists {
		pw.counter++
	}
}

func (pw *FilterWrapper) Delete(k FilterKey) {
	_, exists := pw.core[k]
	delete(pw.core, k)
	if exists {
		pw.counter--
	}
}
