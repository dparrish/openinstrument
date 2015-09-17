package datastore

import (
	"testing"

	oproto "github.com/dparrish/openinstrument/proto"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct {
	dataDir string
}

var _ = Suite(&MySuite{})

func (s *MySuite) SetUpSuite(c *C) {
	s.dataDir = c.MkDir()
}

func (s *MySuite) TestRead(c *C) {
}

func (s *MySuite) TestWrite(c *C) {
	block := newBlock("/test/foo", "")
	block.NewStreams = append(block.NewStreams, &oproto.ValueStream{
		Variable: &oproto.StreamVariable{Name: "/test/foo"},
		Value: []*oproto.Value{
			{DoubleValue: 0.0},
			{DoubleValue: 1.1},
		},
	})
	block.LogStreams["/test/bar"] = &oproto.ValueStream{
		Variable: &oproto.StreamVariable{Name: "/test/bar"},
		Value: []*oproto.Value{
			{Timestamp: 1, DoubleValue: 1.1},
			{Timestamp: 2, DoubleValue: 1.2},
			{Timestamp: 3, DoubleValue: 1.3},
		},
	}
	block.LogStreams["/test/bar1"] = &oproto.ValueStream{
		Variable: &oproto.StreamVariable{Name: "/test/bar1"},
		Value: []*oproto.Value{
			{Timestamp: 4, DoubleValue: 1.4},
			{Timestamp: 5, DoubleValue: 1.5},
			{Timestamp: 6, DoubleValue: 1.6},
		},
	}
	err := block.Write(s.dataDir, block.LogStreams)
	c.Assert(err, Equals, nil)

	reader, err := block.Read(s.dataDir)
	c.Assert(err, Equals, nil)
	vs := <-reader
	switch vs.Variable.Name {
	case "/test/bar":
		c.Check(vs.Value[0].DoubleValue, Equals, 1.1)
		c.Check(vs.Value[1].DoubleValue, Equals, 1.2)
		c.Check(vs.Value[2].DoubleValue, Equals, 1.3)
	case "/test/bar1":
		c.Check(vs.Value[0].DoubleValue, Equals, 1.4)
		c.Check(vs.Value[1].DoubleValue, Equals, 1.5)
		c.Check(vs.Value[2].DoubleValue, Equals, 1.6)
	default:
		c.Fail()
	}
}

func (s *MySuite) TestNumStreams(c *C) {
	block := newBlock("/test/foo", "")
	block.NewStreams = append(block.NewStreams, &oproto.ValueStream{
		Variable: &oproto.StreamVariable{Name: "/test/foo"},
		Value: []*oproto.Value{
			{DoubleValue: 0.0},
			{DoubleValue: 1.1},
		},
	})
	block.LogStreams["/test/bar"] = &oproto.ValueStream{
		Variable: &oproto.StreamVariable{Name: "/test/bar"},
		Value: []*oproto.Value{
			{DoubleValue: 1.1},
			{DoubleValue: 1.2},
			{DoubleValue: 1.3},
		},
	}
	c.Check(block.NumStreams(), Equals, uint32(2))
}

func (s *MySuite) TestNumValues(c *C) {
	block := newBlock("/test/foo", "")
	block.NewStreams = append(block.NewStreams, &oproto.ValueStream{
		Variable: &oproto.StreamVariable{Name: "/test/foo"},
		Value: []*oproto.Value{
			{DoubleValue: 0.0},
			{DoubleValue: 1.1},
		},
	})
	block.LogStreams["/test/bar"] = &oproto.ValueStream{
		Variable: &oproto.StreamVariable{Name: "/test/bar"},
		Value: []*oproto.Value{
			{DoubleValue: 1.2},
			{DoubleValue: 1.2},
			{DoubleValue: 1.2},
		},
	}
	c.Check(block.NumLogValues(), Equals, uint32(3))
	c.Check(block.NumValues(), Equals, uint32(5))
}
