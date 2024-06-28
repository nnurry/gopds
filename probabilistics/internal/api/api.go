package api

import (
	"gopds/probabilistics/internal/api/controllers"
	"net/http"
)

func SetupMux() *http.ServeMux {
	var mux = http.NewServeMux()
	mux.HandleFunc("/probabilistic/create", controllers.Probabilistics.Create)
	mux.HandleFunc("/probabilistic/exists", controllers.Probabilistics.Exists)
	mux.HandleFunc("/probabilistic/add", controllers.Probabilistics.Add)
	mux.HandleFunc("/probabilistic/card", controllers.Probabilistics.Card)
	return mux
}
