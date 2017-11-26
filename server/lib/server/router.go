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

func buildRoutes(s *Server) routes {
	return routes{
		route{
			"google-login",
			"auth/google",
			s.googleLoginHandler,
		},
		route{
			"google-login-callback",
			"auth/google/callback",
			s.googleLoginCallbackHandler,
		},
	}
}

func (s *Server) buildRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range buildRoutes(s) {
		router.HandleFunc(route.Route, route.Handler).Methods(route.Type)
	}
	return router
}
