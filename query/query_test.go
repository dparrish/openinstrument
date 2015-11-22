package query

import (
	"log"
	"testing"

	"github.com/dparrish/openinstrument/datastore"
	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/variable"
	"github.com/golang/protobuf/proto"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type FakeReadableStore struct{}

func (s *FakeReadableStore) Reader(v *variable.Variable) <-chan *oproto.ValueStream {
	c := make(chan *oproto.ValueStream, 100)
	go func() {
		defer close(c)
		switch v.String() {
		case "/test/offset":
			c <- &oproto.ValueStream{
				Variable: v.AsProto(),
				Value: []*oproto.Value{
					{Timestamp: uint64(0), DoubleValue: float64(20)},
					{Timestamp: uint64(61), DoubleValue: float64(40)},
					{Timestamp: uint64(122), DoubleValue: float64(60)},
					{Timestamp: uint64(185), DoubleValue: float64(80)},
					{Timestamp: uint64(241), DoubleValue: float64(100)},
					{Timestamp: uint64(299), DoubleValue: float64(122)},
					{Timestamp: uint64(330), DoubleValue: float64(132)},
					{Timestamp: uint64(359), DoubleValue: float64(140)},
					{Timestamp: uint64(421), DoubleValue: float64(160)},
					{Timestamp: uint64(488), DoubleValue: float64(180)},
					{Timestamp: uint64(540), DoubleValue: float64(200)},
					{Timestamp: uint64(975), DoubleValue: float64(275)},
				},
			}
		default:
			c <- &oproto.ValueStream{
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
	q, err := Parse("/test/foo")
	c.Assert(err, IsNil)
	query := q.query
	c.Check(variable.ProtoToString(query.Variable[0]), Equals, "/test/foo")
}

func (s *MySuite) TestVariableNoLabels(c *C) {
	q, err := Parse("/test/foo{}")
	c.Assert(err, IsNil)
	query := q.query
	c.Check(variable.ProtoToString(query.Variable[0]), Equals, "/test/foo")
	c.Check(query.Variable[0].MinTimestamp, Equals, int64(0))
	c.Check(query.Variable[0].MaxTimestamp, Equals, int64(0))

	ch, err := q.Run(s.store)
	c.Assert(err, IsNil)
	numStreams := 0
	for stream := range ch {
		log.Println(stream)
		c.Check(variable.ProtoToString(stream.Variable), Equals, "/test/foo")
		c.Check(len(stream.Value), Equals, 11)
		numStreams++
	}
	c.Check(numStreams, Equals, 1)
}

func (s *MySuite) TestVariableWithStartRange(c *C) {
	q, err := Parse("/test/foo[200]")
	c.Assert(err, IsNil)
	query := q.query
	c.Check(variable.ProtoToString(query.Variable[0]), Equals, "/test/foo[200]")
	c.Check(query.Variable[0].MinTimestamp, Equals, int64(200))
	c.Check(query.Variable[0].MaxTimestamp, Equals, int64(0))

	ch, err := q.Run(s.store)
	c.Assert(err, IsNil)
	numStreams := 0
	for stream := range ch {
		log.Println(stream)
		c.Check(variable.ProtoToString(stream.Variable), Equals, "/test/foo[200]")
		c.Check(len(stream.Value), Equals, 7)
		numStreams++
	}
	c.Check(numStreams, Equals, 1)
}

func (s *MySuite) TestVariableWithEndRange(c *C) {
	q, err := Parse("/test/foo{host=a}[200:400]")
	c.Assert(err, IsNil)
	query := q.query
	c.Check(variable.ProtoToString(query.Variable[0]), Equals, "/test/foo{host=a}[200:400]")
	c.Check(query.Variable[0].MinTimestamp, Equals, int64(200))
	c.Check(query.Variable[0].MaxTimestamp, Equals, int64(400))

	ch, err := q.Run(s.store)
	c.Assert(err, IsNil)
	numStreams := 0
	for stream := range ch {
		log.Println(stream)
		c.Check(variable.ProtoToString(stream.Variable), Equals, "/test/foo{host=a}[200:400]")
		c.Check(len(stream.Value), Equals, 3)
		numStreams++
	}
	c.Check(numStreams, Equals, 1)
}

func (s *MySuite) TestVariableOneLabel(c *C) {
	q, err := Parse("/test/foo{host=a}")
	c.Assert(err, IsNil)
	query := q.query
	c.Check(variable.ProtoToString(query.Variable[0]), Equals, "/test/foo{host=a}")
}

func (s *MySuite) TestVariableTwoLabels(c *C) {
	q, err := Parse("/test/foo{x=y,host=a}")
	c.Assert(err, IsNil)
	query := q.query
	c.Check(variable.ProtoToString(query.Variable[0]), Equals, "/test/foo{host=a,x=y}")
}

func (s *MySuite) TestLabelWildcard(c *C) {
	q, err := Parse("/test/foo{host=*}")
	c.Assert(err, IsNil)
	query := q.query
	c.Check(variable.ProtoToString(query.Variable[0]), Equals, "/test/foo{host=*}")
}

func (s *MySuite) TestAggregation(c *C) {
	q, err := Parse("mean by (host, xyz) (/test/foo{host=a}, /test/foo{host=b})")
	c.Assert(err, IsNil)
	query := q.query
	c.Check(variable.ProtoToString(query.Aggregation[0].Query[0].Variable[0]), Equals, "/test/foo{host=a}")
	c.Check(variable.ProtoToString(query.Aggregation[0].Query[0].Variable[1]), Equals, "/test/foo{host=b}")
	c.Check(query.Aggregation[0].Type, Equals, oproto.StreamAggregation_MEAN)
	c.Check(query.Aggregation[0].Label[0], Equals, "host")
	c.Check(query.Aggregation[0].Label[1], Equals, "xyz")
}

func (s *MySuite) TestPercentile(c *C) {
	q, err := Parse("percentile(20) by (host) (/test/foo{host=a})")
	c.Assert(err, IsNil)
	query := q.query
	c.Check(variable.ProtoToString(query.Aggregation[0].Query[0].Variable[0]), Equals, "/test/foo{host=a}")
	c.Check(query.Aggregation[0].Type, Equals, oproto.StreamAggregation_PERCENTILE)
	c.Check(query.Aggregation[0].Percentile, Equals, uint32(20))
	c.Check(query.Aggregation[0].Label[0], Equals, "host")
}

func (s *MySuite) TestMutation(c *C) {
	q, err := Parse("rate(/test/foo{host=a})")
	c.Assert(err, IsNil)
	query := q.query
	c.Check(query.Mutation[0].Type, Equals, oproto.StreamMutation_RATE)
	c.Check(variable.ProtoToString(query.Mutation[0].Query.Variable[0]), Equals, "/test/foo{host=a}")

	ch, err := q.Run(s.store)
	c.Assert(err, IsNil)
	numStreams := 0
	for stream := range ch {
		log.Println(stream)
		c.Check(variable.ProtoToString(stream.Variable), Equals, "/test/foo{host=a}")
		c.Check(len(stream.Value), Equals, 10)
		numStreams++
	}
	c.Check(numStreams, Equals, 1)
}

func (s *MySuite) TestAggregationOfMutations(c *C) {
	q, err := Parse("mean by (host) (rate(/test/foo{host=a}[1200:1500], /test/foo{host=b}[1200:1500]))")
	c.Assert(err, IsNil)
	query := q.query
	c.Check(query.Aggregation[0].Type, Equals, oproto.StreamAggregation_MEAN)
	c.Check(query.Aggregation[0].Query[0].Mutation[0].Type, Equals, oproto.StreamMutation_RATE)
	c.Check(variable.ProtoToString(query.Aggregation[0].Query[0].Mutation[0].Query.Variable[0]), Equals, "/test/foo{host=a}[1200:1500]")
	c.Check(variable.ProtoToString(query.Aggregation[0].Query[0].Mutation[0].Query.Variable[1]), Equals, "/test/foo{host=b}[1200:1500]")
	c.Check(query.Aggregation[0].Query[0].Mutation[0].Query.Variable[0].MinTimestamp, Equals, int64(1200))
	c.Check(query.Aggregation[0].Query[0].Mutation[0].Query.Variable[0].MaxTimestamp, Equals, int64(1500))
	c.Check(query.Aggregation[0].Query[0].Mutation[0].Query.Variable[1].MinTimestamp, Equals, int64(1200))
	c.Check(query.Aggregation[0].Query[0].Mutation[0].Query.Variable[1].MaxTimestamp, Equals, int64(1500))
}

func (s *MySuite) TestAggregationOfPercentile(c *C) {
	q, err := Parse("percentile(90) by (host) (rate(/test/foo{host=a}[1200:1500], /test/foo{host=b}[1200:1500]))")
	c.Assert(err, IsNil)
	query := q.query
	//fmt.Println(openinstrument.ProtoText(query))
	c.Check(query.Aggregation[0].Type, Equals, oproto.StreamAggregation_PERCENTILE)
	c.Check(query.Aggregation[0].Percentile, Equals, uint32(90))
	c.Check(query.Aggregation[0].Label[0], Equals, "host")
	c.Check(variable.ProtoToString(query.Aggregation[0].Query[0].Mutation[0].Query.Variable[0]), Equals, "/test/foo{host=a}[1200:1500]")
	c.Check(variable.ProtoToString(query.Aggregation[0].Query[0].Mutation[0].Query.Variable[1]), Equals, "/test/foo{host=b}[1200:1500]")
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
							name: "/test/foo"
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
							name: "/test/foo"
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

	output, err := NewFromProto(qp).Run(s.store)
	c.Assert(err, IsNil)

	log.Printf("Output from query:")
	for stream := range output {
		log.Println(stream)
	}

	//c.Fail()
}
