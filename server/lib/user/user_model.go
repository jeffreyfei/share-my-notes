package user

import (
	"time"

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

func HandleLogin(db *gorm.DB, user *UserModel) (*UserModel, error) {
	var existingUser UserModel
	if err := db.Where("google_id = ?", user.GoogleID).First(&existingUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, NewUser(db, user)
		}
		return nil, err
	}
	if err := existingUser.UpdateLoginTime(db); err != nil {
		return nil, err
	}
	return &existingUser, nil
}

func NewUser(db *gorm.DB, user *UserModel) error {
	user.CreatedAt = time.Now()
	user.LastLoggedInAt = time.Now()
	return db.Save(user).Error
}

func (u *UserModel) UpdateLoginTime(db *gorm.DB) error {
	u.LastLoggedInAt = time.Now()
	return db.Save(u).Error
}
