package job

// Task is the interface a task must implement to be part of a job.
type Task interface {
	Run(*Job) (interface{}, error)
}

// SuccessfulTask is the interface a task must implement to be notified
// after is executed successfully.
type SuccessfulTask interface {
	OnSuccess(interface{})
}

// FailedTask is the interface a task must implement to be notified when fails.
type FailedTask interface {
	OnError(error)
}

// CancelableTask is the interface a task must implement to be notified when it's canceled.
type CancelableTask interface {
	OnCancel(error)
}
