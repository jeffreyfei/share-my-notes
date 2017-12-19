package load_balancer

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Handles client requests for markdown compilation
func (lb *LoadBalancer) mdClientHandler(w http.ResponseWriter, r *http.Request) {
	action := mux.Vars(r)["action"]
	switch action {
	case "create":
	case "get":
	case "update", "delete":
	default:
		log.Error("Invalid action")
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}
