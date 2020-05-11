package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return fmt.Errorf("error from task %d", i)
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		result := Run(tasks, workersCount, maxErrorsCount)

		require.Equal(t, ErrErrorsLimitExceeded, result)
		require.LessOrEqual(t,
			int32(workersCount+maxErrorsCount), runTasksCount, "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		result := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.Nil(t, result)

		require.Equal(t, int32(tasksCount), runTasksCount, "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("wrong goroutines count", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		maxErrorsCount := 10

		workersCount := 0
		result1 := Run(tasks, workersCount, maxErrorsCount)
		require.Equal(t, ErrRoutinesCount, result1)

		workersCount = -10
		result2 := Run(tasks, workersCount, maxErrorsCount)
		require.Equal(t, ErrRoutinesCount, result2)
	})

	t.Run("tasks with errors less then M", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		maxErrorsCount := 25
		errorsCount := 0
		for i := 0; i < tasksCount; i++ {
			if errorsCount < (maxErrorsCount - 5) {
				tasks = append(tasks, func() error {
					atomic.AddInt32(&runTasksCount, 1)
					return fmt.Errorf("error from task %d", i)
				})
			} else {
				tasks = append(tasks, func() error {
					atomic.AddInt32(&runTasksCount, 1)
					return nil
				})
			}
			errorsCount++
		}

		workersCount := 10

		result := Run(tasks, workersCount, maxErrorsCount)

		require.NotEqual(t, ErrErrorsLimitExceeded, result)
		require.LessOrEqual(t, runTasksCount, int32(tasksCount), "all tasks were started")
	})
}
