package user

import (
	"fmt"
	"testing"
	"time"

	"github.com/jeffreyfei/share-my-notes/server/lib/db"
	"github.com/jeffreyfei/share-my-notes/server/lib/md_note"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserModelTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func createMockNote(ownerID int64) md_note.MDNoteModel {
	note := md_note.MDNoteModel{}
	note.OwnerID = ownerID
	return note
}

func createMockUser(id int64, googleID string) UserModel {
	user := UserModel{}
	user.ID = id
	user.GoogleID = googleID
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
	s.clearData()
}

func (s *UserModelTestSuite) clearData() {
	assert.NoError(s.T(), s.db.Exec("TRUNCATE TABLE user_models").Error)
	assert.NoError(s.T(), s.db.Exec("TRUNCATE TABLE md_note_models").Error)
}

func (s *UserModelTestSuite) TestNewUser() {
	user := createMockUser(int64(1), "1")
	assert.NoError(s.T(), NewUser(s.db, &user))
	var savedUser UserModel
	assert.NoError(s.T(), s.db.First(&savedUser).Error)
	assert.Equal(s.T(), user.GoogleID, savedUser.GoogleID)
	assert.Equal(s.T(), user.Name, savedUser.Name)
	assert.Equal(s.T(), user.ImageURL, savedUser.ImageURL)
}

func (s *UserModelTestSuite) TestHandleLoginNewUser() {
	user := createMockUser(int64(1), "1")
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

func (s *UserModelTestSuite) TestGetUserByGoogleID() {
	mockUser := createMockUser(int64(1), "1234")
	assert.NoError(s.T(), s.db.Create(&mockUser).Error)
	user, err := GetUserByGoogleID(s.db, "1234")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user.ID, int64(1))
	assert.Equal(s.T(), user.GoogleID, "1234")
}

func (s *UserModelTestSuite) TestHandleLoginExistingUser() {
	user := createMockUser(int64(1), "1")
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

func (s *UserModelTestSuite) TestMDNotes() {
	user1 := createMockUser(int64(1), "1")
	user2 := createMockUser(int64(2), "2")
	note1 := createMockNote(int64(1))
	note2 := createMockNote(int64(1))
	note3 := createMockNote(int64(1))
	note4 := createMockNote(int64(2))
	note5 := createMockNote(int64(2))
	assert.NoError(s.T(), s.db.Create(&note1).Error)
	assert.NoError(s.T(), s.db.Create(&note2).Error)
	assert.NoError(s.T(), s.db.Create(&note3).Error)
	assert.NoError(s.T(), s.db.Create(&note4).Error)
	assert.NoError(s.T(), s.db.Create(&note5).Error)
	notes, err := user1.MDNotes(s.db)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 3, len(notes))
	notes, err = user2.MDNotes(s.db)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 2, len(notes))
}
