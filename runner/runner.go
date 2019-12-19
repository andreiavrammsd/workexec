// Package runner handles execution of jobs at desired concurrency level.
package runner

import (
	"errors"
	"sync"

	"github.com/andreiavrammsd/workexec/job"
	"github.com/cespare/xxhash/v2"
)

const (
	concurrency = 1024
	queueSize   = 1024
)

// Config allows setup of runner.
type Config struct {
	Concurrency uint
	QueueSize   uint
}

// Runner represents a manager of jobs.
type Runner struct {
	concurrency uint
	queue       chan *job.Job
	stop        chan struct{}
	running     map[uint64]*job.Job
	toCancel    map[uint64]struct{}
	wait        chan struct{}
	state       state
	lock        sync.RWMutex
}

// Status represents the current state of the runner, regarding number of routines
// it is working on and count of jobs currently running.
type Status struct {
	Concurrency uint
	RunningJobs uint
}

// state of runner
type state uint

const (
	// stopped means the runner is not working on any routine
	stopped state = iota
	// running means the runner has routines working on
	running
)

// Start starts the runner routines
func (r *Runner) Start() {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.state == running {
		return
	}
	r.state = running

	for i := uint(0); i < r.concurrency; i++ {
		go r.run()
	}
}

// Stop asks the runner to stop all jobs from running.
func (r *Runner) Stop() {
	r.lock.Lock()
	if r.state == stopped {
		r.lock.Unlock()
		return
	}
	r.state = stopped
	r.lock.Unlock()

	for _, j := range r.running {
		j.Cancel(errors.New("runner was stopped"))
	}

	for i := uint(0); i < r.concurrency; i++ {
		r.stop <- struct{}{}
	}
}

// Enqueue puts jobs to the runner queue.
func (r *Runner) Enqueue(jobs ...*job.Job) error {
	r.lock.RLock()
	isStopped := r.state == stopped
	r.lock.RUnlock()

	if isStopped {
		return errors.New("runner is stopped")
	}

	for i := 0; i < len(jobs); i++ {
		r.queue <- jobs[i]
	}

	return nil
}

// Wait blocks until runner is done with running all the queued jobs.
func (r *Runner) Wait() {
	r.lock.RLock()
	isStopped := r.state == stopped
	r.lock.RUnlock()

	if isStopped {
		return
	}

	<-r.wait
}

// Cancel asks a job (by given id) to stop.
func (r *Runner) Cancel(id job.ID) {
	r.lock.Lock()
	r.cancel(id)
	r.lock.Unlock()
}

// ScaleUp increases concurrency by starting new worker routines.
func (r *Runner) ScaleUp(count uint) {
	if count == 0 {
		return
	}

	r.lock.Lock()
	r.concurrency += count
	r.lock.Unlock()

	for i := uint(0); i < count; i++ {
		go r.run()
	}
}

// ScaleDown decreases concurrency by asking routines to stop.
func (r *Runner) ScaleDown(count uint) {
	if count == 0 {
		return
	}

	r.lock.Lock()
	if int(r.concurrency)-int(count) >= 0 {
		r.concurrency -= count
	} else {
		r.concurrency = 0
	}
	r.lock.Unlock()

	for i := uint(0); i < count; i++ {
		r.stop <- struct{}{}
	}
}

// Status returns runner state
func (r *Runner) Status() Status {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return Status{
		Concurrency: r.concurrency,
		RunningJobs: uint(len(r.running)),
	}
}

func (r *Runner) run() {
	for {
		select {
		case j := <-r.queue:
			hash := hash(j.ID())

			r.lock.Lock()

			// Add to running jobs
			r.running[hash] = j

			// Check if scheduled for cancellation
			if _, cancel := r.toCancel[hash]; cancel {
				delete(r.toCancel, hash)
				r.cancel(j.ID())
			}

			r.lock.Unlock()

			j.Run().Wait()

			r.lock.Lock()
			delete(r.running, hash)
			r.lock.Unlock()
		case <-r.stop:
			r.lock.RLock()
			if r.state == stopped && len(r.running) == 0 {
				r.wait <- struct{}{}
			}
			r.lock.RUnlock()

			return
		}
	}
}

func (r *Runner) cancel(id job.ID) {
	hash := hash(id)

	// Cancel now if running
	j, ok := r.running[hash]
	if ok {
		j.Cancel(errors.New("canceled by runner"))
		return
	}

	// Schedule to be canceled before run
	r.toCancel[hash] = struct{}{}
}

// New creates a new job runner.
func New(c Config) *Runner {
	if c.Concurrency == 0 {
		c.Concurrency = concurrency
	}
	if c.QueueSize == 0 {
		c.QueueSize = queueSize
	}

	return &Runner{
		concurrency: c.Concurrency,
		queue:       make(chan *job.Job, c.QueueSize),
		stop:        make(chan struct{}, c.QueueSize),
		running:     make(map[uint64]*job.Job),
		toCancel:    make(map[uint64]struct{}),
		wait:        make(chan struct{}, c.Concurrency),
		state:       stopped,
	}
}

func hash(id job.ID) uint64 {
	return xxhash.Sum64String(string(id))
}
