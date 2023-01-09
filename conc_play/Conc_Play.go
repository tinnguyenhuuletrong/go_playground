package conc_play

import (
	"log"
	"math"
	"time"

	"github.com/sourcegraph/conc"
	"github.com/sourcegraph/conc/iter"
	"github.com/sourcegraph/conc/pool"
)

func doSomeThingWithPanic(v int) {
	if v == 9 {
		log.Panicln("i don't like number 9")
	}

	time.Sleep(10 * time.Millisecond)
	log.Println("done", v)
}

func doSomeThingWithWait(v int, d time.Duration) {
	time.Sleep(d)
	log.Println("done", v)
}

func doSomeThing(v int) {
	time.Sleep(10 * time.Millisecond)
	log.Println("done", v)
}

func Conc_BetterWaitGroup_Play() {

	log.Println("start")

	var wg conc.WaitGroup
	for i := 0; i < 20; i++ {
		inp := i
		wg.Go(func() {
			doSomeThing(inp)
		})
	}

	/*

		number 9 -> die
		others -> still running
		nice stacktrace
	*/

	wg.Wait()
	log.Println("end")
}

func Conc_LimitedWorkerPool_Play() {
	POOL_SIZE := 2
	p := pool.New().WithMaxGoroutines(POOL_SIZE)
	log.Println("start, poolSize", POOL_SIZE)

	for i := 0; i < 20; i++ {
		inp := i
		p.Go(func() {
			doSomeThingWithWait(inp, 2*time.Second)
		})
	}
	p.Wait()
	log.Println("end")
}

func Conc_Iter_Play() {
	log.Println("begin")

	data := make([]int, 100)
	for i := 0; i < 100; i++ {
		data[i] = i
	}

	log.Println("data:", data)

	log.Println("iter.ForEach")
	iter.ForEach(data, func(v *int) {
		log.Println("\t", *v)
	})

	log.Println("iter.Map power of 2")
	data1 := iter.Map(data, func(v *int) int {
		return int(math.Pow(float64((*v)), 2.0))
	})
	log.Println("iter.Map power of 2 res", data1)

	log.Println("end")
}
