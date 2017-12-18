package server

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/jeffreyfei/share-my-notes/server/lib/md_note"
	log "github.com/sirupsen/logrus"

	"github.com/jeffreyfei/share-my-notes/server/lib/user"
)

type mdCompilePayload struct {
	OwnerID  int64
	RawText  string
	Category string
}

// Action function that will be added to the job queue
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

// Callback function that streams the result back to the load balancer once the job is finished
func (s *Server) mdCreateCallback(recPayload interface{}, doneCh chan interface{}, errCh chan error) {
	s.buffer.NewJob(s.mdCreateAction, recPayload, doneCh, make(chan error))
	form := url.Values{}
	select {
	case <-doneCh:
	case err := <-errCh:
		form.Add("err", err.Error())
	}
	http.PostForm(fmt.Sprintf("%s/response/md/create", s.lbPrivateURL), form)
}

func (s *Server) mdCreateHandler(w http.ResponseWriter, r *http.Request) {
	profile := s.getProfileFromSession(r)
	currentUser, err := user.GetUserByGoogleID(s.db, profile.GoogleID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	r.ParseForm()
	payload := mdCompilePayload{}
	payload.OwnerID = currentUser.ID
	payload.RawText = r.PostFormValue("rawText")
	payload.Category = r.PostFormValue("category")
	doneCh := make(chan interface{})
	errCh := make(chan error)
	go s.mdCreateCallback(payload, doneCh, errCh)
	w.WriteHeader(http.StatusOK)
}

type mdGetPayload struct {
	ID           int64
	RawText      string
	CompiledText string
	Category     string
}

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

func (s *Server) mdGetCallback(recPayload interface{}, doneCh chan interface{}, errCh chan error) {
	id := recPayload.(int64)
	s.buffer.NewJob(s.mdGetAction, id, doneCh, errCh)
	form := url.Values{}
	// Determine if the returned value is an error
	select {
	case result := <-doneCh:
		reqPayload := result.(mdGetPayload)
		form.Add("rawText", reqPayload.RawText)
		form.Add("compiledText", reqPayload.CompiledText)
		form.Add("category", reqPayload.Category)
	case err := <-errCh:
		form.Add("err", err.Error())
		log.WithField("err", err).Error("Getting MD notes failed")
	}
	http.PostForm(fmt.Sprintf("%s/response/md/%d/get", s.lbPrivateURL, id), form)
}

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
	ID      int64
	RawText string
}

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

func (s *Server) mdUpdateCallback(recPayload interface{}, doneCh chan interface{}, errCh chan error) {
	s.buffer.NewJob(s.mdUpdateAction, recPayload, doneCh, errCh)
	id := recPayload.(mdUpdatePayload).ID
	form := url.Values{}
	select {
	case <-doneCh:
	case err := <-errCh:
		form.Add("err", err.Error())
	}
	http.PostForm(fmt.Sprintf("%s/response/md/%d/update", s.lbPrivateURL, id), form)
}

func (s *Server) mdUpdateHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		log.WithField("err", err).Error("Invalid ID for MD note")
		w.WriteHeader(http.StatusInternalServerError)
	}
	r.ParseForm()
	updatePayload := mdUpdatePayload{}
	updatePayload.ID = id
	updatePayload.RawText = r.PostFormValue("rawText")
	doneCh := make(chan interface{})
	errCh := make(chan error)
	go s.mdCreateCallback(updatePayload, doneCh, errCh)
	w.WriteHeader(http.StatusOK)
}

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

func (s *Server) mdDeleteCallback(recPayload interface{}, doneCh chan interface{}, errCh chan error) {
	id := recPayload.(int64)
	s.buffer.NewJob(s.mdDeleteAction, id, doneCh, errCh)
	form := url.Values{}
	// Determine if the returned value is an error
	select {
	case <-doneCh:
	case err := <-errCh:
		form.Add("err", err.Error())
		log.WithField("err", err).Error("Deleting MD notes failed")
	}
	http.PostForm(fmt.Sprintf("%s/response/md/%d/get", s.lbPrivateURL, id), form)
}

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
