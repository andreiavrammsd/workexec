package taskexecutor_with_generics_test

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/andreiavrammsd/workexec/taskexecutor_with_generics"
)

func Example() {
	// Set concurrency
	config := taskexecutor_with_generics.Config{
		Concurrency: 1,
		QueueSize:   1,
	}
	// Create new executor
	taskExecutor := taskexecutor_with_generics.New(config)

	// Start the executor
	taskExecutor.Start()

	// Submit tasks
	taskFuture, err := taskexecutor_with_generics.NewFuture[uint](&task{n: 0})
	if err != nil {
		log.Fatal(err)
	}
	if err := taskExecutor.Submit(taskFuture); err != nil {
		fmt.Println(err)
	}

	anotherTaskFuture, err := taskexecutor_with_generics.NewFuture[string](&anotherTask{input: "abc"})
	if err != nil {
		log.Fatal(err)
	}
	if err := taskExecutor.Submit(anotherTaskFuture); err != nil {
		fmt.Println(err)
	}

	canceledTaskFuture, err := taskexecutor_with_generics.NewFuture[struct{}](&canceledTask{})
	if err != nil {
		log.Fatal(err)
	}
	if err := taskExecutor.Submit(canceledTaskFuture); err != nil {
		fmt.Println(err)
	}

	canceledTaskFuture.Cancel()

	// Stop executor
	time.AfterFunc(time.Millisecond*50, func() {
		taskExecutor.Stop()
	})

	taskExecutor.Wait()

	// Output:
	// on error: 0 -> n is zero
	// on success: abc -> abcabc
	// on cancel
}

type task struct {
	n uint
}

func (f *task) OnSuccess(result uint) {
	fmt.Printf("on success: %d -> %d\n", f.n, result)
}

func (f *task) OnError(err error) {
	fmt.Printf("on error: %d -> %s\n", f.n, err)
}

func (f *task) OnCancel() {
	fmt.Printf("on cancel: %d\n", f.n)
}

func (f *task) Run(isCanceled func() bool) (uint, error) {
	n := f.n

	if n == 0 {
		return 0, errors.New("n is zero")
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

type anotherTask struct {
	input string
}

func (f *anotherTask) OnSuccess(result string) {
	fmt.Printf("on success: %s -> %s\n", f.input, result)
}

func (f *anotherTask) OnError(err error) {
	fmt.Printf("on error: %s -> %s\n", f.input, err)
}

func (f *anotherTask) OnCancel() {
	fmt.Printf("on cancel: %s\n", f.input)
}

func (f *anotherTask) Run(isCanceled func() bool) (string, error) {
	if f.input == "" {
		return "", errors.New("input is empty")
	}

	return f.input + f.input, nil
}

type canceledTask struct {
}

func (f *canceledTask) OnSuccess(result string) {
	fmt.Printf("on success: %s\n", result)
}

func (f *canceledTask) OnError(err error) {
	fmt.Printf("on error: %s\n", err)
}

func (f *canceledTask) OnCancel() {
	fmt.Printf("on cancel\n")
}

func (f *canceledTask) Run(isCanceled func() bool) (struct{}, error) {
	time.Sleep(time.Millisecond * 100)
	return struct{}{}, nil
}
