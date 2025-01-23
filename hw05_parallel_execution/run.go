package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}
	// Создаем канал для задач
	taskChannel := make(chan Task, n)

	// Создаем канал для сигнала завершения
	doneChannel := make(chan struct{})
	onceForDoneChannel := sync.Once{}
	closeDoneChannel := func() {
		onceForDoneChannel.Do(func() {
			close(doneChannel)
		})
	}

	wg := sync.WaitGroup{}
	errorCounter := atomic.Int64{}

	worker := func() {
		defer wg.Done()
		for {
			select {
			case <-doneChannel:
				return
			case task, ok := <-taskChannel:
				if !ok {
					return
				}

				err := task()
				if err != nil {
					errorCounter.Add(1)
					if errorCounter.Load() >= int64(m) {
						closeDoneChannel()
						return
					}
				}
			}
		}
	}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker()
	}

	go func() {
		defer close(taskChannel)
		i := 0
		for i < len(tasks) && errorCounter.Load() < int64(m) {
			select {
			case <-doneChannel:
				return
			default:
			}

			select {
			case <-doneChannel:
				return
			case taskChannel <- tasks[i]:
			}
			i++
		}
	}()

	wg.Wait()
	closeDoneChannel()
	if errorCounter.Load() >= int64(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}
