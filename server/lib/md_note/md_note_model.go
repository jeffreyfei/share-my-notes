package md_note

import (
	"time"

	"github.com/jinzhu/gorm"
)

type MDNoteModel struct {
	ID           int64 `gorm:"primary_key"`
	OwnerID        int64
	Category     string `gorm:"type:varchar(50)"`
	RawText      string `sql:"type:text"`
	CompiledText string `sql:"type:text"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&MDNoteModel{}).Error
}
