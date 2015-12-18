package query

import (
	"testing"

	"golang.org/x/net/context"

	"github.com/dparrish/openinstrument/datastore"
	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/variable"
	"github.com/golang/protobuf/proto"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

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
					{Timestamp: uint64(60 * 0), DoubleValue: float64(20 * 1)},
					{Timestamp: uint64(60 * 1), DoubleValue: float64(20 * 2)},
					{Timestamp: uint64(60 * 2), DoubleValue: float64(20 * 3)},
					{Timestamp: uint64(60 * 3), DoubleValue: float64(20 * 4)},
					{Timestamp: uint64(60 * 4), DoubleValue: float64(20 * 5)},
					{Timestamp: uint64(60 * 5), DoubleValue: float64(20 * 6)},
					{Timestamp: uint64(60 * 6), DoubleValue: float64(20 * 7)},
					{Timestamp: uint64(60 * 7), DoubleValue: float64(20 * 8)},
					{Timestamp: uint64(60 * 8), DoubleValue: float64(20 * 9)},
					{Timestamp: uint64(60 * 9), DoubleValue: float64(20 * 10)},
					{Timestamp: uint64(60 * 10), DoubleValue: float64(20 * 11)},
				},
			}
			if stream != nil && v.Labels["host"] == "a" {
				for _, value := range stream.Value {
					value.Timestamp -= 3
				}
			}
			if stream != nil && v.Labels["host"] == "b" {
				for _, value := range stream.Value {
					value.Timestamp += 5
				}
			}
		case "/test":
			stream = &oproto.ValueStream{
				Variable: v.AsProto(),
				Value: []*oproto.Value{
					{Timestamp: uint64(60 * 0), DoubleValue: float64(20 * 1)},
					{Timestamp: uint64(60 * 1), DoubleValue: float64(20 * 2)},
					{Timestamp: uint64(60 * 2), DoubleValue: float64(20 * 3)},
					{Timestamp: uint64(60 * 3), DoubleValue: float64(20 * 4)},
					{Timestamp: uint64(60 * 4), DoubleValue: float64(20 * 5)},
					{Timestamp: uint64(60 * 5), DoubleValue: float64(20 * 6)},
					{Timestamp: uint64(60 * 6), DoubleValue: float64(20 * 7)},
					{Timestamp: uint64(60 * 7), DoubleValue: float64(20 * 8)},
					{Timestamp: uint64(60 * 8), DoubleValue: float64(20 * 9)},
					{Timestamp: uint64(60 * 9), DoubleValue: float64(20 * 10)},
					{Timestamp: uint64(60 * 10), DoubleValue: float64(20 * 11)},
				},
			}
		}
		if stream != nil && v.Labels["host"] == "b" {
			for _, value := range stream.Value {
				value.DoubleValue *= 2
			}
		}
		c <- stream
	}()
	return c
}

type MySuite struct {
	store datastore.ReadableStore
}

func (s *MySuite) SetUpSuite(c *C) {
	s.store = &FakeReadableStore{}
}

var _ = Suite(&MySuite{})

func (s *MySuite) TestVariableNoLabelsOrBraces(c *C) {
	q, err := Parse("/test")
	c.Assert(err, IsNil)
	query := q.query
	c.Check(variable.ProtoToString(query.Variable[0]), Equals, "/test")
}

func (s *MySuite) TestVariableNoLabels(c *C) {
	q, err := Parse("/test{}")
	c.Assert(err, IsNil)
	query := q.query
	c.Check(variable.ProtoToString(query.Variable[0]), Equals, "/test")
	c.Check(query.Variable[0].MinTimestamp, Equals, int64(0))
	c.Check(query.Variable[0].MaxTimestamp, Equals, int64(0))

	ch, err := q.Run(context.Background(), s.store)
	c.Assert(err, IsNil)
	numStreams := 0
	for stream := range ch {
		c.Check(variable.ProtoToString(stream.Variable), Equals, "/test")
		c.Check(len(stream.Value), Equals, 11)
		numStreams++
	}
	c.Check(numStreams, Equals, 1)
}

func (s *MySuite) TestVariableWithStartRange(c *C) {
	q, err := Parse("/test[200]")
	c.Assert(err, IsNil)
	query := q.query
	c.Check(variable.ProtoToString(query.Variable[0]), Equals, "/test[200]")
	c.Check(query.Variable[0].MinTimestamp, Equals, int64(200))
	c.Check(query.Variable[0].MaxTimestamp, Equals, int64(0))

	ch, err := q.Run(context.Background(), s.store)
	c.Assert(err, IsNil)
	numStreams := 0
	for stream := range ch {
		c.Check(variable.ProtoToString(stream.Variable), Equals, "/test[200]")
		c.Check(len(stream.Value), Equals, 7)
		numStreams++
	}
	c.Check(numStreams, Equals, 1)
}

func (s *MySuite) TestVariableWithEndRange(c *C) {
	q, err := Parse("/test{host=a}[200:400]")
	c.Assert(err, IsNil)
	query := q.query
	c.Check(variable.ProtoToString(query.Variable[0]), Equals, "/test{host=a}[200:400]")
	c.Check(query.Variable[0].MinTimestamp, Equals, int64(200))
	c.Check(query.Variable[0].MaxTimestamp, Equals, int64(400))

	ch, err := q.Run(context.Background(), s.store)
	c.Assert(err, IsNil)
	numStreams := 0
	for stream := range ch {
		c.Check(variable.ProtoToString(stream.Variable), Equals, "/test{host=a}[200:400]")
		c.Check(len(stream.Value), Equals, 3)
		numStreams++
	}
	c.Check(numStreams, Equals, 1)
}

func (s *MySuite) TestVariableOneLabel(c *C) {
	q, err := Parse("/test{host=a}")
	c.Assert(err, IsNil)
	query := q.query
	c.Check(variable.ProtoToString(query.Variable[0]), Equals, "/test{host=a}")
}

func (s *MySuite) TestVariableTwoLabels(c *C) {
	q, err := Parse("/test{x=y,host=a}")
	c.Assert(err, IsNil)
	query := q.query
	c.Check(variable.ProtoToString(query.Variable[0]), Equals, "/test{host=a,x=y}")
}

func (s *MySuite) TestLabelWildcard(c *C) {
	q, err := Parse("/test{host=*}")
	c.Assert(err, IsNil)
	query := q.query
	c.Check(variable.ProtoToString(query.Variable[0]), Equals, "/test{host=*}")
}

func (s *MySuite) TestMean(c *C) {
	q, err := Parse("mean by (xyz) (/test{host=a}, /test{host=b})")
	c.Assert(err, IsNil)
	query := q.query
	c.Check(variable.ProtoToString(query.Aggregation[0].Query[0].Variable[0]), Equals, "/test{host=a}")
	c.Check(variable.ProtoToString(query.Aggregation[0].Query[0].Variable[1]), Equals, "/test{host=b}")
	c.Check(query.Aggregation[0].Type, Equals, oproto.StreamAggregation_MEAN)
	c.Check(query.Aggregation[0].Label[0], Equals, "xyz")

	ch, err := q.Run(context.Background(), s.store)
	c.Assert(err, IsNil)
	output := []*oproto.ValueStream{}
	for stream := range ch {
		output = append(output, stream)
	}

	c.Check(output[0].Value[0].DoubleValue, Equals, float64((20*1+40*1)/2))
	c.Check(output[0].Value[1].DoubleValue, Equals, float64((20*2+40*2)/2))
	c.Check(output[0].Value[2].DoubleValue, Equals, float64((20*3+40*3)/2))
	c.Check(output[0].Value[3].DoubleValue, Equals, float64((20*4+40*4)/2))
	c.Check(output[0].Value[4].DoubleValue, Equals, float64((20*5+40*5)/2))
	c.Check(output[0].Value[5].DoubleValue, Equals, float64((20*6+40*6)/2))
	c.Check(output[0].Value[6].DoubleValue, Equals, float64((20*7+40*7)/2))
	c.Check(output[0].Value[7].DoubleValue, Equals, float64((20*8+40*8)/2))
	c.Check(output[0].Value[8].DoubleValue, Equals, float64((20*9+40*9)/2))
	c.Check(output[0].Value[9].DoubleValue, Equals, float64((20*10+40*10)/2))
	c.Check(output[0].Value[10].DoubleValue, Equals, float64((20*11+40*11)/2))
}

func (s *MySuite) TestPercentile(c *C) {
	q, err := Parse("percentile(20) by (host) (/test{host=a})")
	c.Assert(err, IsNil)
	query := q.query
	c.Check(variable.ProtoToString(query.Aggregation[0].Query[0].Variable[0]), Equals, "/test{host=a}")
	c.Check(query.Aggregation[0].Type, Equals, oproto.StreamAggregation_PERCENTILE)
	c.Check(query.Aggregation[0].Param, Equals, 20.0)
	c.Check(query.Aggregation[0].Label[0], Equals, "host")
}

func (s *MySuite) TestMutation(c *C) {
	q, err := Parse("rate(/test{host=a})")
	c.Assert(err, IsNil)
	query := q.query
	c.Check(query.Mutation[0].Type, Equals, oproto.StreamMutation_RATE)
	c.Check(variable.ProtoToString(query.Mutation[0].Query.Variable[0]), Equals, "/test{host=a}")

	ch, err := q.Run(context.Background(), s.store)
	c.Assert(err, IsNil)
	numStreams := 0
	for stream := range ch {
		c.Check(variable.ProtoToString(stream.Variable), Equals, "/test{host=a}")
		c.Check(len(stream.Value), Equals, 10)
		numStreams++
	}
	c.Check(numStreams, Equals, 1)
}

func (s *MySuite) TestAggregationOfMutations(c *C) {
	q, err := Parse("mean by (host) (rate(/test{host=a}[1200:1500], /test{host=b}[1200:1500]))")
	c.Assert(err, IsNil)
	query := q.query
	c.Check(query.Aggregation[0].Type, Equals, oproto.StreamAggregation_MEAN)
	c.Check(query.Aggregation[0].Query[0].Mutation[0].Type, Equals, oproto.StreamMutation_RATE)
	c.Check(variable.ProtoToString(query.Aggregation[0].Query[0].Mutation[0].Query.Variable[0]), Equals, "/test{host=a}[1200:1500]")
	c.Check(variable.ProtoToString(query.Aggregation[0].Query[0].Mutation[0].Query.Variable[1]), Equals, "/test{host=b}[1200:1500]")
	c.Check(query.Aggregation[0].Query[0].Mutation[0].Query.Variable[0].MinTimestamp, Equals, int64(1200))
	c.Check(query.Aggregation[0].Query[0].Mutation[0].Query.Variable[0].MaxTimestamp, Equals, int64(1500))
	c.Check(query.Aggregation[0].Query[0].Mutation[0].Query.Variable[1].MinTimestamp, Equals, int64(1200))
	c.Check(query.Aggregation[0].Query[0].Mutation[0].Query.Variable[1].MaxTimestamp, Equals, int64(1500))
}

func (s *MySuite) TestAggregationOfPercentile(c *C) {
	q, err := Parse("percentile(90) by (host) (rate(/test{host=a}[1200:1500], /test{host=b}[1200:1500]))")
	c.Assert(err, IsNil)
	query := q.query
	c.Check(query.Aggregation[0].Type, Equals, oproto.StreamAggregation_PERCENTILE)
	c.Check(query.Aggregation[0].Param, Equals, 90.0)
	c.Check(query.Aggregation[0].Label[0], Equals, "host")
	c.Check(variable.ProtoToString(query.Aggregation[0].Query[0].Mutation[0].Query.Variable[0]), Equals, "/test{host=a}[1200:1500]")
	c.Check(variable.ProtoToString(query.Aggregation[0].Query[0].Mutation[0].Query.Variable[1]), Equals, "/test{host=b}[1200:1500]")
	c.Check(query.Aggregation[0].Query[0].Mutation[0].Query.Variable[0].MinTimestamp, Equals, int64(1200))
	c.Check(query.Aggregation[0].Query[0].Mutation[0].Query.Variable[0].MaxTimestamp, Equals, int64(1500))
	c.Check(query.Aggregation[0].Query[0].Mutation[0].Query.Variable[1].MinTimestamp, Equals, int64(1200))
	c.Check(query.Aggregation[0].Query[0].Mutation[0].Query.Variable[1].MaxTimestamp, Equals, int64(1500))
}

func (s *MySuite) TestPipeline(c *C) {
	return
	qs := `
		aggregation <
			type: MEAN
			query <
				mutation <
					type: MEAN
					query <
						variable <
							name: "/test"
							label <
								key: "host"
								value: "a"
							>
							label <
								key: "otherkey"
								value: "x"
							>
						>
						variable <
							name: "/test"
							label <
								key: "host"
								value: "b"
							>
							label <
								key: "otherkey"
								value: "y"
							>
						>
					>
				>
			>
		>
	`

	qp := &oproto.Query{}
	err := proto.UnmarshalText(qs, qp)
	c.Assert(err, IsNil)

	_, err = NewFromProto(qp).Run(context.Background(), s.store)
	c.Assert(err, IsNil)

	//c.Fail()
}
