package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gorilla/mux"

	"github.com/stretchr/testify/suite"
)

type RouterTestSuite struct {
	suite.Suite
}

func TestRouterTestSuite(t *testing.T) {
	s := new(RouterTestSuite)
	suite.Run(t, s)
}

func (s *RouterTestSuite) TestBuildRouter() {
	flag1 := "init"
	flag2 := "init"
	routes := Routes{
		Route{
			"GET",
			"/test/{id}",
			func(w http.ResponseWriter, r *http.Request) {
				id := mux.Vars(r)["id"]
				flag1 = id
			},
		},
		Route{
			"POST",
			"/test/{id}",
			func(w http.ResponseWriter, r *http.Request) {
				id := mux.Vars(r)["id"]
				flag2 = id
			},
		},
	}
	router := BuildRouter(routes)
	req, err := http.NewRequest("POST", "/test/mutated", nil)
	assert.NoError(s.T(), err)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(s.T(), http.StatusOK, rec.Code)
	assert.Equal(s.T(), "init", flag1)
	assert.Equal(s.T(), "mutated", flag2)
}
