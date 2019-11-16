package simplefuture_test

import (
	"errors"
	"testing"

	"github.com/andreiavrammsd/jobrunner/simplefuture"

	"github.com/stretchr/testify/assert"
)

func TestNew_WithNilTaskError(t *testing.T) {
	taskFuture, err := simplefuture.New(nil)
	assert.Nil(t, taskFuture)
	assert.Error(t, err)
}

func TestFuture_Wait(t *testing.T) {
	task := &divideTask{a: 4, b: 2}
	taskFuture, err := simplefuture.New(task)
	assert.NoError(t, err)

	taskFuture.Run()
	taskFuture.Wait()
	taskFuture.Wait()

	assert.Equal(t, float64(2), task.result)
	assert.NoError(t, taskFuture.Error())
	assert.False(t, taskFuture.IsCanceled())
}

func TestFuture_Error(t *testing.T) {
	task := &divideTask{a: 4, b: 0}
	taskFuture, err := simplefuture.New(task)
	assert.NoError(t, err)

	taskFuture.Run()
	taskFuture.Wait()

	assert.Equal(t, float64(0), task.result)
	assert.Error(t, taskFuture.Error())
	assert.False(t, taskFuture.IsCanceled())
}

func TestFuture_Cancel(t *testing.T) {
	task := &longRunningTask{}
	taskFuture, err := simplefuture.New(task)
	assert.NoError(t, err)

	taskFuture.Run()
	taskFuture.Cancel()
	taskFuture.Wait()

	assert.NoError(t, taskFuture.Error())
	assert.True(t, taskFuture.IsCanceled())
	assert.True(t, task.canceled)
}

func TestFuture_TaskOnSuccess(t *testing.T) {
	in := 1
	task := &eventsTask{in: in}
	taskFuture, err := simplefuture.New(task)
	assert.NoError(t, err)

	taskFuture.Run()
	taskFuture.Wait()

	assert.NoError(t, taskFuture.Error())
	assert.False(t, taskFuture.IsCanceled())
	assert.Nil(t, task.set)
}

func TestFuture_TaskOnError(t *testing.T) {
	in := 0
	task := &eventsTask{in: in}
	taskFuture, err := simplefuture.New(task)
	assert.NoError(t, err)

	taskFuture.Run()
	taskFuture.Wait()

	assert.Error(t, taskFuture.Error())
	assert.False(t, taskFuture.IsCanceled())
	assert.Error(t, task.set.(error))
}

type divideTask struct {
	a      int
	b      int
	result float64
}

func (t *divideTask) Run(func() bool) error {
	if t.b == 0 {
		return errors.New("division by zero")
	}

	t.result = float64(t.a / t.b)
	return nil
}

type longRunningTask struct {
	canceled bool
}

func (t *longRunningTask) OnCancel() {
	t.canceled = true
}

func (t *longRunningTask) Run(isCanceled func() bool) error {
	for {
		if isCanceled() {
			break
		}
	}

	return nil
}

type eventsTask struct {
	in  int
	set interface{}
}

func (e *eventsTask) OnSuccess(result interface{}) {
	e.set = result
}

func (e *eventsTask) OnError(err error) {
	e.set = err
}

func (e *eventsTask) Run(func() bool) error {
	if e.in == 0 {
		return errors.New("err")
	}
	return nil
}
