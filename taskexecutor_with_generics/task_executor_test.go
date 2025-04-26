package taskexecutor_with_generics

import (
	"math"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	config := Config{}
	taskExecutor := New(config)

	assert.Equal(t, uint(math.Max(1, float64(runtime.NumCPU())-1)), taskExecutor.concurrency)
}
