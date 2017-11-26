package server

import (
	"net/http"
	"os"
)

func (s *Server) googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := s.oauth2Config.AuthCodeURL(os.Getenv("OAUTH_STATE_STRING"))
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (s *Server) googleLoginCallbackHandler(w http.ResponseWriter, r *http.Request) {

}
