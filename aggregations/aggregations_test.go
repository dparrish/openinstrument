package aggregations_test

// Different package name to resolve cyclic dependencies

import (
	"testing"

	"golang.org/x/net/context"

	. "github.com/dparrish/openinstrument/aggregations"
	"github.com/dparrish/openinstrument/datastore"
	"github.com/dparrish/openinstrument/value"
	"github.com/dparrish/openinstrument/variable"

	oproto "github.com/dparrish/openinstrument/proto"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct {
	store datastore.ReadableStore
}

var _ = Suite(&MySuite{})

func (s *MySuite) SetUpSuite(c *C) {
	s.store = &FakeReadableStore{}
}

func (s *MySuite) TestCartesianProduct(c *C) {
	input := [][]string{
		{"host=a", "host=b"},
		{"job=a", "job=b"},
		{"other=x", "other=y"},
	}
	output := CartesianProduct(input)
	c.Assert(len(output), Equals, 2*2*2)
	c.Check(output[0], DeepEquals, []string{"host=a", "job=a", "other=x"})
	c.Check(output[len(output)-1], DeepEquals, []string{"host=b", "job=b", "other=y"})
}

func (s *MySuite) TestCartesianProductToVariable(c *C) {
	input := [][]string{
		{"host=a", "host=b"},
		{"job=a", "job=b"},
		{"other=x", "other=y"},
	}
	output := CartesianProductToVariable(CartesianProduct(input))
	c.Assert(len(output), Equals, 2*2*2)
	c.Check(output[0].String(), Equals, "{host=a,job=a,other=x}")
	c.Check(output[len(output)-1].String(), Equals, "{host=b,job=b,other=y}")
}

func checkValue(c *C, value *oproto.Value, expectedTimestamp uint64, expectedValue float64) {
	c.Assert(value, Not(IsNil))
	c.Check(value.GetDouble(), Equals, expectedValue)
	c.Check(int(value.Timestamp), Equals, int(expectedTimestamp))
}

func (s *MySuite) TestMean(c *C) {
	input := []*oproto.ValueStream{
		<-s.store.Reader(context.Background(), variable.NewFromString("/test{host=a}")),
		<-s.store.Reader(context.Background(), variable.NewFromString("/test{host=b}")),
	}
	output := Mean(nil, input)
	c.Assert(output, Not(IsNil))
	c.Assert(len(output), Equals, 1)
	stream := output[0]
	c.Assert(len(stream.Value), Equals, 11)

	checkValue(c, stream.Value[0], 60*0, float64((20*1+40*1)/2))
	checkValue(c, stream.Value[1], 60*1, float64((20*2+40*2)/2))
	checkValue(c, stream.Value[2], 60*2, float64((20*3+40*3)/2))
	checkValue(c, stream.Value[3], 60*3, float64((20*4+40*4)/2))
	checkValue(c, stream.Value[4], 60*4, float64((20*5+40*5)/2))
	checkValue(c, stream.Value[5], 60*5, float64((20*6+40*6)/2))
	checkValue(c, stream.Value[6], 60*6, float64((20*7+40*7)/2))
	checkValue(c, stream.Value[7], 60*7, float64((20*8+40*8)/2))
	checkValue(c, stream.Value[8], 60*8, float64((20*9+40*9)/2))
	checkValue(c, stream.Value[9], 60*9, float64((20*10+40*10)/2))
	checkValue(c, stream.Value[10], 60*10, float64((20*11+40*11)/2))
}

func (s *MySuite) TestMeanBy(c *C) {
	input := []*oproto.ValueStream{
		<-s.store.Reader(context.Background(), variable.NewFromString("/test{host=a,job=foo,other=w}")),
		<-s.store.Reader(context.Background(), variable.NewFromString("/test{host=a,job=bar,other=x}")),
		<-s.store.Reader(context.Background(), variable.NewFromString("/test{host=b,job=foo,other=y}")),
		<-s.store.Reader(context.Background(), variable.NewFromString("/test{host=b,job=bar,other=z}")),
	}
	output := Mean([]string{"host"}, input)
	c.Assert(output, Not(IsNil))
	c.Assert(len(output), Equals, 2)
	// Check that there are two output streams, as there are two distinct hosts
	checkHostsAndJobs(c, output, 1, 1, 0, 0)

	for _, stream := range output {
		c.Assert(len(stream.Value), Equals, 11)
		if stream.Variable.Label["host"] == "a" {
			checkValue(c, stream.Value[0], 60*0, float64((20*1+20*1)/2))
			checkValue(c, stream.Value[1], 60*1, float64((20*2+20*2)/2))
			checkValue(c, stream.Value[2], 60*2, float64((20*3+20*3)/2))
			checkValue(c, stream.Value[3], 60*3, float64((20*4+20*4)/2))
			checkValue(c, stream.Value[4], 60*4, float64((20*5+20*5)/2))
			checkValue(c, stream.Value[5], 60*5, float64((20*6+20*6)/2))
			checkValue(c, stream.Value[6], 60*6, float64((20*7+20*7)/2))
			checkValue(c, stream.Value[7], 60*7, float64((20*8+20*8)/2))
			checkValue(c, stream.Value[8], 60*8, float64((20*9+20*9)/2))
			checkValue(c, stream.Value[9], 60*9, float64((20*10+20*10)/2))
			checkValue(c, stream.Value[10], 60*10, float64((20*11+20*11)/2))
		} else if stream.Variable.Label["host"] == "b" {
			checkValue(c, stream.Value[0], 60*0, float64((40*1+40*1)/2))
			checkValue(c, stream.Value[1], 60*1, float64((40*2+40*2)/2))
			checkValue(c, stream.Value[2], 60*2, float64((40*3+40*3)/2))
			checkValue(c, stream.Value[3], 60*3, float64((40*4+40*4)/2))
			checkValue(c, stream.Value[4], 60*4, float64((40*5+40*5)/2))
			checkValue(c, stream.Value[5], 60*5, float64((40*6+40*6)/2))
			checkValue(c, stream.Value[6], 60*6, float64((40*7+40*7)/2))
			checkValue(c, stream.Value[7], 60*7, float64((40*8+40*8)/2))
			checkValue(c, stream.Value[8], 60*8, float64((40*9+40*9)/2))
			checkValue(c, stream.Value[9], 60*9, float64((40*10+40*10)/2))
			checkValue(c, stream.Value[10], 60*10, float64((40*11+40*11)/2))
		} else {
			c.Fail()
		}
	}
}

func (s *MySuite) TestMeanByJob(c *C) {
	input := []*oproto.ValueStream{
		<-s.store.Reader(context.Background(), variable.NewFromString("/test{host=a,job=foo,other=w}")),
		<-s.store.Reader(context.Background(), variable.NewFromString("/test{host=a,job=bar,other=x}")),
		<-s.store.Reader(context.Background(), variable.NewFromString("/test{host=b,job=foo,other=y}")),
		<-s.store.Reader(context.Background(), variable.NewFromString("/test{host=b,job=bar,other=z}")),
	}
	output := Mean([]string{"job"}, input)
	c.Assert(output, Not(IsNil))
	c.Assert(len(output), Equals, 2)
	// Check that there are two output streams, as there are two distinct hosts
	checkHostsAndJobs(c, output, 0, 0, 1, 1)

	// First stream is job=foo
	stream := output[0]
	c.Assert(len(stream.Value), Equals, 11)
	checkValue(c, stream.Value[0], 60*0, float64((20*1+40*1)/2))
	checkValue(c, stream.Value[1], 60*1, float64((20*2+40*2)/2))
	checkValue(c, stream.Value[2], 60*2, float64((20*3+40*3)/2))
	checkValue(c, stream.Value[3], 60*3, float64((20*4+40*4)/2))
	checkValue(c, stream.Value[4], 60*4, float64((20*5+40*5)/2))
	checkValue(c, stream.Value[5], 60*5, float64((20*6+40*6)/2))
	checkValue(c, stream.Value[6], 60*6, float64((20*7+40*7)/2))
	checkValue(c, stream.Value[7], 60*7, float64((20*8+40*8)/2))
	checkValue(c, stream.Value[8], 60*8, float64((20*9+40*9)/2))
	checkValue(c, stream.Value[9], 60*9, float64((20*10+40*10)/2))
	checkValue(c, stream.Value[10], 60*10, float64((20*11+40*11)/2))

	// Second stream is job=bar
	stream = output[1]
	c.Assert(len(stream.Value), Equals, 11)
	checkValue(c, stream.Value[0], 60*0, float64((20*1+40*1)/2))
	checkValue(c, stream.Value[1], 60*1, float64((20*2+40*2)/2))
	checkValue(c, stream.Value[2], 60*2, float64((20*3+40*3)/2))
	checkValue(c, stream.Value[3], 60*3, float64((20*4+40*4)/2))
	checkValue(c, stream.Value[4], 60*4, float64((20*5+40*5)/2))
	checkValue(c, stream.Value[5], 60*5, float64((20*6+40*6)/2))
	checkValue(c, stream.Value[6], 60*6, float64((20*7+40*7)/2))
	checkValue(c, stream.Value[7], 60*7, float64((20*8+40*8)/2))
	checkValue(c, stream.Value[8], 60*8, float64((20*9+40*9)/2))
	checkValue(c, stream.Value[9], 60*9, float64((20*10+40*10)/2))
	checkValue(c, stream.Value[10], 60*10, float64((20*11+40*11)/2))
}

func checkHostsAndJobs(c *C, streams []*oproto.ValueStream, hostA, hostB, jobFoo, jobBar int) {
	job := map[string]int{"foo": 0, "bar": 0}
	host := map[string]int{"a": 0, "b": 0}
	for _, stream := range streams {
		job[stream.Variable.Label["job"]]++
		host[stream.Variable.Label["host"]]++
	}
	c.Check(host["a"], Equals, hostA)
	c.Check(host["b"], Equals, hostB)
	c.Check(job["foo"], Equals, jobFoo)
	c.Check(job["bar"], Equals, jobBar)
}

func (s *MySuite) TestMeanByTwoLabels(c *C) {
	input := []*oproto.ValueStream{
		<-s.store.Reader(context.Background(), variable.NewFromString("/test{host=a,job=foo,other=w}")),
		<-s.store.Reader(context.Background(), variable.NewFromString("/test{host=a,job=bar,other=x}")),
		<-s.store.Reader(context.Background(), variable.NewFromString("/test{host=b,job=foo,other=y}")),
		<-s.store.Reader(context.Background(), variable.NewFromString("/test{host=b,job=bar,other=z}")),
	}
	output := Mean([]string{"host", "job"}, input)
	c.Assert(output, Not(IsNil))
	c.Assert(len(output), Equals, 4)
	// Check that there are 4 output streams with the correct number of output labels
	checkHostsAndJobs(c, output, 2, 2, 2, 2)
}

func (s *MySuite) TestOffsetMean(c *C) {
	return
	input := []*oproto.ValueStream{
		<-s.store.Reader(context.Background(), variable.NewFromString("/test/offset{host=a}")),
		<-s.store.Reader(context.Background(), variable.NewFromString("/test/offset{host=b}")),
	}
	output := Mean(nil, input)
	c.Assert(output, Not(IsNil))
	c.Assert(len(output), Equals, 1)
	stream := output[0]
	c.Assert(len(stream.Value), Equals, 11)

	checkValue(c, stream.Value[0], 60*0, float64((20*1+40*1)/2))
	checkValue(c, stream.Value[1], 60*1, float64((20*2+40*2)/2))
	checkValue(c, stream.Value[2], 60*2, float64((20*3+40*3)/2))
	checkValue(c, stream.Value[3], 60*3, float64((20*4+40*4)/2))
	checkValue(c, stream.Value[4], 60*4, float64((20*5+40*5)/2))
	checkValue(c, stream.Value[5], 60*5, float64((20*6+40*6)/2))
	checkValue(c, stream.Value[6], 60*6, float64((20*7+40*7)/2))
	checkValue(c, stream.Value[7], 60*7, float64((20*8+40*8)/2))
	checkValue(c, stream.Value[8], 60*8, float64((20*9+40*9)/2))
	checkValue(c, stream.Value[9], 60*9, float64((20*10+40*10)/2))
	checkValue(c, stream.Value[10], 60*10, float64((20*11+40*11)/2))
}

func (s *MySuite) TestMin(c *C) {
	input := []*oproto.ValueStream{
		<-s.store.Reader(context.Background(), variable.NewFromString("/test{host=a}")),
		<-s.store.Reader(context.Background(), variable.NewFromString("/test{host=b}")),
	}
	output := Min(nil, input)
	c.Assert(output, Not(IsNil))
	c.Assert(len(output), Equals, 1)
	stream := output[0]
	c.Assert(len(stream.Value), Equals, 11)

	checkValue(c, stream.Value[0], 60*0, float64(20*1))
	checkValue(c, stream.Value[1], 60*1, float64(20*2))
	checkValue(c, stream.Value[2], 60*2, float64(20*3))
	checkValue(c, stream.Value[3], 60*3, float64(20*4))
	checkValue(c, stream.Value[4], 60*4, float64(20*5))
	checkValue(c, stream.Value[5], 60*5, float64(20*6))
	checkValue(c, stream.Value[6], 60*6, float64(20*7))
	checkValue(c, stream.Value[7], 60*7, float64(20*8))
	checkValue(c, stream.Value[8], 60*8, float64(20*9))
	checkValue(c, stream.Value[9], 60*9, float64(20*10))
	checkValue(c, stream.Value[10], 60*10, float64(20*11))
}

func (s *MySuite) TestMax(c *C) {
	input := []*oproto.ValueStream{
		<-s.store.Reader(context.Background(), variable.NewFromString("/test{host=a}")),
		<-s.store.Reader(context.Background(), variable.NewFromString("/test{host=b}")),
	}
	output := Max(nil, input)
	c.Assert(output, Not(IsNil))
	c.Assert(len(output), Equals, 1)
	stream := output[0]
	c.Assert(len(stream.Value), Equals, 11)

	checkValue(c, stream.Value[0], 60*0, float64(40*1))
	checkValue(c, stream.Value[1], 60*1, float64(40*2))
	checkValue(c, stream.Value[2], 60*2, float64(40*3))
	checkValue(c, stream.Value[3], 60*3, float64(40*4))
	checkValue(c, stream.Value[4], 60*4, float64(40*5))
	checkValue(c, stream.Value[5], 60*5, float64(40*6))
	checkValue(c, stream.Value[6], 60*6, float64(40*7))
	checkValue(c, stream.Value[7], 60*7, float64(40*8))
	checkValue(c, stream.Value[8], 60*8, float64(40*9))
	checkValue(c, stream.Value[9], 60*9, float64(40*10))
	checkValue(c, stream.Value[10], 60*10, float64(40*11))
}

func (s *MySuite) TestSum(c *C) {
	input := []*oproto.ValueStream{
		<-s.store.Reader(context.Background(), variable.NewFromString("/test{host=a}")),
		<-s.store.Reader(context.Background(), variable.NewFromString("/test{host=b}")),
	}
	output := Sum(nil, input)
	c.Assert(output, Not(IsNil))
	c.Assert(len(output), Equals, 1)
	stream := output[0]
	c.Assert(len(stream.Value), Equals, 11)

	checkValue(c, stream.Value[0], 60*0, float64((20*1 + 40*1)))
	checkValue(c, stream.Value[1], 60*1, float64((20*2 + 40*2)))
	checkValue(c, stream.Value[2], 60*2, float64((20*3 + 40*3)))
	checkValue(c, stream.Value[3], 60*3, float64((20*4 + 40*4)))
	checkValue(c, stream.Value[4], 60*4, float64((20*5 + 40*5)))
	checkValue(c, stream.Value[5], 60*5, float64((20*6 + 40*6)))
	checkValue(c, stream.Value[6], 60*6, float64((20*7 + 40*7)))
	checkValue(c, stream.Value[7], 60*7, float64((20*8 + 40*8)))
	checkValue(c, stream.Value[8], 60*8, float64((20*9 + 40*9)))
	checkValue(c, stream.Value[9], 60*9, float64((20*10 + 40*10)))
	checkValue(c, stream.Value[10], 60*10, float64((20*11 + 40*11)))
}

func (s *MySuite) TestStdDev(c *C) {
	input := []*oproto.ValueStream{
		<-s.store.Reader(context.Background(), variable.NewFromString("/test{host=a}")),
		<-s.store.Reader(context.Background(), variable.NewFromString("/test{host=b}")),
	}
	output := StdDev(nil, input)
	c.Assert(output, Not(IsNil))
	c.Assert(len(output), Equals, 1)
	stream := output[0]
	c.Assert(len(stream.Value), Equals, 11)

	checkValue(c, stream.Value[0], 60*0, 10.0)
	checkValue(c, stream.Value[1], 60*1, 20.0)
	checkValue(c, stream.Value[2], 60*2, 30.0)
	checkValue(c, stream.Value[3], 60*3, 40.0)
	checkValue(c, stream.Value[4], 60*4, 50.0)
	checkValue(c, stream.Value[5], 60*5, 60.0)
	checkValue(c, stream.Value[6], 60*6, 70.0)
	checkValue(c, stream.Value[7], 60*7, 80.0)
	checkValue(c, stream.Value[8], 60*8, 90.0)
	checkValue(c, stream.Value[9], 60*9, 100.0)
	checkValue(c, stream.Value[10], 60*10, 110.0)
}

type FakeReadableStore struct{}

func (s *FakeReadableStore) Reader(ctx context.Context, v *variable.Variable) <-chan *oproto.ValueStream {
	c := make(chan *oproto.ValueStream, 100)
	go func() {
		defer close(c)
		var stream *oproto.ValueStream
		switch v.Variable {
		case "/test/offset":
			stream = &oproto.ValueStream{
				Variable: v.AsProto(),
				Value: []*oproto.Value{
					value.NewDouble(uint64(60*0), float64(20*1)),
					value.NewDouble(uint64(60*1), float64(20*2)),
					value.NewDouble(uint64(60*2), float64(20*3)),
					value.NewDouble(uint64(60*3), float64(20*4)),
					value.NewDouble(uint64(60*4), float64(20*5)),
					value.NewDouble(uint64(60*5), float64(20*6)),
					value.NewDouble(uint64(60*6), float64(20*7)),
					value.NewDouble(uint64(60*7), float64(20*8)),
					value.NewDouble(uint64(60*8), float64(20*9)),
					value.NewDouble(uint64(60*9), float64(20*10)),
					value.NewDouble(uint64(60*10), float64(20*11)),
				},
			}
			if v.Labels["host"] == "a" {
				for _, value := range stream.Value {
					value.Timestamp -= 3
				}
			}
			if v.Labels["host"] == "b" {
				for _, value := range stream.Value {
					value.Timestamp += 5
				}
			}
		case "/test":
			stream = &oproto.ValueStream{
				Variable: v.AsProto(),
				Value: []*oproto.Value{
					value.NewDouble(uint64(60*0), float64(20*1)),
					value.NewDouble(uint64(60*1), float64(20*2)),
					value.NewDouble(uint64(60*2), float64(20*3)),
					value.NewDouble(uint64(60*3), float64(20*4)),
					value.NewDouble(uint64(60*4), float64(20*5)),
					value.NewDouble(uint64(60*5), float64(20*6)),
					value.NewDouble(uint64(60*6), float64(20*7)),
					value.NewDouble(uint64(60*7), float64(20*8)),
					value.NewDouble(uint64(60*8), float64(20*9)),
					value.NewDouble(uint64(60*9), float64(20*10)),
					value.NewDouble(uint64(60*10), float64(20*11)),
				},
			}
		}
		if v.Labels["host"] == "b" {
			for _, value := range stream.Value {
				value.Value.(*oproto.Value_Double).Double *= 2
			}
		}
		c <- stream
	}()
	return c
}
