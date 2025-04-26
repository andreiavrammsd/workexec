package job_with_generics

// Task is the interface a task must implement to be part of a job.
type Task[T any] interface {
	Run(*Job[T]) (T, error)
}

// SuccessfulTask is the interface a task must implement to be notified
// after is executed successfully. OnSuccess receives the task result.
type SuccessfulTask[T any] interface {
	OnSuccess(T)
}

// FailedTask is the interface a task must implement to be notified when fails.
// OnError receives the returned error.
type FailedTask interface {
	OnError(error)
}

// CancelableTask is the interface a task must implement to be notified when it's canceled.
// OnCancel receives the error the task was canceled with.
type CancelableTask interface {
	OnCancel(error)
}
