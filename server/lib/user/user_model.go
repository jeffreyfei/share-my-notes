package user

import (
	"time"

	"github.com/jeffreyfei/share-my-notes/server/lib/md_note"

	"github.com/jinzhu/gorm"
)

type UserModel struct {
	ID             int64  `gorm:"primary_key"`
	GoogleID       string `gorm:"type:varchar(25);unique_index"`
	Name           string `gorm:"type:varchar(64)"`
	ImageURL       string `gorm:"type:varchar(255)"`
	CreatedAt      time.Time
	LastLoggedInAt time.Time
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&UserModel{}).Error
}

// Handles user login
// If the user exists, retrieve and return UserModel fom database
// If the user does not exist, create a new user in the database and return the UserModel
func HandleLogin(db *gorm.DB, user *UserModel) (*UserModel, error) {
	var existingUser UserModel
	if err := db.Where("google_id = ?", user.GoogleID).First(&existingUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, NewUser(db, user)
		}
		return nil, err
	}
	if err := existingUser.updateLoginTime(db); err != nil {
		return nil, err
	}
	return &existingUser, nil
}

// Retrieves the UserModel corresponding to the given Google ID
func GetUserByGoogleID(db *gorm.DB, googleID string) (*UserModel, error) {
	var user UserModel
	if err := db.Where("google_id = ?", googleID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Creates a new user in the database
func NewUser(db *gorm.DB, user *UserModel) error {
	user.CreatedAt = time.Now()
	user.LastLoggedInAt = time.Now()
	return db.Create(user).Error
}

// Update the login time of a user
func (u *UserModel) updateLoginTime(db *gorm.DB) error {
	u.LastLoggedInAt = time.Now()
	return db.Save(u).Error
}

// Retrieve all MD note entries owned by the given user
func (u *UserModel) MDNotes(db *gorm.DB) ([]md_note.MDNoteModel, error) {
	var mdNotes []md_note.MDNoteModel
	if err := db.Where("owner_id = ?", u.ID).Find(&mdNotes).Error; err != nil {
		return mdNotes, err
	}
	return mdNotes, nil
}
