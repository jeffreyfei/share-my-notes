package server

import (
	"testing"
	"time"

	"github.com/jeffreyfei/share-my-notes/server/lib/md_note"

	"github.com/stretchr/testify/assert"

	"github.com/jeffreyfei/share-my-notes/server/lib/db"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
)

type MDHandlerTestSuite struct {
	suite.Suite
	db *gorm.DB
	s  *Server
}

func createMockMDCompilePayload() mdCompilePayload {
	return mdCompilePayload{
		1,
		"mock-text",
		"mock-cat",
	}
}

func createMockUpdatePayload() mdUpdatePayload {
	return mdUpdatePayload{
		1,
		"mock-raw-text",
	}
}

func createMockMDNote(id int64) md_note.MDNoteModel {
	return md_note.MDNoteModel{
		id,
		id,
		"mock-cat",
		"mock-raw",
		"mock-comp",
		time.Now(),
		time.Now(),
	}
}

func TestMDHandlerTestSuite(t *testing.T) {
	s := new(MDHandlerTestSuite)
	var err error
	s.db, err = db.GetDB()
	s.s = NewServer(s.db, "", "", "", "", "/mockaddr")
	assert.NoError(t, err)
	suite.Run(t, s)
}

func (s *MDHandlerTestSuite) SetupTest() {
	s.clearData()
}

func (s *MDHandlerTestSuite) clearData() {
	assert.NoError(s.T(), s.db.Exec("TRUNCATE TABLE md_note_models").Error)
}

func (s *MDHandlerTestSuite) TestMDCreateAction() {
	payload := createMockMDCompilePayload()
	done := make(chan interface{})
	go s.s.mdCreateAction(payload, done, make(chan error))
	result := <-done
	assert.Equal(s.T(), struct{}{}, result)
}

func (s *MDHandlerTestSuite) TestMDCreateCallback() {
	payload := createMockMDCompilePayload()
	done := make(chan interface{})
	go func() {
		done <- struct{}{}
	}()
	s.s.mdCreateCallback(payload, done, make(chan error))
}

func (s *MDHandlerTestSuite) TestMDGetAction() {
	mockMDNote := createMockMDNote(1)
	assert.NoError(s.T(), s.db.Create(&mockMDNote).Error)
	done := make(chan interface{})
	errCh := make(chan error)
	go s.s.mdGetAction(mockMDNote.ID, done, errCh)
	select {
	case result := <-done:
		payload, ok := result.(mdGetPayload)
		assert.True(s.T(), ok)
		assert.Equal(s.T(), mockMDNote.ID, payload.ID)
		assert.Equal(s.T(), mockMDNote.Category, payload.Category)
		assert.Equal(s.T(), mockMDNote.CompiledText, payload.CompiledText)
		assert.Equal(s.T(), mockMDNote.RawText, payload.RawText)
	case err := <-errCh:
		assert.NoError(s.T(), err)
	}
}

func (s *MDHandlerTestSuite) TestMDGetCallback() {
	mockRetPayload := mdGetPayload{
		1,
		"mock-raw",
		"mock-comp",
		"mock-cat",
	}
	done := make(chan interface{})
	go func() {
		done <- mockRetPayload
	}()
	s.s.mdGetCallback(int64(1), done, make(chan error))
}

func (s *MDHandlerTestSuite) TestMDUpdateAction() {
	mockMDNote := createMockMDNote(1)
	assert.NoError(s.T(), s.db.Create(&mockMDNote).Error)
	doneCh := make(chan interface{})
	errCh := make(chan error)
	mockPayload := createMockUpdatePayload()
	go s.s.mdUpdateAction(mockPayload, doneCh, errCh)
	select {
	case <-doneCh:
		var note md_note.MDNoteModel
		assert.NoError(s.T(), s.db.First(&note, 1).Error)
		assert.Equal(s.T(), mockPayload.ID, note.ID)
		assert.Equal(s.T(), mockPayload.RawText, note.RawText)
	case err := <-errCh:
		assert.NoError(s.T(), err)
	}
}

func (s *MDHandlerTestSuite) TestMDUpdateCallback() {
	done := make(chan interface{})
	mockPayload := createMockUpdatePayload()
	go func() {
		done <- struct{}{}
	}()
	s.s.mdUpdateCallback(mockPayload, done, make(chan error))
}

func (s *MDHandlerTestSuite) TestMDDeleteAction() {
	mockMDNote1 := createMockMDNote(1)
	mockMDNote2 := createMockMDNote(2)
	assert.NoError(s.T(), s.db.Create(&mockMDNote1).Error)
	assert.NoError(s.T(), s.db.Create(&mockMDNote2).Error)
	doneCh := make(chan interface{})
	errCh := make(chan error)
	go s.s.mdDeleteAction(int64(1), doneCh, errCh)
	select {
	case <-doneCh:
		var savedNote md_note.MDNoteModel
		assert.True(s.T(), s.db.First(&savedNote, 1).RecordNotFound())
		assert.False(s.T(), s.db.First(&savedNote, 2).RecordNotFound())
	case err := <-errCh:
		assert.NoError(s.T(), err)
	}
}

func (s *MDHandlerTestSuite) TestMDDeleteCallback() {
	done := make(chan interface{})
	go func() {
		done <- struct{}{}
	}()
	s.s.mdDeleteCallback(int64(1), done, make(chan error))
}
