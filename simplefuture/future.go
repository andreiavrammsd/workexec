// Package simplefuture represents an async Task call.
package simplefuture

import (
	"errors"
	"sync"
)

// ErrNilTask is returned when a nil task is provided.
var ErrNilTask = errors.New("nil task passed")

// Future represents a task which executes async work.
type Future struct {
	sync.RWMutex

	task     Task
	done     chan struct{}
	err      error
	canceled bool
	on       bool
}

// Run executes the Task async.
func (f *Future) Run() {
	f.Lock()

	if f.on {
		f.Unlock()

		return
	}

	f.on = true
	f.Unlock()

	go f.run()
}

// Wait blocks until task is done.
func (f *Future) Wait() {
	f.RLock()

	if !f.on {
		f.RUnlock()

		return
	}

	f.RUnlock()

	<-f.done
}

func (f *Future) Error() error {
	f.RLock()
	defer f.RUnlock()

	return f.err
}

// Cancel asks the task to stop.
func (f *Future) Cancel() {
	f.Lock()
	f.canceled = true
	f.Unlock()
}

// IsCanceled returns true if task was canceled.
func (f *Future) IsCanceled() bool {
	f.RLock()
	defer f.RUnlock()

	return f.canceled
}

func (f *Future) run() {
	defer close(f.done)

	f.err = f.task.Run(f.IsCanceled)

	if task, ok := f.task.(CanceledTask); ok {
		f.RLock()

		if f.canceled {
			f.RUnlock()

			task.OnCancel()

			return
		}

		f.RUnlock()
	}

	if task, ok := f.task.(SuccessfulTask); ok {
		if f.err == nil {
			task.OnSuccess()
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
		return nil, ErrNilTask
	}

	future := &Future{
		task: task,
		done: make(chan struct{}),
	}

	return future, nil
}
