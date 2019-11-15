// Package job provides a method to executed an async task and be notified when executions
// is done (successfully, with error or canceled).
// A job contains a task which actually does the work. A task is a struct which implements the
// Task interface. Other task related interfaces can be used to be notified about events.
package job

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

// ID of a job.
type ID string

// Job contains a Task.
type Job struct {
	id     uuid.UUID
	task   Task
	cancel error
	lock   sync.RWMutex
}

// ID returns the job unique identifier.
func (j *Job) ID() ID {
	return ID(j.id.String())
}

// Run starts executing the job task and returns a Future.
func (j *Job) Run() *Future {
	future := &Future{
		done: make(chan struct{}),
	}

	go func() {
		defer close(future.done)
		future.result, future.err = j.run()
		future.canceled = j.IsCanceled()
	}()

	return future
}

// Cancel asks the job to stop.
func (j *Job) Cancel(err error) {
	if _, ok := j.task.(CancelableTask); !ok {
		return
	}

	if err == nil {
		err = errors.New("canceled")
	}

	j.lock.Lock()
	j.cancel = err
	j.lock.Unlock()
}

// IsCanceled returns true if job was canceled.
func (j *Job) IsCanceled() bool {
	j.lock.RLock()
	defer j.lock.RUnlock()
	return j.cancel != nil
}

func (j *Job) run() (result interface{}, err error) {
	result, err = j.task.Run(j)

	if task, ok := j.task.(CancelableTask); ok {
		j.lock.RLock()
		cancel := j.cancel
		j.lock.RUnlock()

		if cancel != nil {
			task.OnCancel(j.cancel)
			return
		}
	}

	if task, ok := j.task.(SuccessfulTask); ok {
		if err == nil {
			task.OnSuccess(result)
		}
	}

	if task, ok := j.task.(FailedTask); ok {
		if err != nil {
			task.OnError(err)
		}
	}

	return
}

// New creates a new job with a given task.
func New(task Task) (*Job, error) {
	if task == nil {
		return nil, errors.New("nil task passed to job")
	}

	job := &Job{
		task: task,
		id:   uuid.New(),
	}

	return job, nil
}
