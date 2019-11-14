package job

import (
	"errors"
	"time"
)

// Job options
type Option interface {
	Apply(*Job)
}

type WithTimeout struct {
	Timeout time.Duration
}

func (o *WithTimeout) Apply(j *Job) {
	time.AfterFunc(o.Timeout, func() {
		j.Cancel(errors.New("timeout"))
	})
}
