package runner_test

import (
	"testing"
	"time"

	"github.com/andreiavrammsd/jobrunner/job"
	"github.com/andreiavrammsd/jobrunner/runner"
)

func TestRunner_EnqueueWhenRunnerIsStopped(t *testing.T) {
	r, err := runner.New(runner.Config{})
	if err != nil {
		t.Fatal(err)
	}

	r.Stop()

	testJob, err := job.New(&task{})
	if err != nil {
		t.Fatal(err)
	}

	err = r.Enqueue(testJob)
	if err == nil {
		t.Error("expected runner is stopped error")
	}
}

func TestRunner_StopMultipleTimes(t *testing.T) {
	r, err := runner.New(runner.Config{})
	if err != nil {
		t.Fatal(err)
	}

	r.Stop()
	r.Stop()

	r.Wait()
}

func TestRunner_StopWithCancelRunningJobs(t *testing.T) {
	r, err := runner.New(runner.Config{
		Concurrency: 1,
		QueueSize:   1,
	})
	if err != nil {
		t.Fatal(err)
	}

	testJob, err := job.New(&task{duration: time.Millisecond * 10})
	if err != nil {
		t.Fatal(err)
	}

	err = r.Enqueue(testJob, testJob)

	r.Stop()

	r.Wait()
}

type task struct {
	duration time.Duration
}

func (t *task) Run(*job.Job) (interface{}, error) {
	if t.duration > 0 {
		time.Sleep(t.duration)
	}
	return nil, nil
}
