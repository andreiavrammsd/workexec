package simplefuture_test

import (
	"errors"
	"runtime"
	"testing"

	"github.com/andreiavrammsd/workexec/simplefuture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_WithNilTaskError(t *testing.T) {
	t.Parallel()

	taskFuture, err := simplefuture.New(nil)
	assert.Nil(t, taskFuture)
	require.Error(t, err)
}

func TestFuture_Wait(t *testing.T) {
	t.Parallel()

	task := &divideTask{a: 4, b: 2}
	taskFuture, err := simplefuture.New(task)
	require.NoError(t, err)

	taskFuture.Run()
	taskFuture.Wait()
	taskFuture.Wait()

	assert.InDelta(t, 2.0, task.result, 0.0001)
	require.NoError(t, taskFuture.Error())
	assert.False(t, taskFuture.IsCanceled())
}

func TestFuture_Error(t *testing.T) {
	t.Parallel()

	task := &divideTask{a: 4, b: 0}
	taskFuture, err := simplefuture.New(task)
	require.NoError(t, err)

	taskFuture.Run()
	taskFuture.Wait()

	assert.InDelta(t, 0.0, task.result, 0.0001)
	require.Error(t, taskFuture.Error())
	assert.False(t, taskFuture.IsCanceled())
}

func TestFuture_Cancel(t *testing.T) {
	t.Parallel()

	task := &longRunningTask{}
	taskFuture, err := simplefuture.New(task)
	require.NoError(t, err)

	taskFuture.Run()
	taskFuture.Cancel()
	taskFuture.Wait()

	require.NoError(t, taskFuture.Error())
	assert.True(t, taskFuture.IsCanceled())
	assert.True(t, task.canceled)
}

func TestFuture_TaskOnSuccess(t *testing.T) {
	t.Parallel()

	in := 1
	task := &eventsTask{in: in}
	taskFuture, err := simplefuture.New(task)
	require.NoError(t, err)

	taskFuture.Run()
	taskFuture.Wait()

	require.NoError(t, taskFuture.Error())
	assert.False(t, taskFuture.IsCanceled())
	assert.Nil(t, task.set)
}

func TestFuture_TaskOnError(t *testing.T) {
	t.Parallel()

	in := 0
	task := &eventsTask{in: in}
	taskFuture, err := simplefuture.New(task)
	require.NoError(t, err)

	taskFuture.Run()
	taskFuture.Wait()

	require.Error(t, taskFuture.Error())
	assert.False(t, taskFuture.IsCanceled())
	errValue, ok := task.set.(error)
	require.True(t, ok)
	require.Error(t, errValue)
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
	for !isCanceled() {
		runtime.Gosched()
	}

	return nil
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

func (e *eventsTask) Run(func() bool) error {
	if e.in == 0 {
		return errors.New("err")
	}

	return nil
}
