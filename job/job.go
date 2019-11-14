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

// Job contains a Task
type ID string

type (
	OnSuccess func(interface{})
	OnError   func(error)
	OnCancel  func(error)
)

type Job struct {
	id        uuid.UUID
	task      Task
	err       error
	cancel    error
	onSuccess OnSuccess
	onError   OnError
	onCancel  OnCancel
	lock      sync.RWMutex
	options   []Option
}

func (j *Job) ID() ID {
	return ID(j.id.String())
}

func (j *Job) Run() (interface{}, error) {
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

func (j *Job) Go() *Wait {
	wait := &Wait{
		done: make(chan struct{}),
	}

	go func() {
		defer close(wait.done)
		wait.result, wait.err = j.Run()
		wait.canceled = j.IsCanceled()
	}()

	return wait
}

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

func (j *Job) IsCanceled() bool {
	j.lock.RLock()
	defer j.lock.RUnlock()
	return j.cancel != nil
}

func New(task Task, o ...Option) (*Job, error) {
	if task == nil {
		return nil, errors.New("task function not passed")
	}

	job := &Job{
		task: task,
		id:   uuid.New(),
	}

	for i := 0; i < len(o); i++ {
		o[i].Apply(job)
	}

	return job, nil
}
