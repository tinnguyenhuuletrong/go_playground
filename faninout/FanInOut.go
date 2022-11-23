package faninout

import "sync"

type ChanAnyReadOnly[T any] <-chan T
type ChanAny[T any] chan T

func FanOut[T any](source ChanAnyReadOnly[T], numWorker uint) []ChanAnyReadOnly[T] {
	res := make([]ChanAnyReadOnly[T], 0, numWorker)

	// n worker listen same source -> dispatch it
	for i := 0; i < int(numWorker); i++ {
		worker := make(chan T)
		go func() {
			defer close(worker)
			for v := range source {
				worker <- v
			}
		}()
		res = append(res, worker)
	}

	return res
}

func FanIn[T any](sources []ChanAnyReadOnly[T]) ChanAny[T] {
	out := make(ChanAny[T])

	var wg sync.WaitGroup
	wg.Add(len(sources))

	for _, source := range sources {
		go func(source ChanAnyReadOnly[T]) {
			defer wg.Done()
			for data := range source {
				out <- data
			}
		}(source)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
