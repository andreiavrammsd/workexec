package future_test

import (
	"errors"
	"testing"

	"github.com/andreiavrammsd/jobrunner/future"
	"github.com/stretchr/testify/assert"
)

func TestNew_WithNilTaskError(t *testing.T) {
	taskFuture, err := future.New(nil)
	assert.Nil(t, taskFuture)
	assert.Error(t, err)
}

func TestFuture_Result(t *testing.T) {
	taskFuture, err := future.New(&divideTask{a: 4, b: 2})
	assert.NoError(t, err)

	result, err := taskFuture.Result()

	assert.Equal(t, 2, result)
	assert.NoError(t, err)
	assert.False(t, taskFuture.IsCanceled())
}

func TestFuture_Wait(t *testing.T) {
	taskFuture, err := future.New(&divideTask{a: 4, b: 2})
	assert.NoError(t, err)

	taskFuture.Wait()
	taskFuture.Wait()

	result, err := taskFuture.Result()
	assert.Equal(t, 2, result)
	assert.NoError(t, err)

	resultAgain, errAgain := taskFuture.Result()
	assert.Equal(t, 2, resultAgain)
	assert.NoError(t, errAgain)

	assert.False(t, taskFuture.IsCanceled())
}

func TestFuture_Error(t *testing.T) {
	taskFuture, err := future.New(&divideTask{a: 4, b: 0})
	assert.NoError(t, err)

	taskFuture.Wait()

	result, err := taskFuture.Result()

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.False(t, taskFuture.IsCanceled())
}

func TestFuture_Cancel(t *testing.T) {
	task := &longRunningTask{}
	taskFuture, err := future.New(task)
	assert.NoError(t, err)

	taskFuture.Cancel()
	taskFuture.Wait()

	result, err := taskFuture.Result()
	assert.Nil(t, result)
	assert.NoError(t, err)
	assert.True(t, taskFuture.IsCanceled())

	assert.True(t, task.canceled)
}

func TestFuture_TaskOnSuccess(t *testing.T) {
	in := 1
	task := &eventsTask{in: in}
	taskFuture, err := future.New(task)
	assert.NoError(t, err)

	taskFuture.Wait()

	result, err := taskFuture.Result()

	assert.Equal(t, in, result)
	assert.NoError(t, err)
	assert.False(t, taskFuture.IsCanceled())

	assert.Equal(t, in, task.set)
}

func TestFuture_TaskOnError(t *testing.T) {
	in := 0
	task := &eventsTask{in: in}
	taskFuture, err := future.New(task)
	assert.NoError(t, err)

	taskFuture.Wait()

	result, err := taskFuture.Result()

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.False(t, taskFuture.IsCanceled())

	assert.Error(t, task.set.(error))
}

type divideTask struct {
	a int
	b int
}

func (t *divideTask) Run(func() bool) (interface{}, error) {
	if t.b == 0 {
		return nil, errors.New("division by zero")
	}

	return t.a / t.b, nil
}

type longRunningTask struct {
	canceled bool
}

func (t *longRunningTask) OnCancel() {
	t.canceled = true
}

func (t *longRunningTask) Run(isCanceled func() bool) (interface{}, error) {
	for {
		if isCanceled() {
			break
		}
	}

	return nil, nil
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

func (e *eventsTask) Run(func() bool) (interface{}, error) {
	if e.in == 0 {
		return nil, errors.New("err")
	}
	return e.in, nil
}
