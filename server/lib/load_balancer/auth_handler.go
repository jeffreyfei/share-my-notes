package load_balancer

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (lb *LoadBalancer) googleAuthHandler(w http.ResponseWriter, r *http.Request) {
	action := mux.Vars(r)["action"]
	done := make(chan struct{})
	req := new(request)
	req.w = w
	req.r = r
	req.route = fmt.Sprintf("/auth/google/%s", action)
	go lb.forwardSyncRequest(req, done)
	<-done
	w.WriteHeader(http.StatusOK)
}
