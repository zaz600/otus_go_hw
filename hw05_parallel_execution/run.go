package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"context"
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
var ErrInvalidWorkerCount = errors.New("num of workers must be >= 1")

type Task func() error

type errCounter struct {
	limit int

	mu      *sync.Mutex
	current int
}

func (e *errCounter) inc() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.current++
}

func (e *errCounter) isLimitExceeded() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.limit > 0 && e.current > e.limit
}

func newErrCounter(limit int) *errCounter {
	return &errCounter{
		limit: limit,
		mu:    &sync.Mutex{},
	}
}

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks.
func Run(tasks []Task, n int, m int) error {
	if n <= 0 {
		return ErrInvalidWorkerCount
	}

	wg := sync.WaitGroup{}
	counter := newErrCounter(m)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	taskCh := make(chan Task, n)

	worker := func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return

			case task, ok := <-taskCh:
				if ok {
					if err := task(); err != nil {
						counter.inc()
					}
					if counter.isLimitExceeded() {
						cancel()
						return
					}
				} else {
					return
				}
			}
		}
	}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker()
	}

	go func() {
		defer close(taskCh)
		for _, task := range tasks {
			select {
			case <-ctx.Done():
				return
			case taskCh <- task:
			}
		}
	}()

	wg.Wait()

	if counter.isLimitExceeded() {
		return ErrErrorsLimitExceeded
	}
	return nil
}
