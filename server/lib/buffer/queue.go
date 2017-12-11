package buffer

import (
	"errors"
	"fmt"
	"time"
)

type jobQueueEntry struct {
	Op        JobActionFunc
	Payload   interface{}
	DoneCh    chan interface{}
	ErrCh     chan error
	CreatedAt time.Time
}

type jobQueue struct {
	storage []jobQueueEntry
	size    int
}

type JobActionFunc func(interface{}, chan interface{}, chan error)

func (q *jobQueue) Size() int {
	return q.size
}

func (q *jobQueue) Enqueue(op JobActionFunc, payload interface{}, doneCh chan interface{}, errCh chan error) {
	q.size++
	q.storage = append(q.storage, jobQueueEntry{
		op,
		payload,
		doneCh,
		errCh,
		time.Now(),
	})
}

func (q *jobQueue) Dequeue(limit int) ([]jobQueueEntry, error) {
	if limit <= 0 {
		return []jobQueueEntry{}, errors.New(fmt.Sprintf("Invalid limit: %d", limit))
	}
	if limit == 0 {
		return []jobQueueEntry{}, nil
	}
	if limit >= len(q.storage) {
		q.size = 0
		result := q.storage
		q.storage = []jobQueueEntry{}
		return result, nil
	}
	q.size -= limit
	result := q.storage[0:limit]
	q.storage = q.storage[limit:len(q.storage)]
	return result, nil
}
