package job_test

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/andreiavrammsd/jobrunner/job"
)

func ExampleSyncJob() {
	task := &fibonacciTask{n: 3}
	fibonacci, err := job.New(task)
	if err != nil {
		log.Fatal(err)
	}

	if fibonacci.ID() == job.ID("") {
		log.Fatal("expected job id")
	}

	fibonacci.OnSuccess(func(result interface{}) {
		fmt.Println("on success result:", result.(uint))
	})

	result, err := fibonacci.Run()
	if err != nil {
		log.Fatal("error", err)
	}

	fmt.Println("returned result:", result.(uint))
	fmt.Println("numbers calculated:", task.count)

	// Output:
	// on success result: 2
	// returned result: 2
	// numbers calculated: 4
}

func ExampleSyncJobError() {
	fibonacci, err := job.New(&fibonacciTask{n: 0})
	if err != nil {
		log.Fatal(err)
	}

	fibonacci.OnError(func(err error) {
		fmt.Println("on error:", err)
	})

	result, err := fibonacci.Run()

	if result != nil {
		log.Fatal("expected no result")
	}

	if err != nil {
		fmt.Println("returned error:", err)
	}

	// Output:
	// on error: n == 0
	// returned error: n == 0
}

func ExampleAsyncJob() {
	fibonacci, err := job.New(&fibonacciTask{n: 3})
	if err != nil {
		log.Fatal(err)
	}

	fibonacci.OnSuccess(func(result interface{}) {
		fmt.Println("on success result:", result.(uint))
	})

	wait := fibonacci.Go()
	wait.Wait()

	// Output:
	// on success result: 2
}

func ExampleAsyncJobError() {
	fibonacci, err := job.New(&fibonacciTask{n: 0})
	if err != nil {
		log.Fatal(err)
	}

	fibonacci.OnError(func(err error) {
		fmt.Println("on error:", err)
	})

	wait := fibonacci.Go()
	wait.Wait()

	if wait.Result() != nil {
		log.Fatal("expected no result")
	}

	fmt.Println("wait error:", wait.Error())

	// Output:
	// on error: n == 0
	// wait error: n == 0
}

func ExampleAsyncJobCancel() {
	fibonacci, err := job.New(&fibonacciTask{n: 3})
	if err != nil {
		log.Fatal(err)
	}

	fibonacci.OnCancel(func(err error) {
		fmt.Println(err)
	})

	fibonacci.Cancel(nil)

	wait := fibonacci.Go()
	wait.Wait()

	if !fibonacci.IsCanceled() {
		log.Fatal("expected to be canceled")
	}

	// Output:
	// canceled
}

func ExampleAsyncJobCancelWithCustomError() {
	fibonacci, err := job.New(&fibonacciTask{n: 3})
	if err != nil {
		log.Fatal(err)
	}

	fibonacci.OnCancel(func(err error) {
		fmt.Println(err)
	})

	fibonacci.Cancel(errors.New("canceled by user"))

	wait := fibonacci.Go()
	wait.Wait()

	if !fibonacci.IsCanceled() {
		log.Fatal("expected to be canceled")
	}

	// Output:
	// canceled by user
}

func ExampleJobWithOptions() {
	fibonacci, err := job.New(&fibonacciTask{n: 6}, &job.WithTimeout{
		Timeout: time.Millisecond,
	})
	if err != nil {
		log.Fatal(err)
	}

	fibonacci.OnCancel(func(err error) {
		fmt.Println(err)
	})

	_, err = fibonacci.Run()
	if err != nil {
		log.Fatal(err)
	}

	if !fibonacci.IsCanceled() {
		log.Fatal("expected to be canceled")
	}

	// Output:
	// timeout
}

func ExampleNewJobError() {
	j, err := job.New(nil)

	if j != nil {
		log.Fatal("expected no job")
	}

	fmt.Println(err)

	// Output:
	// task function not passed
}

type fibonacciTask struct {
	n     uint
	count uint
}

func (f *fibonacciTask) Run(job *job.Job) (interface{}, error) {
	n := f.n

	if n == 0 {
		return nil, errors.New("n == 0")
	}

	nums := make([]uint, n+1, n+2)
	if n < 2 {
		nums = nums[0:2]
	}
	nums[0] = 0
	nums[1] = 1
	for i := uint(2); i <= n; i++ {
		if job.IsCanceled() {
			return nil, nil
		}

		nums[i] = nums[i-1] + nums[i-2]
		if i > 5 {
			time.Sleep(time.Millisecond * 50)
		}
	}

	f.count = uint(len(nums))

	return nums[n], nil
}
