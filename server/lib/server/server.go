package server

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Server struct {
	db           *gorm.DB
	Router       *mux.Router
	baseURL      string
	oauth2Config *oauth2.Config
	sessionStore *sessions.CookieStore
}

func NewServer(db *gorm.DB, baseURL, clientID, clientSecret string) *Server {
	server := Server{}
	server.db = db
	server.Router = buildRouter()
	server.baseURL = baseURL
	server.oauth2Config = server.getOauthConfig(clientID, clientSecret)
	return &server
}

func (s *Server) getOauthConfig(clientID, clientSecret string) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  fmt.Sprintf("%s/auth/google/callback", s.baseURL),
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}
}
