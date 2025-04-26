package job_with_generics_test

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/andreiavrammsd/workexec/job_with_generics"
)

func ExampleJob() {
	job1, err := job_with_generics.New[int](&defaultTask{num: 1})
	if err != nil {
		log.Fatal(err)
	}
	future1 := job1.Run()

	job2, err := job_with_generics.New[int](&defaultTask{num: -1})
	if err != nil {
		log.Fatal(err)
	}
	future2 := job2.Run()

	job3, err := job_with_generics.New[int](&defaultTask{num: 0})
	if err != nil {
		log.Fatal(err)
	}
	job3.Cancel(errors.New("job 3 was canceled"))
	future3 := job3.Run()

	fmt.Println(future1.Result())

	future2.Wait()

	future3.Wait()

	// Unordered output:
	// on success result: 2
	// 2
	// on error: negative number
	// on cancel: job 3 was canceled
}

func Example_cancel_long_running_tasks() {
	fibonacci10, err := job_with_generics.New[int64](&fibonacciTask{n: 10})
	if err != nil {
		log.Fatal(err)
	}

	fibonacci100, err := job_with_generics.New[int64](&fibonacciTask{n: 100})
	if err != nil {
		log.Fatal(err)
	}

	time.AfterFunc(time.Millisecond*10, func() {
		fibonacci10.Cancel(errors.New("10"))
		fibonacci100.Cancel(errors.New("100"))
	})

	future10 := fibonacci10.Run()
	future100 := fibonacci100.Run()

	future10.Wait()
	future100.Wait()

	// Unordered output:
	// on cancel: 10
	// on cancel: 100
}

type defaultTask struct {
	num int
}

func (t *defaultTask) Run(*job_with_generics.Job[int]) (int, error) {
	if t.num < 0 {
		return 0, errors.New("negative number")
	}

	return t.num + 1, nil
}

func (t *defaultTask) OnSuccess(result int) {
	fmt.Println("on success result:", result)
}

func (t *defaultTask) OnError(err error) {
	fmt.Println("on error:", err)
}

func (t *defaultTask) OnCancel(err error) {
	fmt.Println("on cancel:", err)
}

type fibonacciTask struct {
	n     int64
	count int64
}

func (f *fibonacciTask) OnCancel(err error) {
	fmt.Println("on cancel:", err)
}

func (f *fibonacciTask) Run(j *job_with_generics.Job[int64]) (int64, error) {
	n := f.n

	if n == 0 {
		return 0, errors.New("n == 0")
	}

	nums := make([]int64, n+1, n+2)
	if n < 2 {
		nums = nums[0:2]
	}
	nums[0] = 0
	nums[1] = 1
	for i := int64(2); i <= n; i++ {
		if j.IsCanceled() {
			return 0, nil
		}

		nums[i] = nums[i-1] + nums[i-2]
		if i > 5 {
			time.Sleep(time.Millisecond * 50)
		}
	}

	f.count = int64(len(nums))

	return nums[n], nil
}
