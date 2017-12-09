package buffer

import (
	"errors"
	"fmt"
	"time"
)

type jobQueueEntry struct {
	Op        int
	Payload   interface{}
	DoneCh    chan interface{}
	CreatedAt time.Time
}

type jobQueue struct {
	storage []jobQueueEntry
	size    int
}

type JobActionFunc func(interface{}, chan interface{})

func (q *jobQueue) Size() int {
	return q.size
}

func (q *jobQueue) Enqueue(op int, payload interface{}, doneCh chan interface{}) {
	q.size++
	q.storage = append(q.storage, jobQueueEntry{
		op,
		payload,
		doneCh,
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
