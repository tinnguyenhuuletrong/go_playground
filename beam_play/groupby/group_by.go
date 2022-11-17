package groupby

import (
	"context"
	"flag"
	"fmt"
	"sort"

	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/log"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/register"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/x/beamx"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/x/debug"
)

type stringPair struct {
	K, V string
}

func splitStringPair(e stringPair) (string, string) {
	return e.K, e.V
}

func init() {
	// Register DoFn.
	register.Function1x2(splitStringPair)

	register.Function3x1(formatCoGBKResults)
	// 1 input of type string => Iter1[string]
	register.Iter1[string]()
}

func formatCoGBKResults(key string, emailIter, phoneIter func(*string) bool) string {
	var s string
	var user_emails, user_phones []string
	for emailIter(&s) {
		user_emails = append(user_emails, s)
	}
	for phoneIter(&s) {
		user_phones = append(user_phones, s)
	}
	// Values have no guaranteed order, sort for deterministic output.
	sort.Strings(user_emails)
	sort.Strings(user_phones)
	return fmt.Sprintf("%s; Emails: %s; Phones: %s", key, user_emails, user_phones)
}

// createAndSplit is a helper function that creates
func createAndSplit(s beam.Scope, input []stringPair) beam.PCollection {
	initial := beam.CreateList(s, input)
	return beam.ParDo(s, splitStringPair, initial)
}

func StartGroupBy() {
	flag.Parse()
	beam.Init()

	ctx := context.Background()

	p := beam.NewPipeline()
	s := p.Root()

	var emailSlice = []stringPair{
		{"amy", "amy@example.com"},
		{"amy", "amy1@example.com"},
		{"amy", "amy2@example.com"},
		{"amy", "amy3@example.com"},
		{"carl", "carl@example.com"},
		{"carl", "carl1@example.com"},
		{"julia", "julia@example.com"},
		{"carl", "carl@email.com"},
		{"carl", "carl2@email.com"},
		{"carl", "carl3@email.com"},
		{"carl", "carl4@email.com"},
		{"carl", "carl5@email.com"},
	}

	var phoneSlice = []stringPair{
		{"amy", "111-222-3333"},
		{"james", "222-333-4444"},
		{"amy", "333-444-5555"},
		{"carl", "444-555-6666"},
	}

	// 1. Convert to (k,v) PCollection
	emails := createAndSplit(s, emailSlice)
	phones := createAndSplit(s, phoneSlice)

	debug.Printf(s, "emails: %v", emails)
	debug.Printf(s, "phones: %v", phones)

	// 2. Simple Group by key
	emailGroupBy := beam.GroupByKey(s, emails)
	debug.Printf(s, "email_group_by_key: %v", emailGroupBy)

	// 3. CoGroupBy aka Group by multiple collection
	//	Note: Can not direct display joinRes. Need format it  - see formatCoGBKResults
	joinRes := beam.CoGroupByKey(s, emails, phones)
	formatedRes := beam.ParDo(s, formatCoGBKResults, joinRes)
	debug.Printf(s, "email_join_phones_group_by: %v", formatedRes)

	if err := beamx.Run(ctx, p); err != nil {
		log.Exitf(ctx, "Failed to execute job: %v", err)
	}
}
