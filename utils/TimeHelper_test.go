package utils

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHumanrize(t *testing.T) {
	var (
		from time.Time
		to   time.Time
		err  error
		res  string
	)

	from, err = time.Parse(time.RFC3339, "2021-11-12T18:49:00+07:00")
	assert.Equal(t, err, nil, "parse error")
	to, err = time.Parse(time.RFC3339, "2021-11-12T18:50:00+07:00")
	assert.Equal(t, err, nil, "parse error")
	log.Printf("t: %v %v\n", from, to)
	res = HumanizeDurationAsString(from, to)
	assert.Equal(t, res, "1 min")

	from, err = time.Parse(time.RFC3339, "2021-11-12T18:49:00+07:00")
	assert.Equal(t, err, nil, "parse error")
	to, err = time.Parse(time.RFC3339, "2021-11-12T18:49:05+07:00")
	assert.Equal(t, err, nil, "parse error")
	res = HumanizeDurationAsString(from, to)
	assert.Equal(t, res, "5 sec")

	from, err = time.Parse(time.RFC3339, "2021-11-12T18:49:00+07:00")
	assert.Equal(t, err, nil, "parse error")
	to, err = time.Parse(time.RFC3339, "2021-11-12T19:49:05+07:00")
	assert.Equal(t, err, nil, "parse error")
	res = HumanizeDurationAsString(from, to)
	assert.Equal(t, res, "1 hour")

	from, err = time.Parse(time.RFC3339, "2021-11-12T18:49:00+07:00")
	assert.Equal(t, err, nil, "parse error")
	to, err = time.Parse(time.RFC3339, "2021-11-19T19:49:05+07:00")
	assert.Equal(t, err, nil, "parse error")
	res = HumanizeDurationAsString(from, to)
	assert.Equal(t, res, "7 day")
}
