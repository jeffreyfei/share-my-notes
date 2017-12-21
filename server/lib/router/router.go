package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Type    string
	Route   string
	Handler http.HandlerFunc
}

type Routes []Route

// Build router from the list of routes provided
func BuildRouter(routes Routes) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.HandleFunc(route.Route, route.Handler).Methods(route.Type)
	}
	return router
}
