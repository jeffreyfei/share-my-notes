package server

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/jeffreyfei/share-my-notes/server/lib/buffer"
	"github.com/jeffreyfei/share-my-notes/server/lib/router"
	log "github.com/sirupsen/logrus"

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
	lbPublicURL  string
	lbPrivateURL string
}

// Initializes a new server instance
func NewServer(db *gorm.DB, baseURL, sessionKey, clientID, clientSecret, lbPrivateURL, lbPublicURL string) *Server {
	server := Server{}
	server.db = db
	server.Router = router.BuildRouter(server.buildRoutes())
	server.baseURL = baseURL
	server.lbPublicURL = lbPublicURL
	server.lbPrivateURL = lbPrivateURL
	server.oauth2Config = getOauthConfig(lbPublicURL, clientID, clientSecret)
	server.sessionStore = getSessionStore(sessionKey)
	server.buffer = buffer.NewBuffer(5000, 100)
	gob.Register(user.UserModel{})
	gob.Register(oauth2.Token{})
	return &server
}

// Start processing jobs on the buffer
func (s *Server) StartBufferProc() {
	go s.buffer.StartProc()
}

// Sends a registration request to the load balancer
func (s *Server) RegisterLoadBalancer() {
	payload := url.Values{}
	payload.Add("url", s.baseURL)
	route := fmt.Sprintf("%s/provider/register", s.lbPrivateURL)
	if res, err := http.PostForm(route, payload); err != nil || res.StatusCode != http.StatusOK {
		log.WithField("err", err).Error("Failed to register to load balancer. Trying again in 1s.")
		time.Sleep(1000 * time.Millisecond)
		go s.RegisterLoadBalancer()
	} else {
		log.Info("Successfully registered to load balancer")
	}
}

// Returns a list of available routes
func (s *Server) buildRoutes() router.Routes {
	return router.Routes{
		router.Route{
			"GET",
			"/report/status",
			s.reportStatusHandler,
		},
		router.Route{
			"GET",
			"/auth/google/login",
			s.googleLoginHandler,
		},
		router.Route{
			"GET",
			"/auth/google/callback",
			s.googleLoginCallbackHandler,
		},
		router.Route{
			"GET",
			"/auth/google/logout",
			s.googleLogoutHandler,
		},
		router.Route{
			"GET",
			"/note/md/get/{id}",
			s.mdGetHandler,
		},
		router.Route{
			"POST",
			"/note/md/update/{id}",
			s.mdUpdateHandler,
		},
		router.Route{
			"POST",
			"/note/md/delete/{id}",
			s.mdDeleteHandler,
		},
		router.Route{
			"POST",
			"/note/md/create",
			s.mdCreateHandler,
		},
	}
}

// Create a new session store instance with the given sessionKey
func getSessionStore(sessionKey string) *sessions.CookieStore {
	cookieStore := sessions.NewCookieStore([]byte(sessionKey))
	return cookieStore
}

// Returns the oauth configuration for Google
func getOauthConfig(baseURL, clientID, clientSecret string) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  fmt.Sprintf("%s/auth/google/callback", baseURL),
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}
}
