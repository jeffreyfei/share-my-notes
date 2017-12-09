package buffer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"
)

type BufferTestSuite struct {
	suite.Suite
}

func mockActionFunc(payload interface{}, ch chan interface{}) {
	ch <- payload.(int) + 1
}

func TestBufferTestSuite(t *testing.T) {
	s := new(BufferTestSuite)
	suite.Run(t, s)
}

func (s *BufferTestSuite) TestNewJob() {
	b := NewBuffer(1000, 4)
	b.NewJob(1, "mock-payload", make(chan interface{}))
	assert.Equal(s.T(), 1, len(b.queue.storage))
	b.NewJob(1, "mock-payload", make(chan interface{}))
	assert.Equal(s.T(), 2, len(b.queue.storage))
}

func (s *BufferTestSuite) TestJobCount() {
	b := NewBuffer(1000, 4)
	jobs := []jobQueueEntry{
		createMockJobQueueEntry(1, "mock-payload", make(chan interface{})),
		createMockJobQueueEntry(1, "mock-payload", make(chan interface{})),
	}
	b.queue = jobQueue{
		jobs,
		2,
	}
	assert.Equal(s.T(), 2, b.JobCount())
}

func (s *BufferTestSuite) TestProcJobQueue() {
	b := NewBuffer(1000, 4)
	b.actionFunc = map[int]JobActionFunc{
		1: mockActionFunc,
	}
	ch1 := make(chan interface{})
	ch2 := make(chan interface{})
	jobs := []jobQueueEntry{
		createMockJobQueueEntry(1, 1, ch1),
		createMockJobQueueEntry(1, 2, ch2),
	}
	b.queue = jobQueue{
		jobs,
		2,
	}
	b.procJobQueue()
	result1 := <-ch1
	result2 := <-ch2
	assert.Equal(s.T(), 2, result1.(int))
	assert.Equal(s.T(), 3, result2.(int))
}
