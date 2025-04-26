// Package taskexecutor handles incoming tasks to be executed.
package taskexecutor_with_generics

import (
	"errors"
	"math"
	"runtime"
	"sync"
)

const (
	defaultQueueSize = 1024
)

// Config allows setup of executor.
type Config struct {
	// Concurrency is the number of routines the executor will start working on.
	Concurrency uint

	// QueueSize is the number of tasks accepted before blocking.
	QueueSize uint
}

// TaskExecutor represents the executor instance.
type TaskExecutor struct {
	concurrency  uint
	queue        chan Future[any]
	wait         chan struct{}
	stop         chan struct{}
	runningTasks uint
	lock         sync.RWMutex
	stopped      bool
}

// Start opens the working routines.
func (te *TaskExecutor) Start() {
	for i := uint(0); i < te.concurrency; i++ {
		go te.run()
	}
}

// Stop asks the working routines to stop after tasks are finished.
func (te *TaskExecutor) Stop() {
	te.lock.Lock()
	if te.stopped {
		te.lock.Unlock()
		return
	}
	te.stopped = true
	te.lock.Unlock()

	for i := uint(0); i < te.concurrency; i++ {
		te.stop <- struct{}{}
	}
}

// Wait blocks until executor is done with running all the queued tasks.
func (te *TaskExecutor) Wait() {
	<-te.wait
}

// Submit puts a task into the executor queue.
func (te *TaskExecutor) Submit(future Future[any]) error {
	te.lock.RLock()
	if te.stopped {
		te.lock.RUnlock()
		return errors.New("executor is stopped")
	}
	te.lock.RUnlock()

	te.queue <- future

	return nil
}

func (te *TaskExecutor) run() {
	for {
		select {
		case future := <-te.queue:
			te.lock.Lock()
			te.runningTasks++
			te.lock.Unlock()

			future.Run()
			future.Wait()

			te.lock.Lock()
			te.runningTasks--
			te.lock.Unlock()
		case <-te.stop:
			te.lock.RLock()
			if te.stopped && te.runningTasks == 0 {
				te.wait <- struct{}{}
			}
			te.lock.RUnlock()

			return
		}
	}
}

// New creates a new task executor.
func New(c Config) *TaskExecutor {
	if c.Concurrency == 0 {
		c.Concurrency = uint(math.Max(1, float64(runtime.NumCPU())-1))
	}
	if c.QueueSize == 0 {
		c.QueueSize = defaultQueueSize
	}

	return &TaskExecutor{
		concurrency: c.Concurrency,
		queue:       make(chan Future[any], c.QueueSize),
		wait:        make(chan struct{}, c.Concurrency),
		stop:        make(chan struct{}, c.Concurrency),
	}
}
