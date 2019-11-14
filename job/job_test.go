package job_test

import (
	"log"
	"testing"

	"github.com/andreiavrammsd/jobrunner/job"
)

func TestNew(t *testing.T) {
	j, err := job.New(&fibonacciTask{n: 1})

	if j == nil {
		log.Fatal("expected job")
	}

	if err != nil {
		t.Error("expected no error")
	}

	if j.ID() == job.ID("") {
		t.Error("expected job ID")
	}
}

func TestNewWithError(t *testing.T) {
	j, err := job.New(nil)

	if j != nil {
		t.Error("expected no job")
	}

	if err == nil {
		t.Error("expected task function not passed error")
	}
}
