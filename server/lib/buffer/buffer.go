package buffer

import (
	"time"
)

const (
	SaveMDNote = 1
)

type Buffer struct {
	queue      jobQueue
	actionFunc map[int]JobActionFunc
	timeout    int
	maxProc    int
}

func NewBuffer(timeout, maxProc int) *Buffer {
	buffer := new(Buffer)
	buffer.queue = jobQueue{}
	buffer.timeout = timeout
	buffer.maxProc = maxProc
	buffer.actionFunc = buildActionFunc()
	return buffer
}

func buildActionFunc() map[int]JobActionFunc {
	return map[int]JobActionFunc{}
}

func (b *Buffer) NewJob(op int, payload interface{}, doneCh chan interface{}) {
	b.queue.Enqueue(op, payload, doneCh)
}

func (b *Buffer) JobCount() int {
	return b.queue.Size()
}

func (b *Buffer) StartProc() {
	for {
		jobs := b.queue.Dequeue(b.maxProc)
		for _, job := range jobs {
			go b.actionFunc[job.Op](job.Payload, job.DoneCh)
		}
		time.Sleep(time.Duration(b.timeout) * time.Millisecond)
	}
}
