package load_balancer

import (
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func copyReqHeader(inReq *http.Request, outReq *http.Request) {
	for k, v := range inReq.Header {
		values := ""
		for _, headerValue := range v {
			values += headerValue + " "
		}
		outReq.Header.Add(k, values)
	}
}

func copyResHeader(inRes *http.Response, outRes http.ResponseWriter) {
	for k, v := range inRes.Header {
		values := ""
		for _, headerValue := range v {
			values += headerValue + " "
		}
		outRes.Header().Add(k, values)
	}
}

type request struct {
	w     http.ResponseWriter
	r     *http.Request
	route string
}

// Forward sync tasks towards the providers (e.g. Login/Logout)
// Waits for a response from the provider before responding to the client
func (lb *LoadBalancer) forwardSyncRequest(req *request, done chan struct{}) {
	newReq, err := http.NewRequest(req.r.Method, fmt.Sprintf("%s/%s", lb.getNextProvider(), req.route), req.r.Body)
	if err != nil {
		log.WithField("err", err).Error("Failed to create provider http request")
		req.w.WriteHeader(http.StatusInternalServerError)
		done <- struct{}{}
		return
	}
	copyReqHeader(req.r, newReq)
	res, err := lb.providerClient.Do(newReq)
	if err != nil || res.StatusCode != http.StatusOK {
		log.WithField("err", err).Error("Failed to contact provider server. Forwarding to next provider.")
		go lb.forwardSyncRequest(req, done)
		return
	}
	copyResHeader(res, req.w)
	io.Copy(req.w, res.Body)
	done <- struct{}{}
}

// Foward async tasks towards the providers (e.g. MD Compilation)
// Returns success to the client as long as the server receives the payload
func (lb *LoadBalancer) forwardAsyncRequest(req *request, done chan struct{}) {

}
