package rle

import (
	"testing"

	oproto "github.com/dparrish/openinstrument/proto"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestEncodeDouble(c *C) {
	// Create pipeline
	input := &oproto.ValueStream{}

	push := func(timestamp uint64, value float64) {
		input.Value = append(input.Value, &oproto.Value{
			Timestamp:   timestamp,
			DoubleValue: value,
		})
	}
	push(1, 1.1)
	push(2, 1.1)
	push(3, 1.2)
	push(4, 1.2)
	push(5, 1.3)

	output := Encode(input)

	test := func(output *oproto.ValueStream, index int, startTimestamp, endTimestamp uint64, testValue float64) {
		value := output.Value[index]
		c.Assert(value.DoubleValue, Equals, testValue)
		c.Assert(value.Timestamp, Equals, startTimestamp)
		c.Assert(value.EndTimestamp, Equals, endTimestamp)
	}

	test(output, 0, 1, 2, 1.1)
	test(output, 1, 3, 4, 1.2)
	test(output, 2, 5, 5, 1.3)
}

func (s *MySuite) TestEncodeString(c *C) {
	// Create pipeline
	input := &oproto.ValueStream{}

	push := func(timestamp uint64, value string) {
		input.Value = append(input.Value, &oproto.Value{
			Timestamp:   timestamp,
			StringValue: value,
		})
	}
	push(1, "Hello")
	push(2, "Hello")
	push(3, "Hello")
	push(4, "Goodbye")
	push(5, "Friend")

	output := Encode(input)

	test := func(output *oproto.ValueStream, index int, startTimestamp, endTimestamp uint64, testValue string) {
		value := output.Value[index]
		c.Assert(value.StringValue, Equals, testValue)
		c.Assert(value.Timestamp, Equals, startTimestamp)
		c.Assert(value.EndTimestamp, Equals, endTimestamp)
	}

	test(output, 0, 1, 3, "Hello")
	test(output, 1, 4, 4, "Goodbye")
	test(output, 2, 5, 5, "Friend")
}

func (s *MySuite) BenchmarkEncodeDouble(c *C) {
	// Create pipeline
	input := &oproto.ValueStream{}
	// Number of repetitions of each value
	const blockSize = 100

	for i := 1; i <= 1000000; i++ {
		input.Value = append(input.Value, &oproto.Value{
			Timestamp:   uint64(i),
			DoubleValue: float64(i - (i % blockSize)),
		})
	}

	for run := 0; run < c.N; run++ {
		Encode(input)
	}
}
