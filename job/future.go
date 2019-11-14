package job

type Future struct {
	done     chan struct{}
	result   interface{}
	err      error
	canceled bool
}

func (w *Future) Wait() {
	<-w.done
}

func (w *Future) Result() interface{} {
	<-w.done
	return w.result
}

func (w *Future) Error() error {
	return w.err
}

func (w *Future) IsCanceled() bool {
	return w.canceled
}
