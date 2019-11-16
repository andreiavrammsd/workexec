package future_test

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/andreiavrammsd/jobrunner/future"
)

func Example() {
	future0, err := future.New(&fibonacciTask{n: 0})
	if err != nil {
		log.Fatal(err)
	}

	future1, err := future.New(&fibonacciTask{n: 1})
	if err != nil {
		log.Fatal(err)
	}

	future3, err := future.New(&fibonacciTask{n: 3})
	if err != nil {
		log.Fatal(err)
	}

	fibonacci4 := &fibonacciTask{n: 4}
	future4, err := future.New(fibonacci4)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("future0")
	future0.Wait()
	fmt.Println()

	fmt.Println("future1")
	future1.Wait()
	fmt.Println()

	fmt.Println("future3")
	fmt.Println(future3.Result())
	fmt.Println()

	fmt.Println("future4")
	go future4.Cancel()
	future4.Wait()
	fmt.Println(future4.Result())
	fmt.Println(fibonacci4.nums)
	fmt.Println(future4.IsCanceled())

	// Output:
	// future0
	// on error: n == 0
	//
	// future1
	// on success: 1
	//
	// future3
	// on success: 2
	// 2 <nil>
	//
	// future4
	// on cancel
	// 0 <nil>
	// [0 1 0 0 0]
	// true
}

type fibonacciTask struct {
	n    uint
	nums []uint
}

func (f *fibonacciTask) OnSuccess(result interface{}) {
	fmt.Println("on success:", result.(uint))
}

func (f *fibonacciTask) OnError(err error) {
	fmt.Println("on error:", err)
}

func (f *fibonacciTask) OnCancel() {
	fmt.Println("on cancel")
}

func (f *fibonacciTask) Run(isCanceled func() bool) (interface{}, error) {
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
		time.Sleep(time.Millisecond * 20)
		if isCanceled() {
			break
		}

		nums[i] = nums[i-1] + nums[i-2]
	}

	f.nums = nums

	return nums[n], nil
}
