package server

import (
	"encoding/gob"
	"fmt"

	"github.com/jeffreyfei/share-my-notes/server/lib/buffer"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jeffreyfei/share-my-notes/server/lib/user"
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
	buffer       *buffer.Buffer
	lbURL        string
}

func NewServer(db *gorm.DB, baseURL, sessionKey, clientID, clientSecret, lbURL string) *Server {
	server := Server{}
	server.db = db
	server.Router = server.buildRouter()
	server.baseURL = baseURL
	server.lbURL = lbURL
	server.oauth2Config = getOauthConfig(baseURL, clientID, clientSecret)
	server.sessionStore = getSessionStore(sessionKey)
	server.buffer = buffer.NewBuffer(5000, 100)
	gob.Register(user.UserModel{})
	gob.Register(oauth2.Token{})
	return &server
}

func (s *Server) StartBufferProc() {
	go s.buffer.StartProc()
}

func getSessionStore(sessionKey string) *sessions.CookieStore {
	cookieStore := sessions.NewCookieStore([]byte(sessionKey))
	return cookieStore
}

func getOauthConfig(baseURL, clientID, clientSecret string) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  fmt.Sprintf("%s/auth/google/callback", baseURL),
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}
}
