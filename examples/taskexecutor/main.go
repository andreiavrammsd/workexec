package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/andreiavrammsd/workexec/future"
	"github.com/andreiavrammsd/workexec/taskexecutor"
)

func main() {
	config := taskexecutor.Config{
		Concurrency: 6,
		QueueSize:   6,
	}

	taskExecutor := taskexecutor.New(config)
	taskExecutor.Start()

	go func() {
		for {
			n, err := rand.Int(rand.Reader, big.NewInt(4))
			if err != nil {
				log.Println("rand error:", err)
				continue
			}

			futureTask, err := future.New(&task{in: int(n.Int64())})
			if err != nil {
				log.Println("future error:", err)
				continue
			}

			if n.Int64() == 3 {
				time.AfterFunc(time.Millisecond*100, func() {
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
		time.AfterFunc(time.Second*3, func() {
			taskExecutor.Stop()
		})
	}()

	taskExecutor.Wait()
}

type task struct {
	in int
}

func (t *task) OnSuccess(result interface{}) {
	fmt.Println("success:", t.in, result, time.Now().Unix())
}

func (t *task) OnError(err error) {
	fmt.Println("error:", t.in, err, time.Now().Unix())
}

func (t *task) OnCancel() {
	fmt.Println("canceled:", t.in, time.Now().Unix())
}

func (t *task) Run(func() bool) (interface{}, error) {
	if t.in == 0 {
		return nil, errors.New("n == 0")
	}

	time.Sleep(time.Second)

	return t.in + 1, nil
}
