package load_balancer

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"
)

type ForwarderTestSuite struct {
	suite.Suite
	lb *LoadBalancer
}

type mockClient struct{}

func (c *mockClient) Do(req *http.Request) (*http.Response, error) {
	res := http.Response{}
	res.Header.Add("url", "mock-url")
	return &res, nil
}

func TestForwarderTestSuite(t *testing.T) {
	s := new(ForwarderTestSuite)
	s.lb = NewLoadBalancer(1000)
	s.lb.providerClient = new(mockClient)
	suite.Run(t, s)
}

func (s *ForwarderTestSuite) TestCopyReqHeader() {
	req, err := http.NewRequest("GET", "bogus", nil)
	assert.NoError(s.T(), err)
	req.Header.Add("mock-header1", "mock-value1")
	req.Header.Add("mock-header2", "mock-value2")
	req.Header.Add("mock-header3", "mock-value3")
	copiedReq, err := http.NewRequest("GET", "bogus", nil)
	copyReqHeader(req, copiedReq)
	assert.Equal(s.T(), "mock-value1", copiedReq.Header.Get("mock-header1"))
	assert.Equal(s.T(), "mock-value2", copiedReq.Header.Get("mock-header2"))
	assert.Equal(s.T(), "mock-value3", copiedReq.Header.Get("mock-header3"))
}

func (s *ForwarderTestSuite) TestGetForm() {
	req, err := http.NewRequest("GET", "bogus", strings.NewReader(""))
	req.ParseForm()
	assert.NoError(s.T(), err)
	req.Form.Add("mock-form1", "mock-value1")
	req.Form.Add("mock-form2", "mock-value2")
	req.Form.Add("mock-form3", "mock-value3")
	form := getForm(req)
	assert.Equal(s.T(), "mock-value1", form.Get("mock-form1"))
	assert.Equal(s.T(), "mock-value2", form.Get("mock-form2"))
	assert.Equal(s.T(), "mock-value3", form.Get("mock-form3"))
}

func (s *ForwarderTestSuite) TestCopyResHeader() {
	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(s.T(), err)
	res := new(http.Response)
	res.Header = make(http.Header)
	res.Header.Add("mock-header1", "mock-value1")
	res.Header.Add("mock-header2", "mock-value2")
	res.Header.Add("mock-header3", "mock-value3")
	router := mux.NewRouter()
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		copyResHeader(res, w)
		assert.Equal(s.T(), "mock-value1", w.Header().Get("mock-header1"))
		assert.Equal(s.T(), "mock-value2", w.Header().Get("mock-header2"))
		assert.Equal(s.T(), "mock-value3", w.Header().Get("mock-header3"))
	})
	emptyRes := httptest.NewRecorder()
	router.ServeHTTP(emptyRes, req)
}

func (s *ForwarderTestSuite) TestForwardRedirectRequest() {
	req, err := http.NewRequest("GET", "bogus", strings.NewReader(""))
	req.ParseForm()
	assert.NoError(s.T(), err)
	req.Form.Add("mock-form1", "mock-value1")
	req.Form.Add("mock-form2", "mock-value2")
	req.Header.Add("mock-header1", "mock-value1")
	req.Header.Add("mock-header2", "mock-value2")
	router := mux.NewRouter()
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		done := make(chan string)
		testReq := new(request)
		testReq.r = r
		testReq.w = w
		go s.lb.forwardRedirectRequest(testReq, done)
		url := <-done
		assert.Equal(s.T(), "mock-url", url)
	})
	emptyRes := httptest.NewRecorder()
	router.ServeHTTP(emptyRes, req)
}
