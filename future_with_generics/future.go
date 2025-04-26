// Package future_with_generics represents an async Task call.
package future_with_generics

import (
	"errors"
	"sync"
)

// Future represents a task which executes async work.
type Future[T any] struct {
	task     Task[T]
	done     chan struct{}
	result   T
	err      error
	canceled bool
	on       bool
	lock     sync.RWMutex
}

// Run executes the Task async.
func (f *Future[T]) Run() {
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
func (f *Future[T]) Wait() {
	f.lock.RLock()
	if !f.on {
		f.lock.RUnlock()
		return
	}
	f.lock.RUnlock()

	<-f.done
}

// Result returns task result and error. Blocks until task is done.
func (f *Future[T]) Result() (T, error) {
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
func (f *Future[T]) Cancel() {
	f.lock.Lock()
	f.canceled = true
	f.lock.Unlock()
}

// IsCanceled returns true if task was canceled.
func (f *Future[T]) IsCanceled() bool {
	f.lock.RLock()
	defer f.lock.RUnlock()
	return f.canceled
}

func (f *Future[T]) run() {
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

	if task, ok := f.task.(SuccessfulTask[T]); ok {
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
func New[T any](task Task[T]) (*Future[T], error) {
	if task == nil {
		return nil, errors.New("nil task passed")
	}

	future := &Future[T]{
		task: task,
		done: make(chan struct{}),
	}

	return future, nil
}
