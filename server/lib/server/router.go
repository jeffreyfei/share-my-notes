package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type route struct {
	Type    string
	Route   string
	Handler http.HandlerFunc
}

type routes []route

func buildRoutes() routes {
	return routes{}
}

func buildRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range buildRoutes() {
		router.HandleFunc(route.Route, route.Handler).Methods(route.Type)
	}
	return router
}
