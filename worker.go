package pool

import (
	"log"
	"sync/atomic"
)

var (
	DefaultMaxJobPerWorker = 1024
)

const (
	STATUS_INIT    = iota
	STATUS_WORKING
	STATUS_IDLE
)

type Job interface {
	Do()         error
	Name()       string
	Error(error)
}

type Worker struct {
	id     uint
	task   chan Job
	stopCh chan struct{}
	d      *Dispatcher
	status uint32
}

func NewWorker(id uint, d *Dispatcher) *Worker {
	w := &Worker{
		id:     id,
		task:   make(chan Job, 100),
		stopCh: make(chan struct{}),
		d:      d,
	}
	atomic.StoreUint32(&w.status, STATUS_INIT)
	w.Start()
	return w
}

func (w *Worker) Start() {
	atomic.StoreUint32(&w.status, STATUS_IDLE)
	go func() {
		defer w.d.wg.Done()
		for {
			select {
			case job, ok := <-w.task:
				if !ok {
					return
				}
				log.Printf("Worker %d start do job: %s\n", w.id, job.Name())
				atomic.StoreUint32(&w.status, uint32(STATUS_WORKING))
				err := job.Do()
				if err != nil {
					log.Printf("Worker %d do %s job error: %s\n", w.id, job.Name(), err.Error())
					job.Error(err)
				}
				atomic.StoreUint32(&w.status, uint32(STATUS_IDLE))
			case <-w.stopCh:
				log.Printf("Worker %d closing . . .\n", w.id)
				return
			}
		}
	}()
}

func (w *Worker) SendJob(job Job) {
	w.task <- job
}

func (w *Worker) Status() uint32 {
	return w.status
}

func (w *Worker) Close() {
	close(w.stopCh)
}