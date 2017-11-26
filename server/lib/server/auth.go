package server

import (
	"errors"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"

	"github.com/satori/go.uuid"

	log "github.com/sirupsen/logrus"
)

func validateRedirectURL(path string) (string, error) {
	if path == "" {
		return "/", nil
	}

	parsedURL, err := url.Parse(path)
	if err != nil {
		return "/", err
	}
	if parsedURL.IsAbs() {
		return "/", errors.New("URL must not be absolute")
	}
	return path, nil
}

func (s *Server) googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := uuid.NewV4().String()
	oauthFlowSession, err := s.sessionStore.New(r, sessionID)
	if err != nil {
		log.WithField("err", err).Error("Could not create oauth session")
	}
	redirectURL, err := validateRedirectURL(r.FormValue("redirect"))
	if err != nil {
		log.WithField("err", err).Error("Invalid URL")
	}
	oauthFlowSession.Values["redirect"] = redirectURL
	if err := oauthFlowSession.Save(r, w); err != nil {
		log.WithField("err", err).Error("Could not save session")
	}
	url := s.oauth2Config.AuthCodeURL(sessionID, oauth2.ApprovalForce, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (s *Server) googleLoginCallbackHandler(w http.ResponseWriter, r *http.Request) {

}
