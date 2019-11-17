package promise_test

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/andreiavrammsd/jobrunner/promise"
)

func ExamplePromise_Async() {
	exec := &division{
		a: 4,
		b: 2,
	}

	promise.New(exec).Then(func() {
		fmt.Println("Done")
	}).Async()

	time.Sleep(time.Millisecond * 50)

	// Output:
	// Done
}

func ExamplePromise_Await() {
	exec := &division{
		a: 4,
		b: 2,
	}

	p := promise.New(exec).
		Then(func() {
			fmt.Println("Done")
		}).
		Error(func(err error) {
			fmt.Println("Error:", err)
		})

	err := p.Await()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(exec.result)

	// Output:
	// Done
	// 2
}

func ExamplePromise_multiple_executors() {
	p1 := &printer{text: "1\n", w: os.Stdout}
	p2 := &printer{text: "2\n", w: os.Stdout}
	p3 := &printer{text: "3\n", w: os.Stdout}
	p4 := &printer{text: "4\n", w: os.Stdout}

	if err := promise.New(p1, p2, p3, p4).Await(); err != nil {
		log.Fatal(err)
	}

	// Unordered output:
	// 1
	// 2
	// 3
	// 4
}

type printer struct {
	text string
	w    io.StringWriter
	l    sync.Mutex
}

func (p *printer) Execute() (err error) {
	p.l.Lock()
	_, err = p.w.WriteString(p.text)
	p.l.Unlock()
	return
}
