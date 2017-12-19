package buffer

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Implementation of a job queue to handle bursts of client requests without timing out

type Buffer struct {
	queue   jobQueue
	timeout int
	maxProc int
}

// Initializes a new buffer instance
// timeout - the time interval that a batch of jobs will be processed
// maxProc - the max number of jobs to be processed in a batch
func NewBuffer(timeout, maxProc int) *Buffer {
	buffer := new(Buffer)
	buffer.queue = jobQueue{}
	buffer.timeout = timeout
	buffer.maxProc = maxProc
	return buffer
}

// Creates a new job on the buffer
// op - the action function that will be called when the job runs
// payload - the parameter that the action function takes in
// doneCh - the response that the action function returns when it finishes processing
// errCh - the response that the action function returns when an error occurs
func (b *Buffer) NewJob(op JobActionFunc, payload interface{}, doneCh chan interface{}, errCh chan error) {
	b.queue.Enqueue(op, payload, doneCh, errCh)
}

// Returns the number of unprocessed jobs currently in the buffer
func (b *Buffer) JobCount() int {
	return b.queue.Size()
}

// Start processing jobs in the buffer
func (b *Buffer) StartProc() {
	for {
		b.procJobQueue()
		time.Sleep(time.Duration(b.timeout) * time.Millisecond)
	}
}

// Processes a batch of jobs
// The number of jobs in the batch is determined by the maxProcs attribute
func (b *Buffer) procJobQueue() {
	if jobs, err := b.queue.Dequeue(b.maxProc); err != nil {
		log.WithField("err", err).Error("Dequeue failed")
	} else {
		for _, job := range jobs {
			go job.Op(job.Payload, job.DoneCh, job.ErrCh)
		}
	}
}
