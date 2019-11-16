package futurereflect

import (
	"reflect"
)

// Future represents a function which executes async work.
type Future struct {
	functionType  reflect.Type
	functionValue reflect.Value
	args          []interface{}
	wait          chan struct{}
	result        interface{}
	err           error
}

// Wait blocks until function is done.
func (f *Future) Wait() {
	<-f.wait
}

// Result retrieves result and error. It blocks if function is not done.
func (f *Future) Result() (interface{}, error) {
	<-f.wait
	return f.result, f.err
}

func (f *Future) run() {
	defer close(f.wait)

	numParams := f.functionType.NumIn()
	values := make([]reflect.Value, numParams)
	for i := 0; i < numParams; i++ {
		values[i] = reflect.ValueOf(f.args[i])
	}

	ret := f.functionValue.Call(values)

	if len(ret) == 0 {
		return
	}

	f.result = ret[0].Interface()

	if f.functionType.NumOut() > 1 && !ret[1].IsNil() {
		// nolint
		f.err = ret[1].Interface().(error)
	}
}

// Callable is the function to call in order to start running the passed function.
type Callable func(args ...interface{}) *Future

// New returns a Callable.
func New(function interface{}) Callable {
	if function == nil {
		return nil
	}

	functionType := reflect.TypeOf(function)
	if functionType.Kind() != reflect.Func {
		return nil
	}

	errorInterface := reflect.TypeOf((*error)(nil)).Elem()
	if functionType.NumOut() > 1 && !functionType.Out(1).Implements(errorInterface) {
		return nil
	}

	return func(args ...interface{}) *Future {
		future := &Future{
			functionType:  functionType,
			functionValue: reflect.ValueOf(function),
			args:          args,
			wait:          make(chan struct{}),
		}

		go future.run()

		return future
	}
}
