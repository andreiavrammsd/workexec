package simplefuture

// Task is the interface a task must implement to be part of a job.
type Task interface {
	// Run is the function which must do the actual work.
	// If the function passed as argument returns true, a cancellation has been
	// requested and the work should stopped.
	Run(func() bool) error
}

// SuccessfulTask is the interface a task must implement to be notified
// after is executed successfully.
type SuccessfulTask interface {
	OnSuccess()
}

// FailedTask is the interface a task must implement to be notified when it fails.
// OnError receives the returned error.
type FailedTask interface {
	OnError(error)
}

// CanceledTask is the interface a task must implement to be notified when it's canceled.
type CanceledTask interface {
	OnCancel()
}
