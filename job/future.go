package job

// Future represents a job running async.
type Future struct {
	done     chan struct{}
	result   interface{}
	err      error
	canceled bool
}

// Wait blocks until job is done.
func (w *Future) Wait() {
	<-w.done
}

// Result returns job result. Blocks until job is done.
func (w *Future) Result() interface{} {
	<-w.done
	return w.result
}

// Error returns job error.
func (w *Future) Error() error {
	return w.err
}

// IsCanceled returns true if job was canceled.
func (w *Future) IsCanceled() bool {
	return w.canceled
}
