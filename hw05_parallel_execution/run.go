package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if len(tasks) == 0 {
		return nil
	}

	// Place your code here.
	var taskChannel = make(chan Task, n)
	var errorChannel = make(chan error)

	waitGroup := sync.WaitGroup{}

	for i := 0; i < n; i++ {
		go doTask(&taskChannel, &errorChannel, &waitGroup)
	}

	foundErrors := 0
	once := sync.Once{}
	for len(tasks) > 0 {
		select {
		case taskChannel <- tasks[0]:
			waitGroup.Add(1)
			tasks = tasks[1:]
		case <-errorChannel:
			foundErrors++
			if foundErrors >= m {
				once.Do(func() {
					close(taskChannel)
					close(errorChannel)
				})
				return ErrErrorsLimitExceeded
			}
		}
	}
	once.Do(func() {
		close(taskChannel)
		close(errorChannel)
	})
	waitGroup.Wait()
	return nil
}

func doTask(taskChannel *chan Task, errorChannel *chan error, group *sync.WaitGroup) {
	for {
		task, ok := <-*taskChannel
		if !ok {
			return
		}

		err := task()
		if err != nil {
			*errorChannel <- err
		}
		group.Done()
	}
}
