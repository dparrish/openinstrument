package valuestream

import (
	"fmt"
	"testing"

	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/value"
	"github.com/dparrish/openinstrument/variable"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestSort(c *C) {
}

func checkValueOrder(c <-chan *oproto.Value) (int, error) {
	var outCount int
	var lastTimestamp uint64
	for value := range c {
		outCount++
		if value.Timestamp < lastTimestamp {
			return 0, fmt.Errorf("Found unsorted value")
		}
		lastTimestamp = value.Timestamp
	}
	return outCount, nil
}

func (s *MySuite) TestMerge(c *C) {
	input := []*oproto.ValueStream{
		&oproto.ValueStream{
			Variable: variable.NewFromString("/test{host=a}").AsProto(),
			Value: []*oproto.Value{
				value.NewDouble(1, 1.0),
				value.NewDouble(4, 4.0),
			},
		},
		&oproto.ValueStream{
			Variable: variable.NewFromString("/test{host=b}").AsProto(),
			Value: []*oproto.Value{
				value.NewDouble(2, 2.0),
				value.NewDouble(5, 5.0),
			},
		},
		&oproto.ValueStream{
			Variable: variable.NewFromString("/test{host=c}").AsProto(),
			Value: []*oproto.Value{
				value.NewDouble(3, 3.0),
				value.NewDouble(6, 6.0),
			},
		},
	}
	outCount, err := checkValueOrder(Merge(input))
	if err != nil {
		c.Error(err)
	}
	c.Check(outCount, Equals, 6)
}

func (s *MySuite) TestMergeBy(c *C) {
	input := []*oproto.ValueStream{
		&oproto.ValueStream{
			Variable: variable.NewFromString("/test{host=a,other=x}").AsProto(),
			Value: []*oproto.Value{
				value.NewDouble(1, 1.0),
				value.NewDouble(4, 4.0),
			},
		},
		&oproto.ValueStream{
			Variable: variable.NewFromString("/test{host=b,other=y}").AsProto(),
			Value: []*oproto.Value{
				value.NewDouble(2, 2.0),
				value.NewDouble(5, 5.0),
			},
		},
		&oproto.ValueStream{
			Variable: variable.NewFromString("/test{host=a,other=z}").AsProto(),
			Value: []*oproto.Value{
				value.NewDouble(3, 3.0),
				value.NewDouble(6, 6.0),
			},
		},
	}
	numSets := 0
	for streams := range MergeBy(input, "host") {
		stv := variable.NewFromProto(streams[0].Variable)
		if stv.Match(variable.NewFromString("/test{host=a}")) {
			c.Assert(len(streams), Equals, 2)
			outCount, err := checkValueOrder(Merge(streams))
			if err != nil {
				c.Error(err)
			}
			c.Check(outCount, Equals, 4)
		} else {
			c.Check(len(streams), Equals, 1)
			outCount, err := checkValueOrder(Merge(streams))
			if err != nil {
				c.Error(err)
			}
			c.Check(outCount, Equals, 2)
		}
		numSets++
	}
	c.Check(numSets, Equals, 2)
}

func (s *MySuite) TestToChan(c *C) {
	stream := &oproto.ValueStream{
		Value: make([]*oproto.Value, 0),
	}
	for i := 0; i < 10; i++ {
		stream.Value = append(stream.Value, value.NewDouble(0, float64(i)))
	}
	output := ToChan(stream)
	for i := 0; i < 10; i++ {
		v := <-output
		c.Check(v.GetDouble(), Equals, float64(i))
	}
}
