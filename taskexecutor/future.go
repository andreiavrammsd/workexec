package taskexecutor

// Future represents a task which executes async work.
type Future interface {
	Run()
	Wait()
	Cancel()
	Result() (any, error)
	IsCanceled() bool
}
