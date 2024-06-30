package controllers

import (
	"fmt"
	"gopds/probabilistics/internal/database/postgres"
	"gopds/probabilistics/internal/service"
	"gopds/probabilistics/pkg/models/wrapper"
	"net/http"
)

func cardinalCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Println(GetBasicInfo(r))
	var err error
	body := &service.CardinalCreateBody{}
	loadJson(r.Body, body)

	tx, _ := postgres.Client.Begin()

	pw := service.CreateCardinal(body)
	err = service.SaveCardinal(pw, true, true, true, tx)

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
	fmt.Println(GetBasicInfo(r))
	body := &service.CardinalAddBody{}
	loadJson(r.Body, body)

	cardinalKey := wrapper.CardinalKey{
		Type: body.Cardinal.Type,
		Key:  body.Meta.Key,
	}

	cardinal := wrapper.GetWrapper().CardinalWrapper().GetCardinal(cardinalKey, false)
	if cardinal != nil {
		cardinal.Core().AddString(body.Meta.Value)
	}

	w.Write([]byte("Added " + body.Meta.Value + " into " + body.Meta.Key))
}

func cardinalCard(w http.ResponseWriter, r *http.Request) {
	fmt.Println(GetBasicInfo(r))
	body := &service.CardinalCardBody{}
	loadJson(r.Body, body)

	cardinalKey := wrapper.CardinalKey{
		Type: body.Cardinal.Type,
		Key:  body.Meta.Key,
	}

	cardinal := wrapper.GetWrapper().CardinalWrapper().GetCardinal(cardinalKey, false)
	var cardinality uint64
	if cardinal != nil {
		cardinality = cardinal.Core().Cardinality()
	}
	w.Write([]byte("Cardinality of " + body.Meta.Key + " = " + fmt.Sprint(cardinality)))
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
