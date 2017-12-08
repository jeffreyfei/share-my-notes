package buffer

import "time"

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
	q.storage = append(q.storage, jobQueueEntry{
		op,
		payload,
		doneCh,
		time.Now(),
	})
}

func (q *jobQueue) Dequeue(maxEntry int) []jobQueueEntry {
	if maxEntry <= len(q.storage) {
		result := q.storage
		q.storage = []jobQueueEntry{}
		return result
	}
	result := q.storage[0:maxEntry]
	q.storage = q.storage[maxEntry-1 : len(q.storage)]
	return result
}
