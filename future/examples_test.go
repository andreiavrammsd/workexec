package future_test

import (
	"errors"
	"fmt"
	"log"

	"github.com/andreiavrammsd/workexec/future"
)

func ExampleFuture_Wait() {
	future1, err := future.New(&fibonacciTask{n: 1})
	if err != nil {
		log.Fatal(err)
	}

	future1.Run()
	future1.Wait()

	// Output:
	// on success: 1 -> 1
}

func ExampleFuture_Result() {
	future2, err := future.New(&fibonacciTask{n: 2})
	if err != nil {
		log.Fatal(err)
	}

	future2.Run()
	result, err := future2.Result()
	fmt.Println(result)
	fmt.Println(err == nil)

	// Output:
	// on success: 2 -> 1
	// 1
	// true
}

func ExampleFuture_Cancel() {
	fibonacci4 := &fibonacciTask{n: 4}
	future4, err := future.New(fibonacci4)
	if err != nil {
		log.Fatal(err)
	}

	future4.Run()
	future4.Cancel()
	fmt.Println(future4.Result())
	fmt.Println(fibonacci4.nums)
	fmt.Println(future4.IsCanceled())

	// Output:
	// on cancel: 4
	// 0 <nil>
	// [0 1 0 0 0]
	// true
}

func ExampleFuture_error() {
	future0, err := future.New(&fibonacciTask{n: 0})
	if err != nil {
		log.Fatal(err)
	}

	future0.Run()
	future0.Wait()

	result, err := future0.Result()
	fmt.Println(result == nil)
	fmt.Println(err == nil)

	// Output:
	// on error: 0 -> n is zero
	// true
	// false
}

type fibonacciTask struct {
	n    uint
	nums []uint
}

func (f *fibonacciTask) OnSuccess(result interface{}) {
	fmt.Printf("on success: %d -> %d\n", f.n, result)
}

func (f *fibonacciTask) OnError(err error) {
	fmt.Printf("on error: %d -> %s\n", f.n, err)
}

func (f *fibonacciTask) OnCancel() {
	fmt.Printf("on cancel: %d\n", f.n)
}

func (f *fibonacciTask) Run(isCanceled func() bool) (interface{}, error) {
	n := f.n

	if n == 0 {
		return nil, errors.New("n is zero")
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

	return nums[n], nil
}
