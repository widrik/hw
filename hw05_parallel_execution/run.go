package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
var ErrRoutinesCount = errors.New("routines count should be more than 0")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving M errors from tasks.
func Run(tasks []Task, n int, m int) error {
	err := validateParams(n, m)

	if err != nil {
		return err
	}

	var errorsCount = 0
	var wg = sync.WaitGroup{}
	var mut = sync.Mutex{}
	var tasksCount = len(tasks)

	if n > tasksCount {
		n = tasksCount
	}

	wg.Add(n)

	tasksCh := make(chan Task, tasksCount)
	for _, task := range tasks {
		tasksCh <- task
	}
	close(tasksCh)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()

			for {
				task, ok := <-tasksCh

				if !ok {
					break
				}

				err := task()
				mut.Lock()

				if err != nil {
					errorsCount++
				}

				if errorsCount == m {
					mut.Unlock()
					return
				}
				mut.Unlock()
			}
		}()
	}

	wg.Wait()

	if errorsCount >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func validateParams(n int, m int) error {
	var err error

	if n <= 0 {
		err = ErrRoutinesCount
	}

	if m <= 0 {
		err = ErrErrorsLimitExceeded
	}

	return err
}
