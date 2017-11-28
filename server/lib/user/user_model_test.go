package user

import (
	"fmt"
	"testing"
	"time"

	"github.com/jeffreyfei/share-my-notes/server/lib/server"

	"github.com/jeffreyfei/share-my-notes/server/lib/db"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserModelTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func createMockProfile(id string) server.Profile {
	return server.Profile{
		id,
		fmt.Sprintf("mock-name-%s", id),
		fmt.Sprintf("mock-image-url-%s", id),
	}
}

func createMockUser(p *server.Profile) UserModel {
	user := UserModel{}
	user.GoogleID = p.ID
	user.Name = p.DisplayName
	user.ImageURL = p.ImageURL
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

func (s *UserModelTestSuite) TestExists() {
	profile1 := createMockProfile("1")
	user1 := createMockUser(&profile1)
	assert.NoError(s.T(), s.db.Save(&user1).Error)
	exists, err := Exists(s.db, &profile1)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), true, exists)
	profile2 := createMockProfile("2")
	exists, err = Exists(s.db, &profile2)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), false, exists)
}

func (s *UserModelTestSuite) TestNewUser() {
	profile := createMockProfile("123456789")
	assert.NoError(s.T(), NewUser(s.db, &profile))
	var newUser UserModel
	assert.NoError(s.T(), s.db.First(&newUser).Error)
	assert.Equal(s.T(), profile.ID, newUser.GoogleID)
	assert.Equal(s.T(), profile.DisplayName, newUser.Name)
	assert.Equal(s.T(), profile.ImageURL, newUser.ImageURL)
}
