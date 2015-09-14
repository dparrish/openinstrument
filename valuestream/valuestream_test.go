package valuestream

import (
	"testing"

	"code.google.com/p/goprotobuf/proto"

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
		stream.Value = append(stream.Value, &oproto.Value{DoubleValue: proto.Float64(float64(i))})
	}
	output := ToChan(stream)
	for i := 0; i < 10; i++ {
		v := <-output
		c.Check(v.GetDoubleValue(), Equals, float64(i))
	}
}

func (s *MySuite) TestFromChan(c *C) {
	input := make(chan *oproto.Value)
	output := &oproto.ValueStream{
		Value: make([]*oproto.Value, 0),
	}
	done := FromChan(input, output)
	for i := 0; i < 10; i++ {
		input <- &oproto.Value{DoubleValue: proto.Float64(float64(i))}
	}
	close(input)
	<-done
	for i := 0; i < 10; i++ {
		c.Check(output.Value[i].GetDoubleValue(), Equals, float64(i))
	}
}
