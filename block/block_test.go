package block

import (
	"fmt"
	"path/filepath"
	"testing"

	"golang.org/x/net/context"

	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/store_config"
	"github.com/dparrish/openinstrument/value"
	"github.com/dparrish/openinstrument/variable"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct {
	dataDir string
	config  store_config.ConfigStore
}

var _ = Suite(&MySuite{})

func (s *MySuite) SetUpSuite(c *C) {
	s.dataDir = c.MkDir()
	s.config = store_config.NewLocalConfigStore(filepath.Join(s.dataDir, "config.txt"), "test")
	s.config.Start(context.Background())
}

func (s *MySuite) TestRead(c *C) {
	block := NewBlock(context.Background(), "/test/foox", "", s.dataDir, s.config)
	streams := make(map[string]*oproto.ValueStream)
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("/test/foo%d", i)
		streams[key] = &oproto.ValueStream{
			Variable: &oproto.StreamVariable{Name: key},
			Value: []*oproto.Value{
				value.NewDouble(0, 0.0),
				value.NewDouble(0, 1.1),
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
	block := NewBlock(context.Background(), "/test/foo", "", s.dataDir, s.config)
	block.UnloggedStreams.Append(&oproto.ValueStream{
		Variable: &oproto.StreamVariable{Name: "/test/foo"},
		Value: []*oproto.Value{
			value.NewDouble(0, 0.0),
			value.NewDouble(0, 1.1),
		},
	})
	block.LoggedStreams.M["/test/bar"] = &oproto.ValueStream{
		Variable: &oproto.StreamVariable{Name: "/test/bar"},
		Value: []*oproto.Value{
			value.NewDouble(1, 1.1),
			value.NewDouble(2, 1.2),
			value.NewDouble(3, 1.3),
		},
	}
	block.LoggedStreams.M["/test/bar1"] = &oproto.ValueStream{
		Variable: &oproto.StreamVariable{Name: "/test/bar1"},
		Value: []*oproto.Value{
			value.NewDouble(4, 1.4),
			value.NewDouble(5, 1.5),
			value.NewDouble(6, 1.6),
		},
	}
	err := block.Write(context.Background(), block.LoggedStreams.M)
	c.Assert(err, Equals, nil)

	reader, err := block.GetIndexedStreams(context.Background())
	c.Assert(err, Equals, nil)
	vs := <-reader
	switch vs.Variable.Name {
	case "/test/bar":
		c.Assert(vs.Value[0].GetDouble(), Equals, 1.1)
		c.Assert(vs.Value[1].GetDouble(), Equals, 1.2)
		c.Assert(vs.Value[2].GetDouble(), Equals, 1.3)
	case "/test/bar1":
		c.Assert(vs.Value[0].GetDouble(), Equals, 1.4)
		c.Assert(vs.Value[1].GetDouble(), Equals, 1.5)
		c.Assert(vs.Value[2].GetDouble(), Equals, 1.6)
	default:
		c.Fail()
	}
}

func (s *MySuite) TestCompact(c *C) {
	block := NewBlock(context.Background(), "/test/foo", "", s.dataDir, s.config)
	for v := 0; v < 10; v++ {
		varName := fmt.Sprintf("/test/bar%d", v)
		block.LoggedStreams.M[varName] = &oproto.ValueStream{
			Variable: &oproto.StreamVariable{Name: varName},
			Value:    []*oproto.Value{},
		}
		for i := 0; i < 1000; i++ {
			block.LoggedStreams.M[varName].Value = append(block.LoggedStreams.M[varName].Value, value.NewDouble(0, float64(i)))
		}
	}
	c.Assert(block.Compact(context.Background()), IsNil)
}

func (s *MySuite) BenchmarkCompact(c *C) {
	block := NewBlock(context.Background(), "/test/foo", "", s.dataDir, s.config)
	for v := 0; v < 10; v++ {
		varName := fmt.Sprintf("/test/bar%d", v)
		block.LoggedStreams.M[varName] = &oproto.ValueStream{
			Variable: &oproto.StreamVariable{Name: varName},
			Value:    []*oproto.Value{},
		}
		for i := 0; i < 100000; i++ {
			block.LoggedStreams.M[varName].Value = append(block.LoggedStreams.M[varName].Value, value.NewDouble(0, float64(i)))
		}
	}
	for run := 0; run < c.N; run++ {
		block.Compact(context.Background())
	}
}

func (s *MySuite) TestGetStreamForVariable(c *C) {
	block := NewBlock(context.Background(), "/test/foo", "", s.dataDir, s.config)
	for v := 0; v < 10; v++ {
		varName := fmt.Sprintf("/test/bar%d", v)
		block.LoggedStreams.M[varName] = &oproto.ValueStream{
			Variable: &oproto.StreamVariable{Name: varName},
			Value:    []*oproto.Value{},
		}
		for i := 0; i < 100; i++ {
			block.LoggedStreams.M[varName].Value = append(block.LoggedStreams.M[varName].Value, value.NewDouble(0, float64(i)))
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

func (s *MySuite) TestGetOptimalSplitPoint(c *C) {
	block := NewBlock(context.Background(), "/test/foo{host=z}", "", s.dataDir, s.config)
	for i := 'a'; i <= 'z'; i++ {
		varName := fmt.Sprintf("/test/foo{host=%c}", i)
		v := []*oproto.Value{}
		for j := 'a'; j <= i; j++ {
			v = append(v, value.NewDouble(0, float64(i)))
		}
		block.LoggedStreams.M[varName] = &oproto.ValueStream{
			Variable: &oproto.StreamVariable{Name: varName},
			Value:    v,
		}
	}

	splitPoint, leftBlockEnd := block.GetOptimalSplitPoint(context.Background())
	c.Check(splitPoint, Equals, 18)
	c.Check(leftBlockEnd, Equals, "/test/foo{host=r}")
}
