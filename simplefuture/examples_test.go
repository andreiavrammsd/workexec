package simplefuture_test

import (
	"errors"
	"fmt"
	"log"

	"github.com/andreiavrammsd/jobrunner/simplefuture"
)

func ExampleFuture_Wait() {
	task := &fibonacciTask{n: 1}
	future1, err := simplefuture.New(task)
	if err != nil {
		log.Fatal(err)
	}

	future1.Run()
	future1.Wait()

	fmt.Println(task.msg)

	// Output:
	// on success: 1 -> 1
}

func ExampleFuture_Cancel() {
	task := &fibonacciTask{n: 4}
	futureTask4, err := simplefuture.New(task)
	if err != nil {
		log.Fatal(err)
	}

	futureTask4.Cancel()
	futureTask4.Run()

	fmt.Println(task.result)
	fmt.Println(task.nums)
	fmt.Println(futureTask4.IsCanceled())

	// Output:
	// 0
	// []
	// true
}

func ExampleFuture_error() {
	task := &fibonacciTask{n: 0}
	future0, err := simplefuture.New(task)
	if err != nil {
		log.Fatal(err)
	}

	future0.Run()
	future0.Wait()

	fmt.Println(task.msg)
	fmt.Println(task.result)
	fmt.Println(future0.Error() == nil)

	// Output:
	// on error: 0 -> n is zero
	// 0
	// false
}

type fibonacciTask struct {
	n      uint
	nums   []uint
	result uint
	msg    string
}

func (f *fibonacciTask) OnSuccess() {
	f.msg = fmt.Sprintf("on success: %d -> %d", f.n, f.result)
}

func (f *fibonacciTask) OnError(err error) {
	f.msg = fmt.Sprintf("on error: %d -> %s", f.n, err)
}

func (f *fibonacciTask) OnCancel() {
	f.msg = fmt.Sprintf("on cancel: %d", f.n)
}

func (f *fibonacciTask) Run(isCanceled func() bool) error {
	n := f.n

	if n == 0 {
		return errors.New("n is zero")
	}

	nums := make([]uint, n+1, n+2)
	if n < 2 {
		nums = nums[0:2]
	}
	nums[0] = 0
	nums[1] = 1
	for i := uint(2); i <= n; i++ {
		if isCanceled() {
			break
		}

		nums[i] = nums[i-1] + nums[i-2]
	}

	f.nums = nums
	f.result = nums[n]
	return nil
}
