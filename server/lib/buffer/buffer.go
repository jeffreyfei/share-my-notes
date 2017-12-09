package buffer

import (
	"time"

	log "github.com/sirupsen/logrus"
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
		b.procJobQueue()
	}
}

func (b *Buffer) procJobQueue() {
	if jobs, err := b.queue.Dequeue(b.maxProc); err != nil {
		log.WithField("err", err).Error("Dequeue failed")
	} else {
		for _, job := range jobs {
			go b.actionFunc[job.Op](job.Payload, job.DoneCh)
		}
	}
	time.Sleep(time.Duration(b.timeout) * time.Millisecond)
}
