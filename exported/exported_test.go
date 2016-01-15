package exported

import (
	"fmt"
	"testing"
	"time"

	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/variable"

	"golang.org/x/net/context"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct {
	dataDir string
	Client  *testClient
}

var _ = Suite(&MySuite{})

type testClient struct {
}

func (c *testClient) List(ctx context.Context, request *oproto.ListRequest) (<-chan *oproto.ListResponse, error) {
	return nil, nil
}

func (c *testClient) Get(ctx context.Context, request *oproto.GetRequest) (<-chan *oproto.GetResponse, error) {
	return nil, nil
}

func (c *testClient) GetConfig(ctx context.Context) (*oproto.ClusterConfig, error) {
	return nil, nil
}

func (c *testClient) Add(ctx context.Context) (chan<- *oproto.AddRequest, <-chan *oproto.AddResponse, error) {
	return nil, nil, nil
}

func (s *MySuite) SetUpSuite(c *C) {
	s.Client = &testClient{}
}

func (s *MySuite) TestInteger(c *C) {
	ctx, cancel := context.WithCancel(context.Background())
	e := NewExporter(ctx, s.Client, 10*time.Second)
	i := NewInteger(e, variable.NewFromString("/test/integer"))
	i.Set(100)

	// Force export
	streams := e.GetStreams()
	c.Assert(len(streams), Equals, 1)
	c.Check(variable.ProtoToString(streams[0].Variable), Equals, "/test/integer")
	c.Check(streams[0].Value[0].DoubleValue, Equals, 100.0)
	cancel()
	<-ctx.Done()
}

func (s *MySuite) TestFloat(c *C) {
	ctx, cancel := context.WithCancel(context.Background())
	e := NewExporter(ctx, s.Client, 10*time.Second)
	f := NewFloat(e, variable.NewFromString("/test/float"))
	f.Set(100.0)

	// Force export
	streams := e.GetStreams()
	c.Assert(len(streams), Equals, 1)
	c.Check(variable.ProtoToString(streams[0].Variable), Equals, "/test/float")
	c.Check(streams[0].Value[0].DoubleValue, Equals, 100.0)
	cancel()
	<-ctx.Done()
}

func (s *MySuite) TestRatio(c *C) {
	ctx, cancel := context.WithCancel(context.Background())
	e := NewExporter(ctx, s.Client, 10*time.Second)
	r := NewRatio(e, variable.NewFromString("/test/ratio"))

	for i := 0; i < 10; i++ {
		r.Success()
	}

	for i := 0; i < 5; i++ {
		r.Failure()
	}

	// Force export
	streams := e.GetStreams()
	c.Assert(len(streams), Equals, 3)
	for _, stream := range streams {
		switch variable.ProtoToString(stream.Variable) {
		case "/test/ratio-success":
			c.Check(stream.Value[0].DoubleValue, Equals, 10.0)
		case "/test/ratio-failure":
			c.Check(stream.Value[0].DoubleValue, Equals, 5.0)
		case "/test/ratio-total":
			c.Check(stream.Value[0].DoubleValue, Equals, 15.0)
		default:
			fmt.Printf("Invalid variable %s", variable.ProtoToString(stream.Variable))
			c.Fail()
		}
	}
	cancel()
	<-ctx.Done()
}

func (s *MySuite) TestAverage(c *C) {
	ctx, cancel := context.WithCancel(context.Background())
	e := NewExporter(ctx, s.Client, 10*time.Second)
	a := NewAverage(e, variable.NewFromString("/test/average"))

	a.Update(95, 1)
	a.Update(100, 1)
	a.Update(105, 1)

	// Force export
	streams := e.GetStreams()
	c.Assert(len(streams), Equals, 2)
	for _, stream := range streams {
		switch variable.ProtoToString(stream.Variable) {
		case "/test/average-total-count":
			c.Check(stream.Value[0].DoubleValue, Equals, 300.0)
		case "/test/average-overall-sum":
			c.Check(stream.Value[0].DoubleValue, Equals, 3.0)
		default:
			fmt.Printf("Invalid variable %s", variable.ProtoToString(stream.Variable))
			c.Fail()
		}
	}
	cancel()
	<-ctx.Done()
}

func (s *MySuite) TestTimer(c *C) {
	ctx, cancel := context.WithCancel(context.Background())
	e := NewExporter(ctx, s.Client, 10*time.Second)
	t := NewTimer(e, variable.NewFromString("/test/timer"))

	t.Start()
	time.Sleep(10 * time.Millisecond)
	t.Stop()

	t.Start()
	time.Sleep(10 * time.Millisecond)
	t.Stop()

	// Force export
	streams := e.GetStreams()
	c.Assert(len(streams), Equals, 2)
	for _, stream := range streams {
		switch variable.ProtoToString(stream.Variable) {
		case "/test/timer-total-count":
			c.Check((stream.Value[0].DoubleValue >= 20.0 && stream.Value[0].DoubleValue <= 25.0), Equals, true)
		case "/test/timer-overall-sum":
			c.Check(stream.Value[0].DoubleValue, Equals, 2.0)
		default:
			fmt.Printf("Invalid variable %s", variable.ProtoToString(stream.Variable))
			c.Fail()
		}
	}
	cancel()
	<-ctx.Done()
}

func (s *MySuite) TestString(c *C) {
	ctx, cancel := context.WithCancel(context.Background())
	e := NewExporter(ctx, s.Client, 10*time.Second)
	v := NewString(e, variable.NewFromString("/test/string"))
	v.Set("Testing")

	// Force export
	streams := e.GetStreams()
	c.Assert(len(streams), Equals, 1)
	c.Check(variable.ProtoToString(streams[0].Variable), Equals, "/test/string")
	c.Check(streams[0].Value[0].DoubleValue, Equals, 0.0)
	c.Check(streams[0].Value[0].StringValue, Equals, "Testing")
	cancel()
	<-ctx.Done()
}
