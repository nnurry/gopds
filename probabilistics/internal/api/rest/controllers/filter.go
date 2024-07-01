package controllers

import (
	"fmt"
	"log"
	"net/http"

	request_schema "github.com/nnurry/gopds/probabilistics/internal/api/rest/schemas/request"
	"github.com/nnurry/gopds/probabilistics/internal/database/postgres"
	"github.com/nnurry/gopds/probabilistics/internal/service"
	"github.com/nnurry/gopds/probabilistics/pkg/models/wrapper"
)

func filterCreate(w http.ResponseWriter, r *http.Request) {
	log.Println(GetBasicInfo(r))
	var err error
	body := &request_schema.FilterCreateBody{}
	loadJson(r.Body, body)

	tx, _ := postgres.Client.Begin()

	pw := service.CreateFilter(body)
	err = service.SaveFilter(pw, true, true, true, tx)

	if err != nil {
		http.Error(w, "Failed to save filter: "+err.Error(), http.StatusInternalServerError)
	}

	wrapperKey := &wrapper.FilterKey{
		Type:           pw.Core().Meta().FilterType(),
		Key:            pw.Core().Meta().Key(),
		MaxCardinality: pw.Core().Meta().MaxCard(),
		ErrorRate:      pw.Core().Meta().MaxFp(),
	}

	wrapper.GetWrapper().AddFilter(*wrapperKey, pw)

	w.Write([]byte(fmt.Sprint("Created filter", pw)))
}

func filterExists(w http.ResponseWriter, r *http.Request) {
	log.Println(GetBasicInfo(r))

	body := &request_schema.FilterExistsBody{}
	loadJson(r.Body, body)

	filterKey := wrapper.FilterKey{
		Type:           body.Filter.Type,
		Key:            body.Meta.Key,
		MaxCardinality: body.Filter.MaxCardinality,
		ErrorRate:      body.Filter.ErrorRate,
	}

	filter := wrapper.GetWrapper().FilterWrapper().GetFilter(filterKey, false)

	var exists bool
	if filter != nil {
		exists = filter.Core().Exists([]byte(body.Meta.Value))
	}

	w.Write([]byte(body.Meta.Value + " exists in " + body.Meta.Key + ": " + fmt.Sprint(exists)))
}

func filterAdd(w http.ResponseWriter, r *http.Request) {
	log.Println(GetBasicInfo(r))
	body := &request_schema.FilterAddBody{}
	loadJson(r.Body, body)

	filterKey := wrapper.FilterKey{
		Type:           body.Filter.Type,
		Key:            body.Meta.Key,
		MaxCardinality: body.Filter.MaxCardinality,
		ErrorRate:      body.Filter.ErrorRate,
	}

	filter := wrapper.GetWrapper().FilterWrapper().GetFilter(filterKey, false)
	if filter != nil {
		filter.Core().AddString(body.Meta.Value)
	}

	w.Write([]byte("Added " + body.Meta.Value + " into " + body.Meta.Key))
}

type AbstractFilter struct {
	Create func(http.ResponseWriter, *http.Request)
	Exists func(http.ResponseWriter, *http.Request)
	Add    func(http.ResponseWriter, *http.Request)
}

var Filter = AbstractFilter{
	Create: filterCreate,
	Exists: filterExists,
	Add:    filterAdd,
}
