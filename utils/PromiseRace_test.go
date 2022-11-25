package utils

import (
	"context"
	"log"
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPromiseRace1(t *testing.T) {

	rand.Seed(0)

	inputs := []float64{1, 2, 3, 4}
	process := func(inp float64) float64 {
		sleepDuration := time.Duration(inp * float64(time.Second))
		log.Println(sleepDuration)
		time.Sleep(sleepDuration)
		return inp
	}

	final_res, err := PromiseRace(inputs, process)

	assert.Equal(t, final_res, float64(1))
	t.Log(final_res)
	t.Log(err)

	// check no panic b/c sent to close channel
	time.Sleep(5 * time.Second)
}

func TestPromiseRace2(t *testing.T) {

	rand.Seed(0)

	inputs := []float64{4, 3, 2, 1}
	process := func(inp float64) float64 {
		sleepDuration := time.Duration(inp * float64(time.Second))
		log.Println(sleepDuration)
		time.Sleep(sleepDuration)
		return inp
	}

	final_res, err := PromiseRace(inputs, process)

	assert.Equal(t, final_res, float64(1))
	t.Log(final_res)
	t.Log(err)
}

func TestPromiseRaceWithContext1(t *testing.T) {
	rand.Seed(0)

	inputs := []float64{1, 2, 3, 4}
	process := func(inp float64) float64 {
		return math.Pow(inp, 2)
	}

	ctx := context.Background()
	all_res, err := PromiseRaceWithCtx(ctx, inputs, process)

	t.Log(all_res)
	t.Log(err)
}

func TestPromiseRaceWithContext2(t *testing.T) {
	rand.Seed(0)

	inputs := []float64{1, 2, 3, 4}
	process := func(inp float64) float64 {
		sleepDuration := time.Duration(math.Ceil((rand.Float64())*100) * float64(time.Second))
		log.Println(sleepDuration)
		time.Sleep(sleepDuration)
		return math.Pow(inp, 2)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	all_res, err := PromiseAllWithCtx(ctx, inputs, process)

	t.Log(all_res)
	t.Log(err)
}
