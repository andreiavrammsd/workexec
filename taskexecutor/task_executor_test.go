package taskexecutor_test

import (
	"math"
	"runtime"
	"testing"

	taskexec "github.com/andreiavrammsd/workexec/taskexecutor"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	config := taskexec.Config{}
	taskExecutor := taskexec.New(config)

	assert.Equal(t, uint(math.Max(1, float64(runtime.NumCPU())-1)), taskExecutor.Concurrency())
}
