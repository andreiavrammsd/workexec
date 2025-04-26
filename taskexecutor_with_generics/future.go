package taskexecutor_with_generics

import "github.com/andreiavrammsd/workexec/future_with_generics"

// Future represents a task which executes async work.
type Future[T any] interface {
	Run()
	Wait()
	Cancel()
	Result() (T, error)
	IsCanceled() bool
}

type AnyFuture[T any] struct {
	Future Future[T]
}

func (a *AnyFuture[T]) Run() {
	a.Future.Run()
}

func (a *AnyFuture[T]) Wait() {
	a.Future.Wait()
}

func (a *AnyFuture[T]) Cancel() {
	a.Future.Cancel()
}

func (a *AnyFuture[T]) IsCanceled() bool {
	return a.Future.IsCanceled()
}

func (a *AnyFuture[T]) Result() (any, error) {
	return a.Future.Result()
}

func NewFuture[T any](task future_with_generics.Task[T]) (Future[any], error) {
	future, err := future_with_generics.New(task)
	if err != nil {
		return nil, err
	}

	return &AnyFuture[T]{Future: future}, nil
}

// func WrapFuture[T any](future Future[T]) Future[any] {
// 	future_with_generics.New[uint](&task{n: 0})
// 	return &AnyFuture[T]{Future: future}
// }
