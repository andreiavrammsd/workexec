package job

type Wait struct {
	done     chan struct{}
	result   interface{}
	err      error
	canceled bool
}

func (w *Wait) Wait() {
	<-w.done
}

func (w *Wait) Result() interface{} {
	return w.result
}

func (w *Wait) Error() error {
	return w.err
}

func (w *Wait) IsCanceled() bool {
	return w.canceled
}
