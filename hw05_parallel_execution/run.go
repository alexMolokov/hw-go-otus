package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type SyncRun struct {
	countRunning chan int
	countErrors  *int
	sync.WaitGroup
	sync.RWMutex
}

func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	if n <= 0 {
		return nil
	}

	var countErrors int
	sr := &SyncRun{
		countRunning: make(chan int, n),
		countErrors:  &countErrors,
	}

	defer close(sr.countRunning)

	for _, task := range tasks {
		sr.countRunning <- 1
		sr.Add(1)
		go func(task Task, sr *SyncRun) {
			err := task()
			if err != nil {
				sr.Lock()
				*sr.countErrors++
				sr.Unlock()
			}
			<-sr.countRunning
			sr.Done()
		}(task, sr)

		sr.RLock()
		if countErrors >= m {
			sr.RUnlock()
			break
		}
		sr.RUnlock()
	}

	sr.Wait()

	if countErrors < m {
		return nil
	}

	return ErrErrorsLimitExceeded
}
