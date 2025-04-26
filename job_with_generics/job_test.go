package job_with_generics_test

import (
	"log"
	"testing"

	"github.com/andreiavrammsd/workexec/job_with_generics"
)

func TestNew(t *testing.T) {
	taskJob, err := job_with_generics.New[int64](&fibonacciTask{n: 1})

	if taskJob == nil {
		log.Fatal("expected job")
	}

	if err != nil {
		t.Error("expected no error")
	}

	if taskJob.ID() == job_with_generics.ID("") {
		t.Error("expected job ID")
	}
}

func TestJob_CancelWithNotCancelableTask(t *testing.T) {
	taskJob, err := job_with_generics.New[string](&normalTask{})

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

func (n *normalTask) Run(*job_with_generics.Job[string]) (string, error) {
	return "", nil
}
