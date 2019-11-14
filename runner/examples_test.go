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
	// Create jobs
	jobA, err := job.New(&myTask{text: "A"})
	if err != nil {
		log.Fatal(err)
	}

	jobB, err := job.New(&myTask{text: "B"})
	if err != nil {
		log.Fatal(err)
	}

	// Setup runner
	c := runner.Config{
		Concurrency: 2,
		QueueSize:   2,
	}
	r, err := runner.New(c)
	if err != nil {
		log.Fatal(err)
	}

	// Add jobs to runner queue
	if err := r.Enqueue(jobA, jobB); err != nil {
		log.Fatal(err)
	}

	// Ask runner to stop
	time.AfterFunc(time.Millisecond*10, func() {
		r.Stop()
		r.Stop() // second stop will be ignored
	})

	// Wait for runner to run jobs
	r.Wait()

	// Add jobs to stopped runner
	if err := r.Enqueue(jobA); err != nil {
		fmt.Println(err)
	}

	// Unordered output:
	// A
	// B
	// runner is stopped
}

func ExampleRunner_cancel_running_job_by_ID() {
	// Create job
	jobA, err := job.New(&myTask{text: "A"})
	if err != nil {
		log.Fatal(err)
	}

	jobA.OnCancel(func(err error) {
		fmt.Println("canceled:", err)
	})

	// Setup runner
	c := runner.Config{
		Concurrency: 1,
		QueueSize:   1,
	}
	r, err := runner.New(c)
	if err != nil {
		log.Fatal(err)
	}

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
	// canceled: canceled by runner
}

func ExampleRunner_cancel_job() {
	// Create job
	jobA, err := job.New(&myTask{text: "A"})
	if err != nil {
		log.Fatal(err)
	}

	jobA.OnCancel(func(err error) {
		fmt.Println("canceled:", err)
	})

	// Setup runner
	c := runner.Config{
		Concurrency: 1,
		QueueSize:   1,
	}
	r, err := runner.New(c)
	if err != nil {
		log.Fatal(err)
	}

	// Add job to runner queue
	if err := r.Enqueue(jobA); err != nil {
		log.Fatal(err)
	}

	// Cancel a job by ID
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
	// canceled: canceled by user
}

func ExampleRunner_cancel_future_job_by_ID() {
	// Create jobs
	jobA, err := job.New(&myTask{text: "A"})
	if err != nil {
		log.Fatal(err)
	}

	jobB, err := job.New(&myTask{text: "B"})
	if err != nil {
		log.Fatal(err)
	}

	jobB.OnCancel(func(err error) {
		fmt.Println("canceled:", err)
	})

	// Setup runner
	c := runner.Config{
		Concurrency: 1,
		QueueSize:   1,
	}
	r, err := runner.New(c)
	if err != nil {
		log.Fatal(err)
	}

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
	// canceled: canceled by runner
}

func ExampleRunner_job_error() {
	// Create job
	jobA, err := job.New(&myTask{text: "err"})
	if err != nil {
		log.Fatal(err)
	}

	// Setup runner
	c := runner.Config{
		Concurrency: 1,
		QueueSize:   1,
	}
	r, err := runner.New(c)
	if err != nil {
		log.Fatal(err)
	}

	// Add job to runner queue
	if err := r.Enqueue(jobA); err != nil {
		log.Fatal(err)
	}

	jobA.OnError(func(err error) {
		fmt.Println("on error:", err)
	})

	// Ask runner to stop
	time.AfterFunc(time.Millisecond*10, func() {
		r.Stop()
	})

	// Wait for runner to run jobs
	r.Wait()

	// Output:
	// on error: err
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
		time.Sleep(time.Millisecond * 50)
	}

	return nil, nil
}
