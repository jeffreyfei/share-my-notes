package server

import "github.com/jinzhu/gorm"
import "github.com/gorilla/mux"

type Server struct {
	db     *gorm.DB
	Router *mux.Router
}

func NewServer(db *gorm.DB) *Server {
	server := Server{}
	server.db = db
	server.Router = BuildRouter()
	return &server
}
