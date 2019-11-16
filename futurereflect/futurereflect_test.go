package futurereflect_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andreiavrammsd/jobrunner/futurereflect"
)

type lines struct {
	out string
}

func (l *lines) str(lines []string) {
	l.out = strings.Join(lines, "\n")
}

func (l *lines) delete() {
	l.out = ""
}

func Test(t *testing.T) {
	sum := func(a, b int) (int, error) {
		return a + b, nil
	}
	future := futurereflect.New(sum)(1, 2)

	result, err := future.Result()
	assert.Equal(t, 3, result)
	assert.NoError(t, err)

	ls := &lines{}

	futureLines := futurereflect.New(ls.str)([]string{"A", "B"})
	futureLines.Wait()
	assert.Equal(t, "A\nB", ls.out)

	futureDelete := futurereflect.New(ls.delete)()
	futureDelete.Wait()
	assert.Equal(t, "", ls.out)

	sumMany := func(a, b, c int) int {
		return a + b + c
	}
	futureSumMany := futurereflect.New(sumMany)(1, 2, 3)
	result, err = futureSumMany.Result()
	assert.Equal(t, 6, result.(int))
	assert.NoError(t, err)
}
