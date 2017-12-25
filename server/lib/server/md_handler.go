package server

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/jeffreyfei/share-my-notes/server/lib/http_helper"
	"github.com/jeffreyfei/share-my-notes/server/lib/md_note"
	"github.com/jeffreyfei/share-my-notes/server/lib/user"
	log "github.com/sirupsen/logrus"
)

func makeLBResponse(url string, body []byte) error {
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(body))
	if err != nil {
		log.WithField("err", err).Error("Failed to create load balancer request")
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.WithField("err", err).Error("Load balancer request failed")
		return err
	}
	if res != nil && res.StatusCode != http.StatusOK {
		log.WithField("code", err).Error("Load balancer request failed")
		return err
	}
	return nil
}

type errorPayload struct {
	Error string `json:"error"`
}

type mdCompilePayload struct {
	OwnerID  int64  `json:"ownerID"`
	RawText  string `json:"rawText"`
	Category string `json:"category"`
}

// Creates the MD note entry in the database
// Action function that gets added to the job queue
func (s *Server) mdCreateAction(payload interface{}, doneCh chan interface{}, errCh chan error) {
	defer close(doneCh)
	defer close(errCh)
	mdPayload := payload.(mdCompilePayload)
	err := md_note.New(s.db, mdPayload.OwnerID, mdPayload.Category, mdPayload.RawText)
	if err != nil {
		errCh <- err
	} else {
		doneCh <- struct{}{}
	}
}

// Creates a MD creation job in the job queue
// Returns the result back to the load balancer once the job is finished
func (s *Server) mdCreateCallback(recPayload interface{}, doneCh chan interface{}, errCh chan error) {
	s.buffer.NewJob(s.mdCreateAction, recPayload, doneCh, make(chan error))
	var body []byte
	select {
	case <-doneCh:
	case err := <-errCh:
		if body, err = http_helper.ParseJSONBody(&errorPayload{err.Error()}); err != nil {
			log.WithField("err", err).Error("Failed to parse error response")
		}
		log.WithField("err", err).Error("Failed to create MD Notes")
	}
	url := fmt.Sprintf("%s/response/md/create", s.lbPrivateURL)
	if err := makeLBResponse(url, body); err != nil {
		log.WithField("err", err).Error("Failed to make LB Request")
	}
}

// Handles Markdown Notes creation
// Calls the MD Create callback function to create a creation job on the job queue
// Does not wait for the creation job to finish
func (s *Server) mdCreateHandler(w http.ResponseWriter, r *http.Request) {
	profile := s.getProfileFromSession(r)
	currentUser, err := user.GetUserByGoogleID(s.db, profile.GoogleID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	payload := mdCompilePayload{}
	if err := http_helper.GetJSONFromRequest(r, &payload); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.WithField("err", err).Error("Failed to parse JSON")
	}
	payload.OwnerID = currentUser.ID
	doneCh := make(chan interface{})
	errCh := make(chan error)
	go s.mdCreateCallback(payload, doneCh, errCh)
	w.WriteHeader(http.StatusOK)
}

type mdGetPayload struct {
	ID           int64  `json:"id"`
	RawText      string `json:"rawText"`
	CompiledText string `json:"compiledText"`
	Category     string `json:"category"`
}

// Retrieves a MD note entry from database
// Action function that gets added to the job queue
func (s *Server) mdGetAction(payload interface{}, doneCh chan interface{}, errCh chan error) {
	defer close(doneCh)
	defer close(errCh)
	id := payload.(int64)
	if note, err := md_note.Get(s.db, id); err != nil {
		errCh <- err
	} else {
		recPayload := mdGetPayload{
			note.ID,
			note.RawText,
			note.CompiledText,
			note.Category,
		}
		doneCh <- recPayload
	}
}

// Creates a mdGetAction job on the job queue
// Returns the retrieve note entry to the load balancer once the job is processed
func (s *Server) mdGetCallback(recPayload interface{}, doneCh chan interface{}, errCh chan error) {
	id := recPayload.(int64)
	s.buffer.NewJob(s.mdGetAction, id, doneCh, errCh)
	var body []byte
	var err error
	// Determine if the returned value is an error
	select {
	case result := <-doneCh:
		reqPayload := result.(mdGetPayload)
		if body, err = http_helper.ParseJSONBody(&reqPayload); err != nil {
			log.WithField("err", err).Error("Failed to parse MD Get response JSON")
		}
	case err := <-errCh:
		if body, err = http_helper.ParseJSONBody(&errorPayload{err.Error()}); err != nil {
			log.WithField("err", err).Error("Failed to parse error response")
		}
		log.WithField("err", err).Error("Failed to get MD Notes")
	}
	url := fmt.Sprintf("%s/response/md/%d/get", s.lbPrivateURL, id)
	if err := makeLBResponse(url, body); err != nil {
		log.WithField("err", err).Error("Failed to make LB Request")
	}
}

// Handles MD note retrieval requests
// Calls mdGetCallback to create a MD note retrieval job on the job queue
// Does not wait for the job to finish
func (s *Server) mdGetHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		log.WithField("err", err).Error("Invalid ID for MD note")
		w.WriteHeader(http.StatusInternalServerError)
	}
	doneCh := make(chan interface{})
	errCh := make(chan error)
	go s.mdGetCallback(id, doneCh, errCh)
	w.WriteHeader(http.StatusOK)
}

type mdUpdatePayload struct {
	ID      int64  `json:"id"`
	RawText string `json:"rawText"`
}

// Updates a MD note entry in the database
// Action function that gets added to the job queue
func (s *Server) mdUpdateAction(payload interface{}, doneCh chan interface{}, errCh chan error) {
	defer close(doneCh)
	defer close(errCh)
	updatePayload := payload.(mdUpdatePayload)
	err := md_note.Update(s.db, updatePayload.ID, updatePayload.RawText)
	if err != nil {
		errCh <- err
	} else {
		doneCh <- struct{}{}
	}
}

// Creates mdUpdateAction on the job queue
// Returns a response to the load balancer when the job is processed
func (s *Server) mdUpdateCallback(recPayload interface{}, doneCh chan interface{}, errCh chan error) {
	s.buffer.NewJob(s.mdUpdateAction, recPayload, doneCh, errCh)
	id := recPayload.(mdUpdatePayload).ID
	var body []byte
	select {
	case <-doneCh:
	case err := <-errCh:
		if body, err = http_helper.ParseJSONBody(&errorPayload{err.Error()}); err != nil {
			log.WithField("err", err).Error("Failed to parse error response")
		}
		log.WithField("err", err).Error("Failed to update MD notes")
	}
	url := fmt.Sprintf("%s/response/md/%d/update", s.lbPrivateURL, id)
	if err := makeLBResponse(url, body); err != nil {
		log.WithField("err", err).Error("Failed to make LB Request")
	}
}

// Handles MD note update requests
// Calls mdUpdateCallback to create a MD update job on the job queue
// Does not wait for the job to finish
func (s *Server) mdUpdateHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		log.WithField("err", err).Error("Invalid ID for MD note")
		w.WriteHeader(http.StatusInternalServerError)
	}
	r.ParseForm()
	updatePayload := mdUpdatePayload{}
	if err := http_helper.GetJSONFromRequest(r, &updatePayload); err != nil {
		log.WithField("err", err).Error("Failed to parse JSON")
	}
	updatePayload.ID = id
	doneCh := make(chan interface{})
	errCh := make(chan error)
	go s.mdCreateCallback(updatePayload, doneCh, errCh)
	w.WriteHeader(http.StatusOK)
}

// Deletes a MD note entry fom database
// Action function that gets added to the job queue
func (s *Server) mdDeleteAction(payload interface{}, doneCh chan interface{}, errCh chan error) {
	defer close(doneCh)
	defer close(errCh)
	id := payload.(int64)
	if err := md_note.Delete(s.db, id); err != nil {
		errCh <- err
	} else {
		doneCh <- struct{}{}
	}
}

// Creates mdDeleteAction on the job queue
// Returns a response to the load balancer when the job is processed
func (s *Server) mdDeleteCallback(recPayload interface{}, doneCh chan interface{}, errCh chan error) {
	id := recPayload.(int64)
	s.buffer.NewJob(s.mdDeleteAction, id, doneCh, errCh)
	var body []byte
	// Determine if the returned value is an error
	select {
	case <-doneCh:
	case err := <-errCh:
		if body, err = http_helper.ParseJSONBody(&errorPayload{err.Error()}); err != nil {
			log.WithField("err", err).Error("Failed to parse error response")
		}
		log.WithField("err", err).Error("Failed to delete MD notes")
	}
	url := fmt.Sprintf("%s/response/md/%d/get", s.lbPrivateURL, id)
	if err := makeLBResponse(url, body); err != nil {
		log.WithField("err", err).Error("Failed to make LB Request")
	}
}

// Handles MD note delete requests
// Calls mdDeleteCallback to create a job on the job queue
// Does not wait for the job to finish
func (s *Server) mdDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		log.WithField("err", err).Error("Invalid ID for MD note")
		w.WriteHeader(http.StatusInternalServerError)
	}
	doneCh := make(chan interface{})
	errCh := make(chan error)
	go s.mdDeleteCallback(id, doneCh, errCh)
	w.WriteHeader(http.StatusOK)
}
