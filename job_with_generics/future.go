package job_with_generics

// Future represents a job running async.
type Future[T any] struct {
	done     chan struct{}
	result   T
	err      error
	canceled bool
}

// Wait blocks until job is done.
func (f *Future[T]) Wait() {
	<-f.done
}

// Result returns job result. Blocks until job is done.
func (f *Future[T]) Result() T {
	<-f.done
	return f.result
}

// Error returns job error.
func (f *Future[T]) Error() error {
	return f.err
}

// IsCanceled returns true if job was canceled.
func (f *Future[T]) IsCanceled() bool {
	return f.canceled
}
