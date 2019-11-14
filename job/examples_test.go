package job_test

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/andreiavrammsd/jobrunner/job"
)

func ExampleJob() {
	fibonacci, err := job.New(&fibonacciTask{n: 3})
	if err != nil {
		log.Fatal(err)
	}

	fibonacci.OnSuccess(func(result interface{}) {
		fmt.Println("on success result:", result.(uint))
	})

	future := fibonacci.Run()
	future.Wait()

	// Output:
	// on success result: 2
}

func ExampleJobGetResult() {
	fibonacci, err := job.New(&fibonacciTask{n: 3})
	if err != nil {
		log.Fatal(err)
	}

	future := fibonacci.Run()
	fmt.Println(future.Result().(uint))

	// Output:
	// 2
}

func ExampleJobError() {
	fibonacci, err := job.New(&fibonacciTask{n: 0})
	if err != nil {
		log.Fatal(err)
	}

	fibonacci.OnError(func(err error) {
		fmt.Println("on error:", err)
	})

	future := fibonacci.Run()
	future.Wait()

	if future.Result() != nil {
		log.Fatal("expected no result")
	}

	fmt.Println("wait error:", future.Error())

	// Output:
	// on error: n == 0
	// wait error: n == 0
}

func ExampleJobCancel() {
	fibonacci, err := job.New(&fibonacciTask{n: 3})
	if err != nil {
		log.Fatal(err)
	}

	fibonacci.OnCancel(func(err error) {
		fmt.Println("on cancel:", err)
	})

	fibonacci.Cancel(nil)

	future := fibonacci.Run()
	future.Wait()

	if !fibonacci.IsCanceled() {
		log.Fatal("expected to be canceled")
	}

	if !future.IsCanceled() {
		log.Fatal("expected to be canceled")
	}

	// Output:
	// on cancel: canceled
}

func ExampleJobCancelWithCustomError() {
	fibonacci, err := job.New(&fibonacciTask{n: 3})
	if err != nil {
		log.Fatal(err)
	}

	fibonacci.OnCancel(func(err error) {
		fmt.Println(err)
	})

	fibonacci.Cancel(errors.New("canceled by user"))

	future := fibonacci.Run()
	future.Wait()

	if !fibonacci.IsCanceled() {
		log.Fatal("expected to be canceled")
	}

	// Output:
	// canceled by user
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
