package controllers

import (
	"fmt"
	"gopds/probabilistics/internal/database/postgres"
	"gopds/probabilistics/internal/service"
	"gopds/probabilistics/pkg/models/wrapper"
	"net/http"
)

func cardinalCreate(w http.ResponseWriter, r *http.Request) {
	GetBasicInfo(r)
	var err error
	body := &service.CardinalCreateBody{}
	loadJson(r.Body, body)

	tx, _ := postgres.Client.Begin()

	pw := service.CreateCardinal(body)
	err = service.SaveCardinal(pw, true, true, tx)

	wrapperKey := &wrapper.CardinalKey{
		Type: pw.Core().Meta().CardinalType(),
		Key:  pw.Core().Meta().Key(),
	}

	wrapper.GetWrapper().AddCardinal(*wrapperKey, pw)

	if err != nil {
		panic(fmt.Sprint("Can't save cardinal for some reason", err))
	}

	w.Write([]byte(fmt.Sprint("Created cardinal", pw)))
}

func cardinalAdd(w http.ResponseWriter, r *http.Request) {
	GetBasicInfo(r)
}

func cardinalCard(w http.ResponseWriter, r *http.Request) {
	GetBasicInfo(r)
}

type AbstractCardinal struct {
	Create func(http.ResponseWriter, *http.Request)
	Add    func(http.ResponseWriter, *http.Request)
	Card   func(http.ResponseWriter, *http.Request)
}

var Cardinal = AbstractCardinal{
	Create: cardinalCreate,
	Card:   cardinalCard,
	Add:    cardinalAdd,
}
