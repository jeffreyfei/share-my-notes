package load_balancer

import (
	"fmt"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

func copyReqHeader(inReq *http.Request, outReq *http.Request) {
	for k, v := range inReq.Header {
		values := ""
		for i, headerValue := range v {
			values += headerValue
			if i != len(v)-1 {
				values += " "
			}
		}
		outReq.Header.Add(k, values)
	}
}

func getForm(inReq *http.Request) url.Values {
	inReq.ParseForm()
	form := url.Values{}
	for k, v := range inReq.Form {
		values := ""
		for i, formValue := range v {
			values += formValue
			if i != len(v)-1 {
				values += " "
			}
		}
		form.Add(k, values)
	}
	return form
}

func copyResHeader(inRes *http.Response, outRes http.ResponseWriter) {
	for k, v := range inRes.Header {
		values := ""
		for i, headerValue := range v {
			values += headerValue
			if i != len(v)-1 {
				values += " "
			}
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
func (lb *LoadBalancer) forwardRedirectRequest(req *request, done chan string) {
	provider, err := lb.getNextProvider()
	if err != nil {
		log.WithField("err", err).Error("Failed to get provider.")
	}
	newReq, err := http.NewRequest(req.r.Method, fmt.Sprintf("%s%s", provider, req.route), req.r.Body)
	newReq.URL.RawQuery = req.r.URL.RawQuery
	if err != nil {
		log.WithField("err", err).Error("Failed to create provider http request")
		req.w.WriteHeader(http.StatusInternalServerError)
		done <- ""
		return
	}
	copyReqHeader(req.r, newReq)
	res, err := lb.providerClient.Do(newReq)
	if err != nil || res.StatusCode != http.StatusOK {
		log.WithField("err", err).Error("Failed to contact provider server. Forwarding to next provider.")
		go lb.forwardRedirectRequest(req, done)
		return
	}
	copyResHeader(res, req.w)
	done <- res.Header.Get("url")
}

// Foward async tasks towards the providers (e.g. MD Compilation)
// Returns success to the client as long as the server receives the payload
func (lb *LoadBalancer) forwardAsyncRequest(req *request, done chan struct{}) {}
