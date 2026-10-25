package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/andreiavrammsd/workexec/future"
	"github.com/andreiavrammsd/workexec/taskexecutor"
)

const (
	defaultConcurrency = 6
	maxRandomValue     = 4
	cancelAtValue      = 3
	cancelDelay        = 100 * time.Millisecond
	stopDelay          = 3 * time.Second
)

var errNZero = errors.New("n == 0")

func main() {
	config := taskexecutor.Config{
		Concurrency: defaultConcurrency,
		QueueSize:   defaultConcurrency,
	}

	taskExecutor := taskexecutor.New(config)
	taskExecutor.Start()

	go func() {
		for {
			n, err := rand.Int(rand.Reader, big.NewInt(maxRandomValue))
			if err != nil {
				log.Println("rand error:", err)

				continue
			}

			futureTask, err := future.New(&task{in: int(n.Int64())})
			if err != nil {
				log.Println("future error:", err)

				continue
			}

			if n.Int64() == cancelAtValue {
				time.AfterFunc(cancelDelay, func() {
					futureTask.Cancel()
				})
			}

			if err := taskExecutor.Submit(futureTask); err != nil {
				log.Println(err)

				return
			}
		}
	}()

	go func() {
		time.AfterFunc(stopDelay, func() {
			taskExecutor.Stop()
		})
	}()

	taskExecutor.Wait()
}

type task struct {
	in int
}

func (t *task) OnSuccess(result any) {
	_, _ = fmt.Fprintf(os.Stdout, "success: %d %v %d\n", t.in, result, time.Now().Unix())
}

func (t *task) OnError(err error) {
	_, _ = fmt.Fprintf(os.Stdout, "error: %d %v %d\n", t.in, err, time.Now().Unix())
}

func (t *task) OnCancel() {
	_, _ = fmt.Fprintf(os.Stdout, "canceled: %d %d\n", t.in, time.Now().Unix())
}

func (t *task) Run(func() bool) (any, error) {
	if t.in == 0 {
		return nil, errNZero
	}

	time.Sleep(time.Second)

	return t.in + 1, nil
}
