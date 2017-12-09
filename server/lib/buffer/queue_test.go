package buffer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"
)

type QueueTestSuite struct {
	suite.Suite
}

func createMockJobQueueEntry(op int, payload interface{}) jobQueueEntry {
	ch := make(chan interface{})
	return jobQueueEntry{
		op,
		payload,
		ch,
		time.Now(),
	}
}

func TestQueueTestSuite(t *testing.T) {
	s := new(QueueTestSuite)
	suite.Run(t, s)
}

func (s *QueueTestSuite) TestSize() {
	q := jobQueue{
		[]jobQueueEntry{},
		10,
	}
	assert.Equal(s.T(), 10, q.Size())
}

func (s *QueueTestSuite) TestEnqueue() {
	q := new(jobQueue)
	q.Enqueue(1, "mock-content", make(chan interface{}))
	q.Enqueue(2, "mock-content", make(chan interface{}))
	assert.Equal(s.T(), 2, q.size)
	assert.Equal(s.T(), 2, len(q.storage))
}

func (s *QueueTestSuite) TestDequeue() {
	mockEntry1 := createMockJobQueueEntry(1, "mock-content")
	mockEntry2 := createMockJobQueueEntry(2, "mock-content")
	mockEntry3 := createMockJobQueueEntry(3, "mock-content")
	mockEntry4 := createMockJobQueueEntry(4, "mock-content")
	q := new(jobQueue)
	q.storage = []jobQueueEntry{mockEntry1, mockEntry2, mockEntry3, mockEntry4}
	q.size = 4
	subQueue, err := q.Dequeue(3)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, len(q.storage))
	assert.Equal(s.T(), 1, q.size)
	assert.Equal(s.T(), 3, len(subQueue))
	assert.Equal(s.T(), mockEntry4, q.storage[0])
	subQueue, err = q.Dequeue(3)
	assert.NoError(s.T(), err)
	assert.Empty(s.T(), q.storage)
	assert.Equal(s.T(), 0, q.size)
	assert.Equal(s.T(), 1, len(subQueue))
}
