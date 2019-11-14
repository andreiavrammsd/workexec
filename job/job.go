package job

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

// Task belongs to user
type Task interface {
	Run(*Job) (interface{}, error)
}

// ID of a job
type ID string

type (
	OnSuccess func(interface{})
	OnError   func(error)
	OnCancel  func(error)
)

// Job contains a Task
type Job struct {
	id        uuid.UUID
	task      Task
	err       error
	cancel    error
	onSuccess OnSuccess
	onError   OnError
	onCancel  OnCancel
	lock      sync.RWMutex
}

// ID returns job ID
func (j *Job) ID() ID {
	return ID(j.id.String())
}

func (j *Job) run() (interface{}, error) {
	result, err := j.task.Run(j)

	j.lock.Lock()
	defer j.lock.Unlock()

	// Cancel
	if j.cancel != nil {
		return nil, nil
	}

	// Error or success
	if err != nil {
		if j.onError != nil {
			j.err = err
			j.onError(err)
		}
	} else {
		if j.onSuccess != nil {
			j.onSuccess(result)
		}
	}

	return result, err
}

// Run starts executing the job task and returns a Future.
func (j *Job) Run() *Future {
	wait := &Future{
		done: make(chan struct{}),
	}

	go func() {
		defer close(wait.done)
		wait.result, wait.err = j.run()
		wait.canceled = j.IsCanceled()
	}()

	return wait
}

// Cancel stops job.
func (j *Job) Cancel(err error) {
	if err == nil {
		err = errors.New("canceled")
	}

	j.lock.Lock()
	defer j.lock.Unlock()

	j.cancel = err

	if j.onCancel != nil {
		j.onCancel(j.cancel)
	}
}

func (j *Job) OnSuccess(f OnSuccess) {
	j.onSuccess = f
}

func (j *Job) OnError(f OnError) {
	j.lock.Lock()
	j.onError = f
	j.lock.Unlock()
}

func (j *Job) OnCancel(f OnCancel) {
	j.onCancel = f
}

// IsCanceled returns true if job was canceled.
func (j *Job) IsCanceled() bool {
	j.lock.RLock()
	defer j.lock.RUnlock()
	return j.cancel != nil
}

// New creates a new job with a given task.
func New(task Task) (*Job, error) {
	if task == nil {
		return nil, errors.New("task function not passed")
	}

	job := &Job{
		task: task,
		id:   uuid.New(),
	}

	return job, nil
}
