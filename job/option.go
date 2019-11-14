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
	Err     error
}

func (o *WithTimeout) Apply(j *Job) {
	o.Err = errors.New("timedout")
	time.AfterFunc(o.Timeout, func() {
		//j.Cancel(ErrTimeout)
		j.Cancel(o.Err)
	})
}
