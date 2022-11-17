package combine

import (
	"context"
	"flag"
	"math"

	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/log"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/x/beamx"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/x/debug"
)

type averageFn struct{}
type averageAccum struct {
	Count, Sum float64
}

func (fn *averageFn) CreateAccumulator() averageAccum {
	return averageAccum{0, 0}
}

func (fn *averageFn) AddInput(a averageAccum, v float64) averageAccum {
	return averageAccum{Count: a.Count + 1, Sum: a.Sum + v}
}

func (fn *averageFn) MergeAccumulators(a, v averageAccum) averageAccum {
	return averageAccum{Count: a.Count + v.Count, Sum: a.Sum + v.Sum}
}

func (fn *averageFn) ExtractOutput(a averageAccum) float64 {
	if a.Count == 0 {
		return math.NaN()
	}
	return float64(a.Sum) / float64(a.Count)
}

func doSum(a, b float64) float64 {
	return a + b
}

func StartCombineSum() {
	flag.Parse()
	beam.Init()

	ctx := context.Background()

	p := beam.NewPipeline()
	s := p.Root()

	samples := []float64{1.0, 2, 3.9, 4, 5, 6, 7}
	input := beam.CreateList(s, samples)

	debug.Printf(s, "input: %v", input)
	sumVal := beam.Combine(s, doSum, input)
	debug.Printf(s, "sum: %v", sumVal)

	avgVal := beam.Combine(s, &averageFn{}, input)
	debug.Printf(s, "avg: %v", avgVal)

	if err := beamx.Run(ctx, p); err != nil {
		log.Exitf(ctx, "Failed to execute job: %v", err)
	}
}
