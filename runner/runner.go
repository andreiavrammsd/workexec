package runner

import (
	"errors"
	"sync"

	"github.com/andreiavrammsd/jobrunner/job"
	"github.com/cespare/xxhash/v2"
)

const (
	queueSize   = 1024
	concurrency = 1024
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
	stopped     bool
	lock        sync.RWMutex
}

// Enqueue puts jobs to the runner queue.
func (r *Runner) Enqueue(jobs ...*job.Job) error {
	r.lock.RLock()
	stopped := r.stopped
	r.lock.RUnlock()

	if stopped {
		return errors.New("runner is stopped")
	}

	for i := 0; i < len(jobs); i++ {
		r.queue <- jobs[i]
	}

	return nil
}

// Wait blocks until runner is done with running all the queued jobs.
func (r *Runner) Wait() {
	<-r.wait
}

// Stop asks the runner to stop all jobs from running.
func (r *Runner) Stop() {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.stopped {
		return
	}

	r.stopped = true

	for _, j := range r.running {
		j.Cancel(errors.New("runner was stopped"))
	}

	for i := uint(0); i < r.concurrency; i++ {
		r.stop <- struct{}{}
	}
}

// Cancel asks a job (by given id) to stop.
func (r *Runner) Cancel(id job.ID) {
	r.lock.Lock()
	r.cancel(id)
	r.lock.Unlock()
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
			if r.stopped && len(r.running) == 0 {
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
func New(c Config) (*Runner, error) {
	if c.Concurrency == 0 {
		c.Concurrency = concurrency
	}
	if c.QueueSize == 0 {
		c.QueueSize = queueSize
	}

	r := &Runner{
		concurrency: c.Concurrency,
		queue:       make(chan *job.Job, c.QueueSize),
		stop:        make(chan struct{}, c.QueueSize),
		running:     make(map[uint64]*job.Job),
		toCancel:    make(map[uint64]struct{}),
		wait:        make(chan struct{}),
	}

	for i := uint(0); i < r.concurrency; i++ {
		go r.run()
	}

	return r, nil
}

func hash(id job.ID) uint64 {
	return xxhash.Sum64String(string(id))
}
