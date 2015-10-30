package valuestream

import (
	"testing"

	oproto "github.com/dparrish/openinstrument/proto"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestSort(c *C) {
}

func (s *MySuite) TestMerge(c *C) {
}

func (s *MySuite) TestMergeBy(c *C) {
}

func (s *MySuite) TestToChan(c *C) {
	stream := &oproto.ValueStream{
		Value: make([]*oproto.Value, 0),
	}
	for i := 0; i < 10; i++ {
		stream.Value = append(stream.Value, &oproto.Value{DoubleValue: float64(i)})
	}
	output := ToChan(stream)
	for i := 0; i < 10; i++ {
		v := <-output
		c.Check(v.DoubleValue, Equals, float64(i))
	}
}
