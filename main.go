package main

import (
	"fmt"
	"runtime"
	"sync"
)

type Worker struct {
	size int
	jobs chan func() error
	wg   sync.WaitGroup
	err  error
}

func NewWorker(size int) *Worker {
	return &Worker{
		size: size,
		jobs: make(chan func() error),
		wg:   sync.WaitGroup{},
		err:  nil,
	}
}

func (w *Worker) worker() {
	for j := range w.jobs {
		if w.err == nil {
			if err := j(); err != nil {
				w.err = err
			}
		}
		w.wg.Done()
	}
}

func (w *Worker) Exec(job func() error) {
	w.wg.Add(1)
	w.jobs <- job
}

func (w *Worker) Start() {
	for i := 0; i < w.size; i++ {
		go w.worker()
	}
}

func (w *Worker) Wait() {
	w.wg.Wait()
}

func (w *Worker) Stop() {
	w.Wait()
	close(w.jobs)
}

func (w *Worker) Error() error {
	return w.err
}

func main() {
	// dummy
	jobs := []uint64{1, 2, 3}
	maxWorkerSize := runtime.NumCPU() - 1

	w := NewWorker(maxWorkerSize)
	w.Start()

	for _, job := range jobs {
		jobFunc := func() error {
			fmt.Println(job)
			return nil
		}
		w.Exec(jobFunc)
	}
	w.Stop()

	if w.Error() != nil {
		panic(w.Error())
	}
}