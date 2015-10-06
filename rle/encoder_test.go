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
	input := make(chan *oproto.Value, 10)
	output := Encode(input)

	push := func(input chan *oproto.Value, timestamp uint64, value float64) {
		input <- &oproto.Value{
			Timestamp:   timestamp,
			DoubleValue: value,
		}
	}

	test := func(output <-chan *oproto.Value, startTimestamp, endTimestamp uint64, testValue float64) {
		value := <-output
		c.Assert(value.DoubleValue, Equals, testValue)
		c.Assert(value.Timestamp, Equals, startTimestamp)
		c.Assert(value.EndTimestamp, Equals, endTimestamp)
	}

	push(input, 1, 1.1)
	push(input, 2, 1.1)
	push(input, 3, 1.2)
	push(input, 4, 1.2)
	push(input, 5, 1.3)
	close(input)

	test(output, 1, 2, 1.1)
	test(output, 3, 4, 1.2)
	test(output, 5, 5, 1.3)
}

func (s *MySuite) TestEncodeString(c *C) {
	// Create pipeline
	input := make(chan *oproto.Value, 10)
	output := Encode(input)

	push := func(input chan *oproto.Value, timestamp uint64, value string) {
		input <- &oproto.Value{
			Timestamp:   timestamp,
			StringValue: value,
		}
	}

	test := func(output <-chan *oproto.Value, startTimestamp, endTimestamp uint64, testValue string) {
		value := <-output
		c.Assert(value.StringValue, Equals, testValue)
		c.Assert(value.Timestamp, Equals, startTimestamp)
		c.Assert(value.EndTimestamp, Equals, endTimestamp)
	}

	push(input, 1, "Hello")
	push(input, 2, "Hello")
	push(input, 3, "Hello")
	push(input, 4, "Goodbye")
	push(input, 5, "Friend")
	close(input)

	test(output, 1, 3, "Hello")
	test(output, 4, 4, "Goodbye")
	test(output, 5, 5, "Friend")
}

func (s *MySuite) BenchmarkEncodeDouble(c *C) {
	// Create pipeline
	for run := 0; run < c.N; run++ {
		input := make(chan *oproto.Value, 100000)
		output := Encode(input)

		// Number of repetitions of each value
		const blockSize = 100

		go func() {
			for i := 1; i <= 1000000; i++ {
				input <- &oproto.Value{
					Timestamp:   uint64(i),
					DoubleValue: float64(i - (i % blockSize)),
				}
			}
			close(input)
		}()

		for range output {
		}
	}
}
