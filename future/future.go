// Package future represents an async Task call.
package future

import (
	"errors"
	"sync"
)

// Future represents a task which executes work.
type Future struct {
	task     Task
	result   interface{}
	err      error
	on       bool
	canceled bool
	sync.RWMutex
}

// Wait blocks until task is done.
func (f *Future) Wait() {
	if f.on {
		return
	}

	done := make(chan struct{})

	go func() {
		defer close(done)
		f.run()
	}()

	<-done
}

// Cancel asks the task to stop.
func (f *Future) Cancel() {
	f.Lock()
	f.canceled = true
	f.Unlock()
}

// Result returns task result and error. Blocks until task is done.
func (f *Future) Result() (interface{}, error) {
	if !f.on {
		f.run()
	}

	return f.result, f.err
}

// IsCanceled returns true if task was canceled.
func (f *Future) IsCanceled() bool {
	f.RLock()
	defer f.RUnlock()
	return f.canceled
}

func (f *Future) run() {
	f.on = true

	f.result, f.err = f.task.Run(f.IsCanceled)

	if task, ok := f.task.(CanceledTask); ok {
		if f.canceled {
			task.OnCancel()
			return
		}
	}

	if task, ok := f.task.(SuccessfulTask); ok {
		if f.err == nil {
			task.OnSuccess(f.result)
		}
	}

	if task, ok := f.task.(FailedTask); ok {
		if f.err != nil {
			task.OnError(f.err)
		}
	}
}

// New creates a new future with a given task.
func New(task Task) (*Future, error) {
	if task == nil {
		return nil, errors.New("nil task passed")
	}

	return &Future{
		task: task,
	}, nil
}
