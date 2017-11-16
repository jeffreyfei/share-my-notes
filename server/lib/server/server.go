package server

import "github.com/jinzhu/gorm"
import "github.com/gorilla/mux"

type Server struct {
	db     *gorm.DB
	router *mux.Router
}

func NewServer(db *gorm.DB, router *muxRouter) *Server {
	server := Server{}
	server.db = db
	server.router = BuildRouter()
	return &server
}
