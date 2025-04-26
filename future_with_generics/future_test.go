package future_with_generics_test

import (
	"errors"
	"testing"

	"github.com/andreiavrammsd/workexec/future_with_generics"
	"github.com/stretchr/testify/assert"
)

func TestNew_WithNilTaskError(t *testing.T) {
	taskFuture, err := future_with_generics.New[int](nil)
	assert.Nil(t, taskFuture)
	assert.Error(t, err)
}

func TestFuture_Result(t *testing.T) {
	taskFuture, err := future_with_generics.New[float64](&divideTask{a: 4, b: 2})
	assert.NoError(t, err)

	taskFuture.Run()
	result, err := taskFuture.Result()

	assert.Equal(t, 2.0, result)
	assert.NoError(t, err)
	assert.False(t, taskFuture.IsCanceled())
}

func TestFuture_Wait(t *testing.T) {
	taskFuture, err := future_with_generics.New[float64](&divideTask{a: 4, b: 2})
	assert.NoError(t, err)

	taskFuture.Run()
	taskFuture.Wait()
	taskFuture.Wait()

	result, err := taskFuture.Result()
	assert.Equal(t, 2.0, result)
	assert.NoError(t, err)

	resultAgain, errAgain := taskFuture.Result()
	assert.Equal(t, 2.0, resultAgain)
	assert.NoError(t, errAgain)

	assert.False(t, taskFuture.IsCanceled())
}

func TestFuture_Error(t *testing.T) {
	taskFuture, err := future_with_generics.New[float64](&divideTask{a: 4, b: 0})
	assert.NoError(t, err)

	taskFuture.Run()
	taskFuture.Wait()

	result, err := taskFuture.Result()

	assert.Equal(t, result, 0.0)
	assert.Error(t, err)
	assert.False(t, taskFuture.IsCanceled())
}

func TestFuture_Cancel(t *testing.T) {
	task := &longRunningTask{}
	taskFuture, err := future_with_generics.New[int](task)
	assert.NoError(t, err)

	taskFuture.Run()
	taskFuture.Cancel()
	taskFuture.Wait()

	result, err := taskFuture.Result()
	assert.Equal(t, result, 0)
	assert.NoError(t, err)
	assert.True(t, taskFuture.IsCanceled())

	assert.True(t, task.canceled)
}

func TestFuture_TaskOnSuccess(t *testing.T) {
	in := 1
	task := &eventsTask{in: in}
	taskFuture, err := future_with_generics.New[int](task)
	assert.NoError(t, err)

	taskFuture.Run()
	taskFuture.Wait()
	result, err := taskFuture.Result()

	assert.Equal(t, in, result)
	assert.NoError(t, err)
	assert.False(t, taskFuture.IsCanceled())

	assert.Equal(t, in, task.set)
	assert.NoError(t, task.err)
}

func TestFuture_TaskOnError(t *testing.T) {
	in := 0
	task := &eventsTask{in: in}
	taskFuture, err := future_with_generics.New[int](task)
	assert.NoError(t, err)

	taskFuture.Run()
	taskFuture.Wait()
	result, err := taskFuture.Result()

	assert.Equal(t, result, 0)
	assert.Error(t, err)
	assert.False(t, taskFuture.IsCanceled())

	assert.Error(t, task.err)
}

type divideTask struct {
	a int
	b int
}

func (t *divideTask) Run(func() bool) (float64, error) {
	if t.b == 0 {
		return 0, errors.New("division by zero")
	}

	return float64(t.a / t.b), nil
}

type longRunningTask struct {
	canceled bool
}

func (t *longRunningTask) OnCancel() {
	t.canceled = true
}

func (t *longRunningTask) Run(isCanceled func() bool) (int, error) {
	for !isCanceled() {
	}

	return 0, nil
}

type eventsTask struct {
	in  int
	set int
	err error
}

func (e *eventsTask) OnSuccess(result int) {
	e.set = result
}

func (e *eventsTask) OnError(err error) {
	e.err = err
}

func (e *eventsTask) Run(func() bool) (int, error) {
	if e.in == 0 {
		return 0, errors.New("err")
	}
	return e.in, nil
}
