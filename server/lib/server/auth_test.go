package server

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"
	plus "google.golang.org/api/plus/v1"
)

type AuthTestSuite struct {
	suite.Suite
	s          *Server
	sessionKey string
}

func TestAuthTestSuite(t *testing.T) {
	s := new(AuthTestSuite)
	s.sessionKey = "mock-session-key"
	s.s = NewServer(nil, "", s.sessionKey, "", "", "", "")
	suite.Run(t, s)
}

func (s *AuthTestSuite) TestValidateRedirectURL() {
	url, err := validateRedirectURL("/test/path")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "/test/path", url)
	url, err = validateRedirectURL("")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "/", url)
	url, err = validateRedirectURL("https://www.google.ca")
	assert.Error(s.T(), err)
	assert.Equal(s.T(), "/", url)
}

func (s *AuthTestSuite) TestFormatProfile() {
	mockProfile := new(plus.Person)
	mockProfile.DisplayName = "mock-name"
	mockProfile.Id = "mock-id"
	mockProfile.Image = new(plus.PersonImage)
	mockProfile.Image.Url = "mock-image-url"
	userModel := formatProfile(mockProfile)
	assert.Equal(s.T(), userModel.GoogleID, mockProfile.Id)
	assert.Equal(s.T(), userModel.Name, mockProfile.DisplayName)
	assert.Equal(s.T(), userModel.ImageURL, mockProfile.Image.Url)
}
