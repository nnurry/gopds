package rest

import (
	"net/http"

	"github.com/nnurry/gopds/probabilistics/internal/api/rest/controllers"
)

func SetupFilterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/filter/create", controllers.Filter.Create)
	mux.HandleFunc("/filter/exists", controllers.Filter.Exists)
	mux.HandleFunc("/filter/add", controllers.Filter.Add)
}

func SetupCardinalRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/cardinal/create", controllers.Cardinal.Create)
	mux.HandleFunc("/cardinal/card", controllers.Cardinal.Card)
	mux.HandleFunc("/cardinal/add", controllers.Cardinal.Add)
}

func SetupMux() *http.ServeMux {
	var mux = http.NewServeMux()
	SetupFilterRoutes(mux)
	SetupCardinalRoutes(mux)
	return mux
}
