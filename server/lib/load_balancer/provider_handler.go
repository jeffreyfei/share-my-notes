package load_balancer

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (lb *LoadBalancer) providerRegisterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	url := r.PostFormValue("url")
	lb.Providers[url] = 0
	log.WithField("url", url).Info("New provider registered.")
	w.WriteHeader(http.StatusOK)
}
