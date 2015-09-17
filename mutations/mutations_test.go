package mutations

import (
	"testing"

	"github.com/golang/protobuf/proto"
	oproto "github.com/dparrish/openinstrument/proto"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func getAndCheck(c *C, output <-chan *oproto.Value, expectedValue float64) {
	value := <-output
	c.Check(value.GetDoubleValue(), Equals, expectedValue)
}

func getAndCheckTimestamp(c *C, output <-chan *oproto.Value, expectedTimestamp uint64, expectedValue float64) {
	value := <-output
	c.Check(value.GetTimestamp(), Equals, expectedTimestamp)
	c.Check(value.GetDoubleValue(), Equals, expectedValue)
}

func writeValues(input chan *oproto.Value) {
	input <- &oproto.Value{Timestamp: proto.Uint64(0), DoubleValue: proto.Float64(20)}
	input <- &oproto.Value{Timestamp: proto.Uint64(60), DoubleValue: proto.Float64(40)}
	input <- &oproto.Value{Timestamp: proto.Uint64(120), DoubleValue: proto.Float64(60)}
	input <- &oproto.Value{Timestamp: proto.Uint64(180), DoubleValue: proto.Float64(80)}
	input <- &oproto.Value{Timestamp: proto.Uint64(240), DoubleValue: proto.Float64(100)}
	input <- &oproto.Value{Timestamp: proto.Uint64(300), DoubleValue: proto.Float64(120)}
	input <- &oproto.Value{Timestamp: proto.Uint64(330), DoubleValue: proto.Float64(130)}
	input <- &oproto.Value{Timestamp: proto.Uint64(360), DoubleValue: proto.Float64(140)}
	input <- &oproto.Value{Timestamp: proto.Uint64(420), DoubleValue: proto.Float64(160)}
	input <- &oproto.Value{Timestamp: proto.Uint64(480), DoubleValue: proto.Float64(180)}
	input <- &oproto.Value{Timestamp: proto.Uint64(540), DoubleValue: proto.Float64(200)}
	close(input)
}
func writeOffsetValues(input chan *oproto.Value) {
	input <- &oproto.Value{Timestamp: proto.Uint64(0), DoubleValue: proto.Float64(20)}
	input <- &oproto.Value{Timestamp: proto.Uint64(61), DoubleValue: proto.Float64(40)}
	input <- &oproto.Value{Timestamp: proto.Uint64(122), DoubleValue: proto.Float64(60)}
	input <- &oproto.Value{Timestamp: proto.Uint64(185), DoubleValue: proto.Float64(80)}
	input <- &oproto.Value{Timestamp: proto.Uint64(241), DoubleValue: proto.Float64(100)}
	input <- &oproto.Value{Timestamp: proto.Uint64(299), DoubleValue: proto.Float64(122)}
	input <- &oproto.Value{Timestamp: proto.Uint64(330), DoubleValue: proto.Float64(132)}
	input <- &oproto.Value{Timestamp: proto.Uint64(359), DoubleValue: proto.Float64(140)}
	input <- &oproto.Value{Timestamp: proto.Uint64(421), DoubleValue: proto.Float64(160)}
	input <- &oproto.Value{Timestamp: proto.Uint64(488), DoubleValue: proto.Float64(180)}
	input <- &oproto.Value{Timestamp: proto.Uint64(540), DoubleValue: proto.Float64(200)}
	input <- &oproto.Value{Timestamp: proto.Uint64(975), DoubleValue: proto.Float64(275)}
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
		c.Check(v.GetDoubleValue(), Equals, float64(1)/float64(3))
	}
}

func (s *MySuite) TestSignedRate(c *C) {
	input := make(chan *oproto.Value, 100)
	output := SignedRate(300, input)
	writeValues(input)
	for v := range output {
		c.Check(v.GetDoubleValue(), Equals, float64(1)/float64(3))
	}
}
