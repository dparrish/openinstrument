package mutations

import (
	"testing"

	oproto "github.com/dparrish/openinstrument/proto"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func getAndCheck(c *C, output <-chan *oproto.Value, expectedValue float64) {
	value := <-output
	c.Check(value.DoubleValue, Equals, expectedValue)
}

func getAndCheckTimestamp(c *C, output <-chan *oproto.Value, expectedTimestamp uint64, expectedValue float64) {
	value := <-output
	c.Check(value.Timestamp, Equals, expectedTimestamp)
	c.Check(value.DoubleValue, Equals, expectedValue)
}

func writeValues(input chan *oproto.Value) {
	input <- &oproto.Value{Timestamp: uint64(0), DoubleValue: float64(20)}
	input <- &oproto.Value{Timestamp: uint64(60), DoubleValue: float64(40)}
	input <- &oproto.Value{Timestamp: uint64(120), DoubleValue: float64(60)}
	input <- &oproto.Value{Timestamp: uint64(180), DoubleValue: float64(80)}
	input <- &oproto.Value{Timestamp: uint64(240), DoubleValue: float64(100)}
	input <- &oproto.Value{Timestamp: uint64(300), DoubleValue: float64(120)}
	input <- &oproto.Value{Timestamp: uint64(330), DoubleValue: float64(130)}
	input <- &oproto.Value{Timestamp: uint64(360), DoubleValue: float64(140)}
	input <- &oproto.Value{Timestamp: uint64(420), DoubleValue: float64(160)}
	input <- &oproto.Value{Timestamp: uint64(480), DoubleValue: float64(180)}
	input <- &oproto.Value{Timestamp: uint64(540), DoubleValue: float64(200)}
	close(input)
}
func writeOffsetValues(input chan *oproto.Value) {
	input <- &oproto.Value{Timestamp: uint64(0), DoubleValue: float64(20)}
	input <- &oproto.Value{Timestamp: uint64(61), DoubleValue: float64(40)}
	input <- &oproto.Value{Timestamp: uint64(122), DoubleValue: float64(60)}
	input <- &oproto.Value{Timestamp: uint64(185), DoubleValue: float64(80)}
	input <- &oproto.Value{Timestamp: uint64(241), DoubleValue: float64(100)}
	input <- &oproto.Value{Timestamp: uint64(299), DoubleValue: float64(122)}
	input <- &oproto.Value{Timestamp: uint64(330), DoubleValue: float64(132)}
	input <- &oproto.Value{Timestamp: uint64(359), DoubleValue: float64(140)}
	input <- &oproto.Value{Timestamp: uint64(421), DoubleValue: float64(160)}
	input <- &oproto.Value{Timestamp: uint64(488), DoubleValue: float64(180)}
	input <- &oproto.Value{Timestamp: uint64(540), DoubleValue: float64(200)}
	input <- &oproto.Value{Timestamp: uint64(975), DoubleValue: float64(275)}
	close(input)
}

func (s *MySuite) TestMean(c *C) {
	input := make(chan *oproto.Value, 100)
	output := Mean(300, input)
	writeValues(input)
	getAndCheck(c, output, float64(70))
	getAndCheck(c, output, float64(162))
}

func (s *MySuite) TestInterpolate(c *C) {
	input := make(chan *oproto.Value, 100)
	output := Interpolate(300, input)
	writeOffsetValues(input)
	getAndCheckTimestamp(c, output, 0, float64(20))
	getAndCheckTimestamp(c, output, 300, float64(121.81818181818181))
	getAndCheckTimestamp(c, output, 600, float64(191.86046511627907))
	getAndCheckTimestamp(c, output, 900, float64(258.37209302325584))
}

func (s *MySuite) TestRate(c *C) {
	input := make(chan *oproto.Value, 100)
	output := Rate(300, input)
	writeValues(input)
	for v := range output {
		c.Check(v.DoubleValue, Equals, float64(1)/float64(3))
	}
}

func (s *MySuite) TestSignedRate(c *C) {
	input := make(chan *oproto.Value, 100)
	output := SignedRate(300, input)
	writeValues(input)
	for v := range output {
		c.Check(v.DoubleValue, Equals, float64(1)/float64(3))
	}
}
