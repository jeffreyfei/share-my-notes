package md_note

import (
	"fmt"
	"testing"

	"github.com/jeffreyfei/share-my-notes/server/lib/db"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MDNoteTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func createMockNote(id, ownerID int64) MDNoteModel {
	note := MDNoteModel{}
	note.ID = id
	note.OwnerID = ownerID
	note.Category = fmt.Sprintf("mock-cat-%d", id)
	note.RawText = fmt.Sprintf("mock-raw-text-%d", id)
	note.CompiledText = fmt.Sprintf("mock-comp-text-%d", id)
	return note
}

func TestMDNoteTestSuite(t *testing.T) {
	s := new(MDNoteTestSuite)
	var err error
	s.db, err = db.GetDB()
	assert.NoError(t, err)
	assert.NoError(t, AutoMigrate(s.db))
	suite.Run(t, s)
}

func (s *MDNoteTestSuite) SetupTest() {
	s.clearData()
}

func (s *MDNoteTestSuite) clearData() {
	assert.NoError(s.T(), s.db.Exec("TRUNCATE TABLE md_note_models").Error)
}

func (s *MDNoteTestSuite) TestNew() {
	assert.NoError(s.T(), New(s.db, int64(1), "mock-cat", "### Test"))
	var newNote MDNoteModel
	assert.NoError(s.T(), s.db.First(&newNote).Error)
	assert.Equal(s.T(), int64(1), newNote.OwnerID)
	assert.Equal(s.T(), "mock-cat", newNote.Category)
	assert.Equal(s.T(), "### Test", newNote.RawText)
	assert.Equal(s.T(), "<h3><a href=\"#test\" rel=\"nofollow\"><span></span></a>Test</h3>\n", newNote.CompiledText)
}

func (s *MDNoteTestSuite) TestGet() {
	note := createMockNote(1, 1)
	assert.NoError(s.T(), s.db.Create(&note).Error)
	savedNote, err := Get(s.db, 1)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), note.ID, savedNote.ID)
	assert.Equal(s.T(), note.OwnerID, savedNote.OwnerID)
	assert.Equal(s.T(), note.Category, savedNote.Category)
	assert.Equal(s.T(), note.CompiledText, savedNote.CompiledText)
	assert.Equal(s.T(), note.RawText, savedNote.RawText)
}

func (s *MDNoteTestSuite) TestUpdate() {
	note := createMockNote(1, 1)
	assert.NoError(s.T(), s.db.Create(&note).Error)
	assert.NoError(s.T(), Update(s.db, int64(1), "### Test"))
	var savedNote MDNoteModel
	assert.NoError(s.T(), s.db.First(&savedNote, 1).Error)
	assert.Equal(s.T(), "<h3><a href=\"#test\" rel=\"nofollow\"><span></span></a>Test</h3>\n", savedNote.CompiledText)
}

func (s *MDNoteTestSuite) TestDelete() {
	note := createMockNote(1, 1)
	assert.NoError(s.T(), s.db.Create(&note).Error)
	assert.NoError(s.T(), Delete(s.db, 1))
	var deletedNote MDNoteModel
	assert.True(s.T(), s.db.First(&deletedNote, 1).RecordNotFound())
}
