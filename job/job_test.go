package job_test

import (
	"log"
	"testing"

	"github.com/andreiavrammsd/workexec/job"
)

func TestNew(t *testing.T) {
	taskJob, err := job.New(&fibonacciTask{n: 1})

	if taskJob == nil {
		log.Fatal("expected job")
	}

	if err != nil {
		t.Error("expected no error")
	}

	if taskJob.ID() == job.ID("") {
		t.Error("expected job ID")
	}
}

func TestNewWithError(t *testing.T) {
	taskJob, err := job.New(nil)

	if taskJob != nil {
		t.Error("expected no job")
	}

	if err == nil {
		t.Error("expected nil task passed to job error")
	}
}

func TestJob_CancelWithNotCancelableTask(t *testing.T) {
	taskJob, err := job.New(&normalTask{})

	if taskJob == nil {
		log.Fatal("expected job")
	}

	if err != nil {
		t.Error("expected no error")
	}

	taskJob.Cancel(nil)
	taskJob.Run()

	if taskJob.IsCanceled() {
		t.Error("expected task to not be canceled")
	}
}

type normalTask struct {
}

func (n *normalTask) Run(*job.Job) (interface{}, error) {
	return nil, nil
}
