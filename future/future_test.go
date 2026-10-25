package future_test

import (
	"errors"
	"testing"
	"time"

	"github.com/andreiavrammsd/workexec/future"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_WithNilTaskError(t *testing.T) {
	t.Parallel()

	taskFuture, err := future.New(nil)
	assert.Nil(t, taskFuture)
	require.Error(t, err)
}

func TestFuture_Result(t *testing.T) {
	t.Parallel()

	taskFuture, err := future.New(&divideTask{a: 4, b: 2})
	require.NoError(t, err)

	taskFuture.Run()
	result, err := taskFuture.Result()

	assert.Equal(t, 2, result)
	require.NoError(t, err)
	assert.False(t, taskFuture.IsCanceled())
}

func TestFuture_Wait(t *testing.T) {
	t.Parallel()

	taskFuture, err := future.New(&divideTask{a: 4, b: 2})
	require.NoError(t, err)

	taskFuture.Run()
	taskFuture.Wait()
	taskFuture.Wait()

	result, err := taskFuture.Result()
	assert.Equal(t, 2, result)
	require.NoError(t, err)

	resultAgain, errAgain := taskFuture.Result()
	assert.Equal(t, 2, resultAgain)
	require.NoError(t, errAgain)

	assert.False(t, taskFuture.IsCanceled())
}

func TestFuture_Error(t *testing.T) {
	t.Parallel()

	taskFuture, err := future.New(&divideTask{a: 4, b: 0})
	require.NoError(t, err)

	taskFuture.Run()
	taskFuture.Wait()

	result, err := taskFuture.Result()

	assert.Nil(t, result)
	require.Error(t, err)
	assert.False(t, taskFuture.IsCanceled())
}

func TestFuture_Cancel(t *testing.T) {
	t.Parallel()

	task := &longRunningTask{}
	taskFuture, err := future.New(task)
	require.NoError(t, err)

	taskFuture.Run()
	taskFuture.Cancel()
	taskFuture.Wait()

	result, err := taskFuture.Result()
	assert.Nil(t, result)
	require.Error(t, err)
	assert.True(t, taskFuture.IsCanceled())

	assert.True(t, task.canceled)
}

func TestFuture_TaskOnSuccess(t *testing.T) {
	t.Parallel()

	in := 1
	task := &eventsTask{in: in}
	taskFuture, err := future.New(task)
	require.NoError(t, err)

	taskFuture.Run()
	taskFuture.Wait()
	result, err := taskFuture.Result()

	assert.Equal(t, in, result)
	require.NoError(t, err)
	assert.False(t, taskFuture.IsCanceled())

	assert.Equal(t, in, task.set)
}

func TestFuture_TaskOnError(t *testing.T) {
	t.Parallel()

	in := 0
	task := &eventsTask{in: in}
	taskFuture, err := future.New(task)
	require.NoError(t, err)

	taskFuture.Run()
	taskFuture.Wait()
	result, err := taskFuture.Result()

	assert.Nil(t, result)
	require.Error(t, err)
	assert.False(t, taskFuture.IsCanceled())

	errValue, ok := task.set.(error)
	require.True(t, ok)
	require.Error(t, errValue)
}

type divideTask struct {
	a int
	b int
}

func (t *divideTask) Run(func() bool) (any, error) {
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

func (t *longRunningTask) Run(isCanceled func() bool) (any, error) {
	for !isCanceled() {
		time.Sleep(time.Millisecond)
	}

	return nil, future.ErrCanceled
}

type eventsTask struct {
	in  int
	set any
}

func (e *eventsTask) OnSuccess(result any) {
	e.set = result
}

func (e *eventsTask) OnError(err error) {
	e.set = err
}

func (e *eventsTask) Run(func() bool) (any, error) {
	if e.in == 0 {
		return nil, errors.New("err")
	}

	return e.in, nil
}
