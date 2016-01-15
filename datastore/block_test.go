package datastore

import (
	"fmt"
	"testing"

	"golang.org/x/net/context"

	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/variable"
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
	block := NewBlock(context.Background(), "/test/foox", "", s.dataDir)
	streams := make(map[string]*oproto.ValueStream)
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("/test/foo%d", i)
		streams[key] = &oproto.ValueStream{
			Variable: &oproto.StreamVariable{Name: key},
			Value: []*oproto.Value{
				{DoubleValue: 0.0},
				{DoubleValue: 1.1},
			},
		}
	}

	err := block.Write(context.Background(), streams)
	c.Assert(err, Equals, nil)

	reader, err := block.GetIndexedStreams(context.Background())
	c.Assert(err, Equals, nil)
	readStreams := 0
	for range reader {
		readStreams++
	}
	c.Assert(readStreams, Equals, len(streams))
}

func (s *MySuite) TestWrite(c *C) {
	block := NewBlock(context.Background(), "/test/foo", "", s.dataDir)
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
	err := block.Write(context.Background(), block.LogStreams)
	c.Assert(err, Equals, nil)

	reader, err := block.GetIndexedStreams(context.Background())
	c.Assert(err, Equals, nil)
	vs := <-reader
	switch vs.Variable.Name {
	case "/test/bar":
		c.Assert(vs.Value[0].DoubleValue, Equals, 1.1)
		c.Assert(vs.Value[1].DoubleValue, Equals, 1.2)
		c.Assert(vs.Value[2].DoubleValue, Equals, 1.3)
	case "/test/bar1":
		c.Assert(vs.Value[0].DoubleValue, Equals, 1.4)
		c.Assert(vs.Value[1].DoubleValue, Equals, 1.5)
		c.Assert(vs.Value[2].DoubleValue, Equals, 1.6)
	default:
		c.Fail()
	}
}

func (s *MySuite) TestCompact(c *C) {
	block := NewBlock(context.Background(), "/test/foo", "", s.dataDir)
	for v := 0; v < 10; v++ {
		varName := fmt.Sprintf("/test/bar%d", v)
		block.LogStreams[varName] = &oproto.ValueStream{
			Variable: &oproto.StreamVariable{Name: varName},
			Value:    []*oproto.Value{},
		}
		for i := 0; i < 1000; i++ {
			block.LogStreams[varName].Value = append(block.LogStreams[varName].Value,
				&oproto.Value{DoubleValue: float64(i)})
		}
	}
	c.Assert(block.Compact(context.Background()), IsNil)
}

func (s *MySuite) TestGetStreamForVariable(c *C) {
	block := NewBlock(context.Background(), "/test/foo", "", s.dataDir)
	for v := 0; v < 10; v++ {
		varName := fmt.Sprintf("/test/bar%d", v)
		block.LogStreams[varName] = &oproto.ValueStream{
			Variable: &oproto.StreamVariable{Name: varName},
			Value:    []*oproto.Value{},
		}
		for i := 0; i < 100; i++ {
			block.LogStreams[varName].Value = append(block.LogStreams[varName].Value,
				&oproto.Value{DoubleValue: float64(i)})
		}
	}
	c.Assert(block.Compact(context.Background()), IsNil)
	found := false
	for _, index := range block.Block.Header.Index {
		cv := variable.NewFromProto(index.Variable)
		if cv.String() != "/test/bar7" {
			continue
		}
		stream := block.getIndexedStream(context.Background(), index)
		c.Assert(variable.ProtoToString(stream.Variable), Equals, "/test/bar7")
		found = true
		break
	}
	c.Assert(found, Equals, true)
}

func (s *MySuite) BenchmarkCompact(c *C) {
	block := NewBlock(context.Background(), "/test/foo", "", s.dataDir)
	for v := 0; v < 10; v++ {
		varName := fmt.Sprintf("/test/bar%d", v)
		block.LogStreams[varName] = &oproto.ValueStream{
			Variable: &oproto.StreamVariable{Name: varName},
			Value:    []*oproto.Value{},
		}
		for i := 0; i < 100000; i++ {
			block.LogStreams[varName].Value = append(block.LogStreams[varName].Value,
				&oproto.Value{DoubleValue: float64(i)})
		}
	}
	for run := 0; run < c.N; run++ {
		block.Compact(context.Background())
	}
}
