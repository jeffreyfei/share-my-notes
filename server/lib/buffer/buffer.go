package buffer

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type Buffer struct {
	queue   jobQueue
	timeout int
	maxProc int
}

func NewBuffer(timeout, maxProc int) *Buffer {
	buffer := new(Buffer)
	buffer.queue = jobQueue{}
	buffer.timeout = timeout
	buffer.maxProc = maxProc
	return buffer
}

func (b *Buffer) NewJob(op JobActionFunc, payload interface{}, doneCh chan interface{}, errCh chan error) {
	b.queue.Enqueue(op, payload, doneCh, errCh)
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
			go job.Op(job.Payload, job.DoneCh, job.ErrCh)
		}
	}
	time.Sleep(time.Duration(b.timeout) * time.Millisecond)
}
