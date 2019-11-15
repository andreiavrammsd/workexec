package runner_test

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/andreiavrammsd/jobrunner/job"
	"github.com/andreiavrammsd/jobrunner/runner"
)

func ExampleRunner() {
	job1, err := job.New(&defaultTask{num: 1})
	if err != nil {
		log.Fatal(err)
	}

	job2, err := job.New(&defaultTask{num: -1})
	if err != nil {
		log.Fatal(err)
	}

	job3, err := job.New(&defaultTask{num: 0})
	if err != nil {
		log.Fatal(err)
	}
	job3.Cancel(errors.New("job 3 was canceled"))

	// Setup runner
	c := runner.Config{
		Concurrency: 2,
		QueueSize:   2,
	}
	r := runner.New(c)
	r.Start()

	// Add jobs to runner queue
	if err := r.Enqueue(job1, job2, job3); err != nil {
		log.Fatal(err)
	}

	// Ask runner to stop
	time.AfterFunc(time.Millisecond*10, func() {
		r.Stop()
	})

	// Wait for runner to run jobs
	r.Wait()

	// Unordered output:
	// on success result: 2
	// on cancel: job 3 was canceled
	// on error: negative number
}

func ExampleRunner_cancel_running_job_by_ID() {
	// Create job
	jobA, err := job.New(&cancelableTask{&myTask{text: "A"}, "canceled by user"})
	if err != nil {
		log.Fatal(err)
	}

	// Setup runner
	c := runner.Config{
		Concurrency: 1,
		QueueSize:   1,
	}
	r := runner.New(c)
	r.Start()

	// Add job to runner queue
	if err := r.Enqueue(jobA); err != nil {
		log.Fatal(err)
	}

	// Cancel a job by ID
	time.AfterFunc(time.Millisecond*1, func() {
		r.Cancel(jobA.ID())
	})

	// Ask runner to stop
	time.AfterFunc(time.Millisecond*100, func() {
		r.Stop()
	})

	// Wait for runner to run jobs
	r.Wait()

	// Output:
	// A
	// on cancel: canceled by user
}

func ExampleRunner_cancel_job() {
	// Create job
	jobA, err := job.New(&cancelableTask{&myTask{text: "A"}, "canceled by user"})
	if err != nil {
		log.Fatal(err)
	}

	// Setup runner
	c := runner.Config{
		Concurrency: 1,
		QueueSize:   1,
	}
	r := runner.New(c)
	r.Start()

	// Add job to runner queue
	if err := r.Enqueue(jobA); err != nil {
		log.Fatal(err)
	}

	// Cancel job
	time.AfterFunc(time.Millisecond*1, func() {
		jobA.Cancel(errors.New("canceled by user"))
	})

	// Ask runner to stop
	time.AfterFunc(time.Millisecond*100, func() {
		r.Stop()
	})

	// Wait for runner to run jobs
	r.Wait()

	// Output:
	// A
	// on cancel: canceled by user
}

func ExampleRunner_cancel_future_job_by_ID() {
	// Create jobs
	jobA, err := job.New(&cancelableTask{&myTask{text: "A"}, "err A"})
	if err != nil {
		log.Fatal(err)
	}

	jobB, err := job.New(&cancelableTask{&myTask{text: "B"}, "err B"})
	if err != nil {
		log.Fatal(err)
	}

	// Setup runner
	c := runner.Config{
		Concurrency: 1,
		QueueSize:   1,
	}
	r := runner.New(c)
	r.Start()

	// Add jobs to runner queue
	if err := r.Enqueue(jobA, jobB); err != nil {
		log.Fatal(err)
	}

	// Cancel a job by ID
	time.AfterFunc(time.Millisecond*1, func() {
		r.Cancel(jobB.ID())
	})

	// Ask runner to stop
	time.AfterFunc(time.Millisecond*200, func() {
		r.Stop()
	})

	// Wait for runner to run jobs
	r.Wait()

	// Output:
	// A
	// A
	// on cancel: err B
}

type defaultTask struct {
	num int
}

func (t *defaultTask) Run(*job.Job) (interface{}, error) {
	if t.num < 0 {
		return nil, errors.New("negative number")
	}

	return t.num + 1, nil
}

func (t *defaultTask) OnSuccess(result interface{}) {
	fmt.Println("on success result:", result.(int))
}

func (t *defaultTask) OnError(err error) {
	fmt.Println("on error:", err)
}

func (t *defaultTask) OnCancel(err error) {
	fmt.Println("on cancel:", err)
}

type cancelableTask struct {
	*myTask
	err string
}

func (c *cancelableTask) Run(j *job.Job) (interface{}, error) {
	return c.myTask.Run(j)
}

func (c *cancelableTask) OnCancel(err error) {
	fmt.Println("on cancel:", c.err)
}

type myTask struct {
	text string
}

func (t *myTask) Run(j *job.Job) (interface{}, error) {
	if t.text == "err" {
		return nil, errors.New("err")
	}

	for i := 0; i < 2; i++ {
		if j.IsCanceled() {
			break
		}
		fmt.Println(t.text)
		time.Sleep(time.Millisecond * 60)
	}

	return nil, nil
}
