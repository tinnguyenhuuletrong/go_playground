package faninout

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// https://medium.com/better-programming/cloud-native-patterns-illustrated-fan-in-and-fan-out-daf77455703c

func stringGenerator(prefix string, count uint) <-chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)
		for i := uint(0); i < count; i++ {
			ch <- fmt.Sprintf("%s-%d", prefix, i)
		}
	}()

	return ch
}

func stringProcessor(id int, source <-chan string) <-chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)
		for data := range source {
			ch <- fmt.Sprintf("P-%d processed message: %s -> %s", id, data, strings.ToUpper(data))
			millisecs := rand.Intn(901) + 100
			time.Sleep(time.Duration(millisecs * int(time.Millisecond)))
		}
	}()

	return ch
}

func Play_FanInOut() {
	const (
		processorCount = 10
	)
	splitSources := FanOut(stringGenerator("itm", 100), processorCount)
	procSources := make([]ChanAnyReadOnly[string], 0, processorCount)

	for i, src := range splitSources {
		procSources = append(procSources, stringProcessor(i, src))
	}

	for out := range FanIn(procSources) {
		fmt.Println(out)
	}
}
