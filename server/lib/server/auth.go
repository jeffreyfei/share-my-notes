package server

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"

	"github.com/satori/go.uuid"
	"google.golang.org/api/plus/v1"

	log "github.com/sirupsen/logrus"
)

const (
	defaultSessionID        = "default"
	redirectKey             = "redirect"
	stateKey                = "state"
	codeKey                 = "code"
	googleProfileSessionKey = "google_profile"
	oauthSessionKey         = "oauth_token"
)

type Profile struct {
	ID,
	DisplayName,
	ImageURL string
}

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

func (s *Server) fetchGooglePlusProfile(ctx context.Context, token *oauth2.Token) (*plus.Person, error) {
	client := oauth2.NewClient(ctx, s.oauth2Config.TokenSource(ctx, token))
	plusService, err := plus.New(client)
	if err != nil {
		return nil, err
	}
	return plusService.People.Get("me").Do()
}

func (s *Server) getProfileFromSession(r *http.Request) *Profile {
	session, err := s.sessionStore.Get(r, defaultSessionID)
	if err != nil {
		return nil
	}
	token, ok := session.Values[oauthSessionKey].(oauth2.Token)
	if !ok || !token.Valid() {
		return nil
	}
	profile, ok := session.Values[googleProfileSessionKey].(Profile)
	if !ok {
		return nil
	}
	return &profile
}

func formatProfile(profile *plus.Person) *Profile {
	return &Profile{
		ID:          profile.Id,
		DisplayName: profile.DisplayName,
		ImageURL:    profile.Image.Url,
	}
}

func (s *Server) googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := uuid.NewV4().String()
	oauthFlowSession, err := s.sessionStore.New(r, sessionID)
	if err != nil {
		log.WithField("err", err).Error("Could not create oauth session")
		return
	}
	redirectURL, err := validateRedirectURL(r.FormValue(redirectKey))
	if err != nil {
		log.WithField("err", err).Error("Invalid URL")
		return
	}
	oauthFlowSession.Values[redirectKey] = redirectURL
	if err := oauthFlowSession.Save(r, w); err != nil {
		log.WithField("err", err).Error("Could not save session")
		return
	}
	url := s.oauth2Config.AuthCodeURL(sessionID, oauth2.ApprovalForce, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (s *Server) googleLoginCallbackHandler(w http.ResponseWriter, r *http.Request) {
	oauthFlowSession, err := s.sessionStore.Get(r, r.FormValue(stateKey))
	if err != nil {
		log.WithField("err", err).Error("Invalid state parameter")
		return
	}
	redirectURL, ok := oauthFlowSession.Values[redirectKey].(string)
	if !ok {
		log.WithField("err", err).Error("Invalid state parameter")
		return
	}
	token, err := s.oauth2Config.Exchange(context.Background(), r.FormValue(codeKey))
	if err != nil {
		log.WithField("err", err).Error("Could not get auth token")
		return
	}
	session, err := s.sessionStore.New(r, defaultSessionID)
	if err != nil {
		log.WithField("err", err).Error("Could not get default session")
	}
	ctx := context.Background()
	profile, err := s.fetchGooglePlusProfile(ctx, token)
	if err != nil {
		log.WithField("err", err).Error("Could not get Google profile")
		return
	}
	session.Values[oauthSessionKey] = token
	session.Values[googleProfileSessionKey] = formatProfile(profile)
	if err := session.Save(r, w); err != nil {
		log.WithField("err", err).Error("Could not save session")
		return
	}
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (s *Server) googleLogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := s.sessionStore.New(r, defaultSessionID)
	if err != nil {
		log.WithField("err", err).Error("Could not get default session")
		return
	}
	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		log.WithField("err", err).Error("Could not save session")
		return
	}
	redirectURL, err := validateRedirectURL(r.FormValue(redirectKey))
	if err != nil {
		log.WithField("err", err).Error("Invalid URL")
	}
	http.Redirect(w, r, redirectURL, http.StatusFound)
}
