package http_helper

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"
)

type JsonTestSuite struct {
	suite.Suite
}

func TestJsonTestSuite(t *testing.T) {
	s := new(JsonTestSuite)
	suite.Run(t, s)
}

func (s *JsonTestSuite) TestGetJSONFromRequest() {
	type testJSON struct {
		TestID     int    `json:"testID"`
		TestString string `json:"testString"`
	}
	jsonStr := []byte(`{"testID":1,"testString":"test-string"}`)
	req, err := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonStr))
	assert.NoError(s.T(), err)
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		test := testJSON{}
		GetJSONFromRequest(r, &test)
		assert.Equal(s.T(), 1, test.TestID)
		assert.Equal(s.T(), "test-string", test.TestString)
	}).Methods("POST")
	router.ServeHTTP(res, req)
}
