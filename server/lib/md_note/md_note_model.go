package md_note

import (
	"time"

	"github.com/jinzhu/gorm"
)

type MDNoteModel struct {
	ID           int64 `gorm:"primary_key"`
	OwnerID      int64
	Category     string `gorm:"type:varchar(50)"`
	RawText      string `sql:"type:text"`
	CompiledText string `sql:"type:text"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&MDNoteModel{}).Error
}

func New(db *gorm.DB, ownerID int64, category, rawText string) error {
	newNote := MDNoteModel{}
	newNote.OwnerID = ownerID
	newNote.Category = category
	newNote.RawText = rawText
	newNote.CompiledText = CompileMD(rawText)
	newNote.CreatedAt = time.Now()
	newNote.UpdatedAt = time.Now()
	return db.Create(&newNote).Error
}

func Get(db *gorm.DB, id int64) (MDNoteModel, error) {
	var note MDNoteModel
	if err := db.First(&note, id).Error; err != nil {
		return MDNoteModel{}, err
	}
	return note, nil
}

func Update(db *gorm.DB, id int64, rawText string) error {
	var note MDNoteModel
	if err := db.First(&note, id).Error; err != nil {
		return err
	}
	note.RawText = rawText
	note.CompiledText = CompileMD(rawText)
	note.UpdatedAt = time.Now()
	return db.Save(&note).Error
}

func Delete(db *gorm.DB, id int64) error {
	deletedNote := MDNoteModel{ID: id}
	return db.Delete(deletedNote).Error
}
