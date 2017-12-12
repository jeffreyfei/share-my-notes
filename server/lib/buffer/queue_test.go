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

func createMockJobQueueEntry(op JobActionFunc, payload interface{}, done chan interface{}, errCh chan error) jobQueueEntry {
	return jobQueueEntry{
		op,
		payload,
		done,
		errCh,
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
	q.Enqueue(mockActionFunc, "mock-content", make(chan interface{}), make(chan error))
	q.Enqueue(mockActionFunc, "mock-content", make(chan interface{}), make(chan error))
	assert.Equal(s.T(), 2, q.size)
	assert.Equal(s.T(), 2, len(q.storage))
}

func (s *QueueTestSuite) TestDequeue() {
	mockEntry1 := createMockJobQueueEntry(mockActionFunc, "mock-content", make(chan interface{}), make(chan error))
	mockEntry2 := createMockJobQueueEntry(mockActionFunc, "mock-content", make(chan interface{}), make(chan error))
	mockEntry3 := createMockJobQueueEntry(mockActionFunc, "mock-content", make(chan interface{}), make(chan error))
	mockEntry4 := createMockJobQueueEntry(mockActionFunc, "mock-content", make(chan interface{}), make(chan error))
	q := new(jobQueue)
	q.storage = []jobQueueEntry{mockEntry1, mockEntry2, mockEntry3, mockEntry4}
	q.size = 4
	subQueue, err := q.Dequeue(3)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, len(q.storage))
	assert.Equal(s.T(), 1, q.size)
	assert.Equal(s.T(), 3, len(subQueue))
	subQueue, err = q.Dequeue(3)
	assert.NoError(s.T(), err)
	assert.Empty(s.T(), q.storage)
	assert.Equal(s.T(), 0, q.size)
	assert.Equal(s.T(), 1, len(subQueue))
}
