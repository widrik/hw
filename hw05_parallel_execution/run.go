package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
var ErrRoutinesCount = errors.New("routines count should be more than 0")

type Task func() error

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks
func Run(tasks []Task, N int, M int) error {
	err := validateParams(N, M)

	if err != nil {
		return err
	}

	var errorsCount = 0
	var wg = sync.WaitGroup{}
	var m = sync.Mutex{}
	var tasksCount = len(tasks)

	if N > tasksCount {
		N = tasksCount
	}

	wg.Add(N)

	tasksCh := make(chan Task, tasksCount)
	for _, task := range tasks {
		tasksCh <- task
	}
	close(tasksCh)

	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()

			for {
				task, ok := <-tasksCh

				if !ok {
					break
				}

				err := task()
				m.Lock()

				if err != nil {
					errorsCount++
				}

				if errorsCount == M {
					m.Unlock()
					return
				}
				m.Unlock()
			}

		}()
	}

	wg.Wait()

	if errorsCount >= M {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func validateParams(N int, M int) error {
	var err error

	if N <= 0 {
		err = ErrRoutinesCount
	}

	if M <= 0 {
		err = ErrErrorsLimitExceeded
	}

	return err
}
