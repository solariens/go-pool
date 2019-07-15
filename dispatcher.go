package pool

import (
	"sync"
	"math/rand"
	"sync/atomic"
)

type Dispatcher struct {
	workerList      []*Worker
	workerNum       uint
	stopCh          chan struct{}
	jobs            chan Job
	maxJobNum       uint
	wg              sync.WaitGroup
	maxJobPerWorker uint
}

func NewDispatcher(cfg Config) *Dispatcher {
	d := &Dispatcher{}
	for i:=0; i<int(cfg.WorkerNum); i++ {
		d.wg.Add(1)
		w := NewWorker(uint(i), cfg.MaxJobPerWorker, d)
		d.workerList = append(d.workerList, w)
	}

	if cfg.MaxJobNum == 0 {
		cfg.MaxJobNum = cfg.WorkerNum * cfg.MaxJobPerWorker
	}
	d.workerNum = cfg.WorkerNum
	d.stopCh    = make(chan struct{})
	d.jobs      = make(chan Job, cfg.MaxJobNum)
	return d
}

func (d *Dispatcher) findIdleWorker() (*Worker, bool) {
	var ret *Worker
	for _, worker := range d.workerList {
		status := uint32(worker.Status())
		if atomic.LoadUint32(&status) == STATUS_IDLE {
			ret = worker
			break
		}
	}
	return ret, ret != nil
}

func (d *Dispatcher) Start() {
	go func() {
		for {
			select {
			case <-d.stopCh:
				return
			case job, ok := <-d.jobs:
				if !ok {
					return
				}
				w, b := d.findIdleWorker()
				if b {
					w.SendJob(job)
				} else {
					index := uint(rand.Int31()) % d.workerNum
					d.workerList[index].SendJob(job)
				}
			}
		}
	}()
}

func (d *Dispatcher) SendJob(job Job) {
	d.jobs <- job
}

func (d *Dispatcher) Stop() {
	for _, worker := range d.workerList {
		worker.Close()
	}
	d.wg.Wait()
	close(d.stopCh)
}
