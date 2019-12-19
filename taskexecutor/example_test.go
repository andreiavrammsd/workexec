package taskexecutor_test

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/andreiavrammsd/workexec/future"

	"github.com/andreiavrammsd/workexec/taskexecutor"
)

func Example() {
	// Set concurrency
	config := taskexecutor.Config{
		Concurrency: 1,
		QueueSize:   1,
	}
	// Create new executor
	taskExecutor := taskexecutor.New(config)

	// Start the executor
	taskExecutor.Start()

	// Create tasks
	taskFuture, err := future.New(&task{n: 0})
	if err != nil {
		log.Fatal(err)
	}

	// Submit tasks
	if err := taskExecutor.Submit(taskFuture); err != nil {
		fmt.Println(err)
	}

	time.AfterFunc(time.Millisecond*50, func() {
		taskExecutor.Stop()
	})

	taskExecutor.Wait()

	// Output:
	// on error: 0 -> n is zero
}

type task struct {
	n uint
}

func (f *task) OnSuccess(result interface{}) {
	fmt.Printf("on success: %d -> %d\n", f.n, result)
}

func (f *task) OnError(err error) {
	fmt.Printf("on error: %d -> %s\n", f.n, err)
}

func (f *task) OnCancel() {
	fmt.Printf("on cancel: %d\n", f.n)
}

func (f *task) Run(isCanceled func() bool) (interface{}, error) {
	n := f.n

	if n == 0 {
		return nil, errors.New("n is zero")
	}

	result := uint(0)
	for i := uint(2); i <= n; i++ {
		result += i + n
		if isCanceled() {
			break
		}
	}

	return result, nil
}
