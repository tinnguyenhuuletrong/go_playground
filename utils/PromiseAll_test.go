package utils

import (
	"log"
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestPromiseAll1(t *testing.T) {

	inputs := []float64{1, 2, 3, 4}
	process := func(inp float64) float64 {
		return math.Pow(inp, 2)
	}

	all_res, err := PromiseAll(inputs, process)

	t.Log(all_res)
	t.Log(err)
}

func TestPromiseAll2(t *testing.T) {

	rand.Seed(0)

	inputs := []float64{1, 2, 3, 4}
	process := func(inp float64) float64 {
		sleepDuration := time.Duration(math.Ceil((rand.Float64())*5) * float64(time.Second))
		log.Println(sleepDuration)
		time.Sleep(sleepDuration)
		return math.Exp(inp)
	}

	all_res, err := PromiseAll(inputs, process)

	t.Log(all_res)
	t.Log(err)
}

func TestPromiseAll3(t *testing.T) {

	rand.Seed(0)

	inputs := []float64{1, 2, 3, 4}
	process := func(inp float64) float64 {
		if inp == 2 {
			panic("🤖 i don't like number 4 🤖")
		}
		sleepDuration := time.Duration(math.Ceil((rand.Float64())*5) * float64(time.Second))
		log.Println(sleepDuration)
		time.Sleep(sleepDuration)
		return math.Exp(inp)
	}

	all_res, err := PromiseAll(inputs, process)

	t.Log(all_res)
	t.Log(err)
}
