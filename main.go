package main

import (
	"fmt"
	"log"
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

func (w *Worker) Close() {
	w.Wait()
	close(w.jobs)
}

func (w *Worker) Error() error {
	return w.err
}

func main() {
	jobNums := []int{1, 2, 3, 4, 5}
	maxWorkerSize := runtime.NumCPU() - 1
	fmt.Printf("run on %d threads.\n", maxWorkerSize)

	w := NewWorker(maxWorkerSize)
	w.Start()
	for _, jobNum := range jobNums {
		w.Exec(jobFunc(jobNum))
	}
	w.Close()

	if w.Error() != nil {
		panic(w.Error())
	}
}

func jobFunc(jobNum int) func() error {
	return func() error {
		log.Printf("num: %d \n", jobNum)
		return nil
	}
}