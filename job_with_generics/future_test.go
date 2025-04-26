package job_with_generics_test

import (
	"errors"
	"testing"

	"github.com/andreiavrammsd/workexec/job_with_generics"
)

func TestFuture_Wait(t *testing.T) {
	futureJob, err := job_with_generics.New[int](&task{in: 1})
	if err != nil {
		t.Fatal(err)
	}

	future := futureJob.Run()
	future.Wait()

	expected := 2
	actual := future.Result()
	if actual != expected {
		t.Errorf("got %d, expected %d", actual, expected)
	}

	if future.Error() != nil {
		t.Fatal("error not expected")
	}
}

func TestFuture_Result(t *testing.T) {
	futureJob, err := job_with_generics.New[int](&task{in: 1})
	if err != nil {
		t.Fatal(err)
	}

	future := futureJob.Run()

	expected := 2
	actual := future.Result()
	if actual != expected {
		t.Errorf("got %d, expected %d", actual, expected)
	}

	if future.Error() != nil {
		t.Fatal("error not expected")
	}
}

func TestFuture_Error(t *testing.T) {
	futureJob, err := job_with_generics.New[int](&task{in: 0})
	if err != nil {
		t.Fatal(err)
	}

	future := futureJob.Run()

	if future.Result() != 0 {
		t.Error("expected nil result")
	}

	if future.Error() == nil {
		t.Error("expected error")
	}
}

func TestFuture_IsCanceled(t *testing.T) {
	futureJob, err := job_with_generics.New[int](&task{in: 0})
	if err != nil {
		t.Fatal(err)
	}

	futureJob.Cancel(nil)

	future := futureJob.Run()
	future.Wait()

	if !future.IsCanceled() {
		t.Error("expected job to be canceled")
	}
}

type task struct {
	in int
}

func (t *task) OnCancel(err error) {
}

func (t *task) Run(*job_with_generics.Job[int]) (int, error) {
	if t.in == 0 {
		return 0, errors.New("err")
	}
	return t.in + 1, nil
}
