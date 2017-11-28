package user

import (
	"time"

	"github.com/jeffreyfei/share-my-notes/server/lib/server"
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

func Exists(db *gorm.DB, profile *server.Profile) (bool, error) {
	var user UserModel
	if err := db.Where("google_id = ?", profile.ID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func NewUser(db *gorm.DB, profile *server.Profile) error {
	newUser := UserModel{}
	newUser.GoogleID = profile.ID
	newUser.Name = profile.DisplayName
	newUser.ImageURL = profile.ImageURL
	newUser.CreatedAt = time.Now()
	newUser.LastLoggedInAt = time.Now()
	return db.Save(&newUser).Error
}
