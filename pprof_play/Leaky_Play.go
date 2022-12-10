package pprof_play

// Guide
// https://jvns.ca/blog/2017/09/24/profiling-go-with-pprof/
// https://medium.com/@wung_s/investigating-goroutine-leak-a10360daae1a

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	// Magic for pprof
	_ "net/http/pprof"
	"sync"
	"time"
)

func memLeakyFunction(wg *sync.WaitGroup) {
	defer wg.Done()
	s := make([]string, 3)
	for i := 0; i < 10000000; i++ {
		s = append(s, "magical pandas")
		if (i % 100000) == 0 {
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func neverReadUnBufferChannelLeakyFunction(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 10000000; i++ {
		runAsync() // 👈 problem point. Un-buffered channel never read -> lead to leak goroutine (runAsync 2 never print)

		// Fixed
		// if err := runAsync(); err != nil {
		//   log.Print(err)
		// }
	}
}

func runAsync() <-chan error {
	ch := make(chan error)
	go func() {
		defer close(ch)
		log.Printf("runAsync 1")
		ch <- doSomething()
		log.Printf("runAsync 2")
	}()
	return ch
}

func doSomething() error {
	return errors.New("OoO")
}

func Leaky_Play() {
	// pprof handler
	// 			http://localhost:6060/debug/pprof/goroutine
	// 			http://localhost:6060/debug/pprof/heap
	// 			http://localhost:6060/debug/pprof/threadcreate
	// 			http://localhost:6060/debug/pprof/block
	// 			http://localhost:6060/debug/pprof/mutex
	// 			and also 2 more: the CPU profile and the CPU trace.
	// 			http://localhost:6060/debug/pprof/profile
	// 			http://localhost:6060/debug/pprof/trace?seconds=5
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	fmt.Println("hello world")

	fmt.Println(`
Visit:
	- http://localhost:6060/debug/pprof/
	`)

	var wg sync.WaitGroup
	wg.Add(1)
	go memLeakyFunction(&wg)

	wg.Add(1)
	go neverReadUnBufferChannelLeakyFunction(&wg)

	wg.Wait()
}

// Guide 1:
// Live view on web
//		- http://localhost:6060/debug/pprof/

// Guide 2:
// Dump help -> analyze on pc
//  -	curl http://localhost:6060/debug/pprof/heap --output heap1.out
//  - go tool pprof -http=":9090" heap1.out
