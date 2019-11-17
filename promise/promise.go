// Package promise implements an asynchronous operation of an executor.
// An executor is a struct which implements the Executor interface.
// The New function takes one or multiple executors and returns a Promise.
// The Promise can be used to attach a success callback (Then) and/or an error callback (Error).
// Async starts calling the executors without blocking.
// Await blocks calling routine until the executors finish.
// If one of the executors returns an error, the others will not be called.
package promise

import (
	"sync"
)

// Executor interface must be implemented to be called inside a Promise.
type Executor interface {
	// Execute is the called method when a Promise starts.
	Execute() error
}

// Then is the function called after Execute returns with no error.
type Then func()

// Error is the function called after Execute returns with error. The error is passed as argument.
type Error func(error)

// Promise represents an async Executor execution.
type Promise struct {
	executors []Executor
	err       error
	then      Then
	error     Error
	lock      sync.RWMutex
}

// Then sets the success callback.
func (p *Promise) Then(f Then) *Promise {
	p.lock.Lock()
	p.then = f
	p.lock.Unlock()
	return p
}

// Error sets the fail callback.
func (p *Promise) Error(f Error) *Promise {
	p.lock.Lock()
	p.error = f
	p.lock.Unlock()
	return p
}

// Async executes executors asynchronous.
func (p *Promise) Async() {
	go func() {
		p.exec()
	}()
}

// Await blocks until the executors finish and returns the error.
func (p *Promise) Await() error {
	p.exec()
	return p.err
}

// New creates a Promise with given executors.
func New(e ...Executor) *Promise {
	return &Promise{
		executors: e,
	}
}

func (p *Promise) exec() {
	wg := sync.WaitGroup{}

	for i := 0; i < len(p.executors); i++ {
		if p.executors[i] == nil {
			continue
		}

		wg.Add(1)

		go func(executor Executor) {
			defer wg.Done()

			// Stop executors on first error
			p.lock.RLock()
			stop := p.err != nil
			p.lock.RUnlock()
			if stop {
				return
			}

			// Call executor
			if err := executor.Execute(); err != nil {
				p.lock.Lock()
				p.err = err
				p.lock.Unlock()

				if p.error != nil {
					p.error(err)
				}
			} else if p.then != nil {
				p.then()
			}
		}(p.executors[i])
	}

	wg.Wait()
}
