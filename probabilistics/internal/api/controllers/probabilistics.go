package controllers

import (
	"encoding/json"
	"fmt"
	"gopds/probabilistics/internal/database/postgres"
	"gopds/probabilistics/internal/service"
	"gopds/probabilistics/pkg/models/wrapper"
	"io"
	"net/http"
)

func probCreate(w http.ResponseWriter, r *http.Request) {
	GetBasicInfo(r)
	var err error
	js, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	body := &service.ProbCreateBody{}
	err = json.Unmarshal(js, body)
	if err != nil {
		panic(err)
	}

	tx, _ := postgres.Client.Begin()

	p := service.CreateProbabilistic(body)
	err = service.SaveProbabilistic(p, true, true, tx)

	wrapper.GetWrapper().Add(wrapper.WrapperKey(*body), p)

	if err != nil {
		panic(fmt.Sprint("Can't save probabilistic for some reason", err))
	}

	w.Write([]byte(fmt.Sprint("Created probabilistic", p)))
}

func probExists(w http.ResponseWriter, r *http.Request) {
	GetBasicInfo(r)
}

func probAdd(w http.ResponseWriter, r *http.Request) {
	GetBasicInfo(r)
}

func probCard(w http.ResponseWriter, r *http.Request) {
	GetBasicInfo(r)
}

type probabilistic struct {
	Create func(http.ResponseWriter, *http.Request)
	Exists func(http.ResponseWriter, *http.Request)
	Add    func(http.ResponseWriter, *http.Request)
	Card   func(http.ResponseWriter, *http.Request)
}

var Probabilistics = probabilistic{
	Create: probCreate,
	Exists: probExists,
	Add:    probAdd,
	Card:   probCard,
}
