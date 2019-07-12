package main

import (
	"time"
	"log"
	"github.com/soalriens/pool"
)

type Task struct {

}

func (t *Task) Do() error {
	time.Sleep(1 * time.Second)
	return nil
}

func (t *Task) Name() string {
	return "TestTask"
}

func (t *Task) Error(err error) {
	log.Println(err)
}

func main() {
	d := pool.NewDispatcher(pool.Config{WorkerNum: 100, MaxJobPerWorker: 1024})

	d.Start()

	d.SendJob(&Task{})

	time.Sleep(5 * time.Second)

	d.Stop()
}
