package job

// Future represents a job running async.
type Future struct {
	done     chan struct{}
	result   any
	err      error
	canceled bool
}

// Wait blocks until job is done.
func (f *Future) Wait() {
	<-f.done
}

// Result returns job result. Blocks until job is done.
func (f *Future) Result() any {
	<-f.done

	return f.result
}

// Error returns job error.
func (f *Future) Error() error {
	return f.err
}

// IsCanceled returns true if job was canceled.
func (f *Future) IsCanceled() bool {
	return f.canceled
}
