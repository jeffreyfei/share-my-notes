package server

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gorilla/mux"

	"github.com/stretchr/testify/suite"
)

type ReportHandlerTestSuite struct {
	suite.Suite
	s *Server
}

func TestReportHandlerTestSuite(t *testing.T) {
	s := new(ReportHandlerTestSuite)
	s.s = NewServer(nil, "", "", "", "", "", "")
	suite.Run(t, s)
}

func (s *ReportHandlerTestSuite) TestReportStatusHandler() {
	emptyFunc := func(i interface{}, done chan interface{}, err chan error) {}
	s.s.buffer.NewJob(emptyFunc, "", make(chan interface{}), make(chan error))
	s.s.buffer.NewJob(emptyFunc, "", make(chan interface{}), make(chan error))
	req, err := http.NewRequest("GET", "/report/status", nil)
	rec := httptest.NewRecorder()
	assert.NoError(s.T(), err)
	router := mux.NewRouter()
	router.HandleFunc("/report/status", s.s.reportStatusHandler)
	router.ServeHTTP(rec, req)
	assert.Equal(s.T(), http.StatusOK, rec.Code)
	body, err := ioutil.ReadAll(rec.Body)
	assert.Equal(s.T(), string(body), "2")
}
