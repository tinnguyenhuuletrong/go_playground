package channelplay

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func stringGenerator(prefix string, count uint) <-chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)
		for i := uint(0); i < count; i++ {
			ch <- fmt.Sprintf("%s-%d", prefix, i)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	return ch
}

func createWorker(id int, ctx context.Context, inp <-chan string) {
	for {
		select {
		case tmp, ok := <-inp:
			if !ok {
				fmt.Printf("[%d] inp channel closed.....\n", id)
				return
			}
			fmt.Printf("[%d] processValue %s\n", id, tmp)
		case <-ctx.Done():
			fmt.Printf("[%d] ctx closed.....\n", id)
			return
		}
	}
}

// Given
//	- Fanout model. 1 worker closed by context
// Conclusion:
// 	- In case 1 worker closed. Others still working. Job still dispatching normally. No missing

func Play_Channel_Ctx_Close() {

	generator := stringGenerator("it", 1e2)

	ctx := context.Background()
	ctx1, cancel1 := context.WithTimeout(ctx, time.Second*2)
	ctx2 := ctx
	defer cancel1()

	go createWorker(1, ctx1, generator)
	go createWorker(2, ctx2, generator)

	// Wait for Ctrl + C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
