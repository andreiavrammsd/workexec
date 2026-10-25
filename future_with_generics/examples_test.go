package future_with_generics_test

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/andreiavrammsd/workexec/future_with_generics"
)

func ExampleFuture_Wait() {
	future1, err := future_with_generics.New[int64](&fibonacciTask{n: 1})
	if err != nil {
		log.Fatal(err)
	}

	future1.Run()
	future1.Wait()

	// Output:
	// on success: 1 -> 1
}

func ExampleFuture_Result() {
	future2, err := future_with_generics.New[int64](&fibonacciTask{n: 2})
	if err != nil {
		log.Fatal(err)
	}

	future2.Run()
	result, err := future2.Result()
	_, _ = fmt.Fprintln(os.Stdout, result)
	_, _ = fmt.Fprintln(os.Stdout, err == nil)

	// Output:
	// on success: 2 -> 1
	// 1
	// true
}

func ExampleFuture_Cancel() {
	fibonacci4 := &fibonacciTask{n: 4}
	future4, err := future_with_generics.New[int64](fibonacci4)
	if err != nil {
		log.Fatal(err)
	}

	future4.Run()
	future4.Cancel()
	result, err := future4.Result()
	_, _ = fmt.Fprintf(os.Stdout, "%v %v\n", result, err)
	_, _ = fmt.Fprintln(os.Stdout, fibonacci4.nums)
	_, _ = fmt.Fprintln(os.Stdout, future4.IsCanceled())

	// Output:
	// on cancel: 4
	// 0 <nil>
	// [0 1 0 0 0]
	// true
}

func ExampleFuture_error() {
	future0, err := future_with_generics.New[int64](&fibonacciTask{n: 0})
	if err != nil {
		log.Fatal(err)
	}

	future0.Run()
	future0.Wait()

	result, err := future0.Result()
	_, _ = fmt.Fprintln(os.Stdout, result == int64(0))
	_, _ = fmt.Fprintln(os.Stdout, err == nil)

	// Output:
	// on error: 0 -> n is zero
	// true
	// false
}

type fibonacciTask struct {
	n    int64
	nums []int64
}

func (f *fibonacciTask) OnSuccess(result int64) {
	_, _ = fmt.Fprintf(os.Stdout, "on success: %d -> %d\n", f.n, result)
}

func (f *fibonacciTask) OnError(err error) {
	_, _ = fmt.Fprintf(os.Stdout, "on error: %d -> %s\n", f.n, err)
}

func (f *fibonacciTask) OnCancel() {
	_, _ = fmt.Fprintf(os.Stdout, "on cancel: %d\n", f.n)
}

func (f *fibonacciTask) Run(isCanceled func() bool) (int64, error) {
	n := f.n

	if n == 0 {
		return 0, errors.New("n is zero")
	}

	nums := make([]int64, n+1, n+2)
	if n < 2 {
		nums = nums[0:2]
	}
	nums[0] = 0
	nums[1] = 1
	for i := int64(2); i <= n; i++ {
		if isCanceled() {
			break
		}

		nums[i] = nums[i-1] + nums[i-2]
	}

	f.nums = nums

	return nums[n], nil
}
