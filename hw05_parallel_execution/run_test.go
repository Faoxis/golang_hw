package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		goroutinsCounter := atomic.Int32{}
		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				goroutinsCounter.Add(1)
				defer goroutinsCounter.Add(-1)

				time.Sleep(time.Millisecond)
				atomic.AddInt32(&runTasksCount, 1)

				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		resChan := make(chan error)
		defer close(resChan)

		go func() {
			err := Run(tasks, workersCount, maxErrorsCount)
			resChan <- err
		}()
		require.Eventually(t, func() bool {
			return int32(workersCount) == goroutinsCounter.Load()
		}, time.Second*2, time.Millisecond)

		err := <-resChan
		require.NoError(t, err)

		require.Equal(t, int32(tasksCount), runTasksCount, "not all tasks were completed")
	})

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("if task list is empty", func(t *testing.T) {
		err := Run([]Task{}, 10, 1)
		require.Nil(t, err)
	})

	t.Run("if goroutins more then tasks", func(t *testing.T) {
		taskCounter := 50
		goroutinCounter := 100

		tasks := make([]Task, 0, taskCounter)
		for i := 0; i < 0; i++ {
			tasks[i] = func() error {
				time.Sleep(time.Millisecond * 10)
				return nil
			}
		}
		err := Run(tasks, goroutinCounter, 1)
		require.Nil(t, err)
	})

	t.Run("if m is zero", func(t *testing.T) {
		err := Run([]Task{}, 1, 0)
		require.ErrorIs(t, ErrErrorsLimitExceeded, err)
	})

	t.Run("if m is negative", func(t *testing.T) {
		err := Run([]Task{}, 1, -42)
		require.ErrorIs(t, ErrErrorsLimitExceeded, err)
	})
}
