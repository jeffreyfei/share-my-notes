package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Type    string
	Route   string
	Handler http.Handler
}

type Routes []Route

func buildRoutes() {
	return Routes{}
}

func BuildRouter() {
	router = mux.NewRouter().StrictSlash(true)
	for _, route := range buildRoutes() {
		router.HandleFunc(route.Route, route.Handler).Methods(route.Type)
	}
}
