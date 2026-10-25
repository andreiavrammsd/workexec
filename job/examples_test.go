package job_test

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/andreiavrammsd/workexec/job"
)

func ExampleJob() {
	job1, err := job.New(&defaultTask{num: 1})
	if err != nil {
		log.Fatal(err)
	}
	future1 := job1.Run()

	job2, err := job.New(&defaultTask{num: -1})
	if err != nil {
		log.Fatal(err)
	}
	future2 := job2.Run()

	job3, err := job.New(&defaultTask{num: 0})
	if err != nil {
		log.Fatal(err)
	}
	job3.Cancel(errors.New("job 3 was canceled"))
	future3 := job3.Run()

	result, ok := future1.Result().(int)
	if !ok {
		log.Fatal("unexpected result type")
	}
	_, _ = fmt.Fprintln(os.Stdout, result)

	future2.Wait()

	future3.Wait()

	// Unordered output:
	// on success result: 2
	// 2
	// on error: negative number
	// on cancel: job 3 was canceled
}

func Example_cancel_long_running_tasks() {
	fibonacci10, err := job.New(&fibonacciTask{n: 10})
	if err != nil {
		log.Fatal(err)
	}

	fibonacci100, err := job.New(&fibonacciTask{n: 100})
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

func (t *defaultTask) Run(*job.Job) (any, error) {
	if t.num < 0 {
		return nil, errors.New("negative number")
	}

	return t.num + 1, nil
}

func (t *defaultTask) OnSuccess(result any) {
	value, ok := result.(int)
	if !ok {
		return
	}
	_, _ = fmt.Fprintf(os.Stdout, "on success result: %d\n", value)
}

func (t *defaultTask) OnError(err error) {
	_, _ = fmt.Fprintf(os.Stdout, "on error: %s\n", err)
}

func (t *defaultTask) OnCancel(err error) {
	_, _ = fmt.Fprintf(os.Stdout, "on cancel: %s\n", err)
}

type fibonacciTask struct {
	n     uint
	count uint
}

func (f *fibonacciTask) OnCancel(err error) {
	_, _ = fmt.Fprintf(os.Stdout, "on cancel: %s\n", err)
}

func (f *fibonacciTask) Run(j *job.Job) (any, error) {
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
		if j.IsCanceled() {
			return nil, job.ErrCanceled
		}

		nums[i] = nums[i-1] + nums[i-2]
		if i > 5 {
			time.Sleep(time.Millisecond * 50)
		}
	}

	f.count = uint(len(nums))

	return nums[n], nil
}
