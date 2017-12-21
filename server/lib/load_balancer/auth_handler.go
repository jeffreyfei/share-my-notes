package load_balancer

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Handles google auth requests from client
// Calls forwardRedirectRequest in forwarder
func (lb *LoadBalancer) googleAuthHandler(w http.ResponseWriter, r *http.Request) {
	action := mux.Vars(r)["action"]
	done := make(chan string)
	req := new(request)
	req.w = w
	req.r = r
	req.route = fmt.Sprintf("/auth/google/%s", action)
	go lb.forwardRedirectRequest(req, done)
	url := <-done
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
