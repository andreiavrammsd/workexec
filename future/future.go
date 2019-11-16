// Package future represents an async Task call.
package future

import (
	"errors"
	"sync"
)

// Future represents a task which executes async work.
type Future struct {
	task     Task
	done     chan struct{}
	result   interface{}
	err      error
	canceled bool
	on       bool
	lock     sync.RWMutex
}

// Run executes the Task async.
func (f *Future) Run() {
	f.lock.Lock()
	if f.on {
		f.lock.Unlock()
		return
	}
	f.on = true
	f.lock.Unlock()

	go f.run()
}

// Wait blocks until task is done.
func (f *Future) Wait() {
	f.lock.RLock()
	if !f.on {
		f.lock.RUnlock()
		return
	}
	f.lock.RUnlock()

	<-f.done
}

// Result returns task result and error. Blocks until task is done.
func (f *Future) Result() (interface{}, error) {
	f.lock.RLock()
	if !f.on {
		f.lock.RUnlock()
		return f.result, f.err
	}
	f.lock.RUnlock()

	<-f.done

	return f.result, f.err
}

// Cancel asks the task to stop.
func (f *Future) Cancel() {
	f.lock.Lock()
	f.canceled = true
	f.lock.Unlock()
}

// IsCanceled returns true if task was canceled.
func (f *Future) IsCanceled() bool {
	f.lock.RLock()
	defer f.lock.RUnlock()
	return f.canceled
}

func (f *Future) run() {
	defer close(f.done)

	f.result, f.err = f.task.Run(f.IsCanceled)

	if task, ok := f.task.(CanceledTask); ok {
		f.lock.RLock()
		if f.canceled {
			f.lock.RUnlock()
			task.OnCancel()
			return
		}
		f.lock.RUnlock()
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

	future := &Future{
		task: task,
		done: make(chan struct{}),
	}

	return future, nil
}
