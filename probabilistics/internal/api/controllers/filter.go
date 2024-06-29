package controllers

import (
	"fmt"
	"gopds/probabilistics/internal/database/postgres"
	"gopds/probabilistics/internal/service"
	"gopds/probabilistics/pkg/models/wrapper"
	"net/http"
)

func filterCreate(w http.ResponseWriter, r *http.Request) {
	GetBasicInfo(r)
	var err error
	body := &service.FilterCreateBody{}
	loadJson(r.Body, body)

	tx, _ := postgres.Client.Begin()

	pw := service.CreateFilter(body)
	err = service.SaveFilter(pw, true, true, tx)

	wrapperKey := &wrapper.FilterKey{
		Type:           pw.Core().Meta().FilterType(),
		Key:            pw.Core().Meta().Key(),
		MaxCardinality: pw.Core().Meta().MaxCard(),
		ErrorRate:      pw.Core().Meta().MaxFp(),
	}

	wrapper.GetWrapper().AddFilter(*wrapperKey, pw)

	if err != nil {
		panic(fmt.Sprint("Can't save filter for some reason", err))
	}

	w.Write([]byte(fmt.Sprint("Created filter", pw)))
}

func filterExists(w http.ResponseWriter, r *http.Request) {
	GetBasicInfo(r)
}

func filterAdd(w http.ResponseWriter, r *http.Request) {
	GetBasicInfo(r)
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
