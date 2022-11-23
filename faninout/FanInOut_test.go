package faninout

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func testStringGenerator(prefix string, count uint, interval time.Duration) <-chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)
		for i := uint(0); i < count; i++ {
			ch <- fmt.Sprintf("%s-%d", prefix, i)
			time.Sleep(interval)
		}
	}()

	return ch
}

func TestFanOut(t *testing.T) {
	strGen := testStringGenerator("t", 100, time.Duration(0))
	chunks := FanOut(strGen, 3)

	var counter atomic.Int32
	counter.Store(0)
	var wg sync.WaitGroup
	wg.Add(len(chunks))

	for id, c := range chunks {
		go func(id int, c <-chan string) {
			defer wg.Done()
			for v := range c {
				counter.Add(1)
				t.Logf("worker-%d: %v\n", id, v)
			}
		}(id, c)
	}

	wg.Wait()
	final := counter.Load()
	t.Logf("Final: %d", final)
	assert.Equal(t, final, int32(100), "not equal")
}

func TestFanIn(t *testing.T) {
	strGen := testStringGenerator("t", 100, time.Duration(0))
	chunks := FanOut(strGen, 3)

	var counter atomic.Int32
	counter.Store(0)
	out := FanIn(chunks)

	for v := range out {
		counter.Add(1)
		t.Logf("out: %v\n", v)
	}

	final := counter.Load()
	t.Logf("Final: %d", final)
	assert.Equal(t, final, int32(100), "not equal")
}
