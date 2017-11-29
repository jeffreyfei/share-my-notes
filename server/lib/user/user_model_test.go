package user

import (
	"fmt"
	"testing"
	"time"

	"github.com/jeffreyfei/share-my-notes/server/lib/db"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserModelTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func createMockUser(id string) UserModel {
	user := UserModel{}
	user.GoogleID = id
	user.Name = fmt.Sprintf("mock-name-%s", id)
	user.ImageURL = fmt.Sprintf("mock-image-url-%s", id)
	user.CreatedAt = time.Now()
	user.LastLoggedInAt = time.Now()
	return user
}

func TestUserModelTestSuite(t *testing.T) {
	s := new(UserModelTestSuite)
	var err error
	s.db, err = db.GetDB()
	assert.NoError(t, err)
	assert.NoError(t, AutoMigrate(s.db))
	suite.Run(t, s)
}

func (s *UserModelTestSuite) SetupTest() {
	assert.NoError(s.T(), s.db.Exec("DELETE FROM user_models").Error)
}

func (s *UserModelTestSuite) TestNewUser() {
	user := createMockUser("1")
	assert.NoError(s.T(), NewUser(s.db, &user))
	var savedUser UserModel
	assert.NoError(s.T(), s.db.First(&savedUser).Error)
	assert.Equal(s.T(), user.GoogleID, savedUser.GoogleID)
	assert.Equal(s.T(), user.Name, savedUser.Name)
	assert.Equal(s.T(), user.ImageURL, savedUser.ImageURL)
}

func (s *UserModelTestSuite) TestHandleLoginNewUser() {
	user := createMockUser("1")
	returnedUser, err := HandleLogin(s.db, &user)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user.GoogleID, returnedUser.GoogleID)
	assert.Equal(s.T(), user.Name, returnedUser.Name)
	assert.Equal(s.T(), user.ImageURL, returnedUser.ImageURL)
	var savedUser UserModel
	assert.NoError(s.T(), s.db.First(&savedUser).Error)
	assert.Equal(s.T(), user.GoogleID, savedUser.GoogleID)
	assert.Equal(s.T(), user.Name, savedUser.Name)
	assert.Equal(s.T(), user.ImageURL, savedUser.ImageURL)
}

func (s *UserModelTestSuite) TestHandleLoginExistingUser() {
	user := createMockUser("1")
	assert.NoError(s.T(), s.db.Save(&user).Error)
	returnedUser, err := HandleLogin(s.db, &user)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user.GoogleID, returnedUser.GoogleID)
	assert.Equal(s.T(), user.Name, returnedUser.Name)
	assert.Equal(s.T(), user.ImageURL, returnedUser.ImageURL)
	var savedUser UserModel
	assert.NoError(s.T(), s.db.First(&savedUser).Error)
	assert.Equal(s.T(), user.GoogleID, savedUser.GoogleID)
	assert.Equal(s.T(), user.Name, savedUser.Name)
	assert.Equal(s.T(), user.ImageURL, savedUser.ImageURL)
}
