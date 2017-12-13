package load_balancer

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

func (l *LoadBalancer) buildClientRoutes() routes {
	return routes{}
}

func (l *LoadBalancer) buildProviderRoutes() routes {
	return routes{}
}

func (l *LoadBalancer) buildRouter(r routes) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range r {
		router.HandleFunc(route.Route, route.Handler).Methods(route.Type)
	}
	return router
}
