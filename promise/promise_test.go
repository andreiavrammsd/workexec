package promise_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/andreiavrammsd/workexec/promise"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.NotNil(t, promise.New(nil))
}

func TestAwait(t *testing.T) {
	exec := &division{a: 4, b: 2}
	err := promise.New(exec).Await()
	assert.NoError(t, err)
	assert.Equal(t, 2, exec.result)

	assert.Nil(t, promise.New(nil).Await())

	exec = &division{a: 2, b: 2}
	err = promise.New(exec).Await()
	assert.NoError(t, err)
	assert.Equal(t, 1, exec.result)
}

func TestPromise_Then(t *testing.T) {
	exec := &division{a: 4, b: 2}
	then := func() {
		exec.result++
	}
	err := promise.New(exec).Then(then).Await()
	assert.NoError(t, err)
	assert.Equal(t, 3, exec.result)
}

func TestPromise_Error(t *testing.T) {
	exec := &division{a: 4, b: 0}
	e := func(err error) {
		exec.result--
		exec.err = err
	}
	err := promise.New(exec).Error(e).Await()
	assert.Error(t, err)
	assert.Equal(t, -1, exec.result)
	assert.Error(t, err, exec.err)
}

func TestAll(t *testing.T) {
	executors := []promise.Executor{
		&division{a: 2, b: 2},
		&division{a: 0, b: 2},
		&division{a: 0, b: 0},
		&division{a: 4, b: 2},
	}

	var errVal error
	err := promise.New(executors...).Error(func(err error) {
		errVal = err
	}).Await()

	assert.Error(t, err)
	assert.Equal(t, err, errVal)
}

func TestAll_WithAllErrors(t *testing.T) {
	executors := []promise.Executor{
		&division{a: 2, b: 0},
		&division{a: 0, b: 0},
		&division{a: 0, b: 0},
		&division{a: 4, b: 0},
	}

	lock := sync.Mutex{}
	var errVal error
	err := promise.New(executors...).Error(func(err error) {
		lock.Lock()
		errVal = err
		lock.Unlock()
	}).Await()

	assert.Error(t, err)
	lock.Lock()
	assert.Equal(t, err, errVal)
	lock.Unlock()
}

type division struct {
	a, b, result int
	err          error
}

func (d *division) Execute() error {
	if d.b == 0 {
		return errors.New("division by zero")
	}

	d.result = d.a / d.b
	return nil
}
