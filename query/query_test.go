package query

import (
	"fmt"
	"testing"

	"github.com/dparrish/openinstrument"
	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/variable"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestVariableNoLabels(c *C) {
	query, err := Parse("/test/foo{}")
	c.Assert(err, IsNil)
	c.Check(variable.ProtoToString(query.Variable[0]), Equals, "/test/foo")
}

func (s *MySuite) TestVariableNoLabelsOrBraces(c *C) {
	query, err := Parse("/test/foo")
	c.Assert(err, IsNil)
	c.Check(variable.ProtoToString(query.Variable[0]), Equals, "/test/foo")
}

func (s *MySuite) TestVariableOneLabel(c *C) {
	query, err := Parse("/test/foo{host=a}")
	c.Assert(err, IsNil)
	c.Check(variable.ProtoToString(query.Variable[0]), Equals, "/test/foo{host=a}")
}

func (s *MySuite) TestVariableTwoLabels(c *C) {
	query, err := Parse("/test/foo{x=y,host=a}")
	c.Assert(err, IsNil)
	c.Check(variable.ProtoToString(query.Variable[0]), Equals, "/test/foo{host=a,x=y}")
}

func (s *MySuite) TestAggregation(c *C) {
	query, err := Parse("mean by (host, xyz) (/test/foo{host=a}, /test/foo{host=b})")
	c.Assert(err, IsNil)
	c.Check(variable.ProtoToString(query.Aggregation[0].Variable[0]), Equals, "/test/foo{host=a}")
	c.Check(variable.ProtoToString(query.Aggregation[0].Variable[1]), Equals, "/test/foo{host=b}")
	c.Check(query.Aggregation[0].Type, Equals, oproto.StreamAggregation_MEAN)
	c.Check(query.Aggregation[0].Label[0], Equals, "host")
	c.Check(query.Aggregation[0].Label[1], Equals, "xyz")
}

func (s *MySuite) TestPercentile(c *C) {
	query, err := Parse("percentile(20) by (host) (/test/foo{host=a})")
	c.Assert(err, IsNil)
	c.Check(variable.ProtoToString(query.Aggregation[0].Variable[0]), Equals, "/test/foo{host=a}")
	c.Check(query.Aggregation[0].Type, Equals, oproto.StreamAggregation_PERCENTILE)
	c.Check(query.Aggregation[0].Percentile, Equals, uint32(20))
	c.Check(query.Aggregation[0].Label[0], Equals, "host")
}

func (s *MySuite) TestMutation(c *C) {
	query, err := Parse("rate(10s, /test/foo{host=a})")
	c.Assert(err, IsNil)
	c.Check(query.Mutation[0].SampleType, Equals, oproto.StreamMutation_RATE)
	c.Check(variable.ProtoToString(query.Mutation[0].Variable[0]), Equals, "/test/foo{host=a}")
}

func (s *MySuite) TestAggregationOfMutations(c *C) {
	query, err := Parse("mean by (host) (rate(10s, /test/foo{host=a}, /test/foo{host=b}))")
	c.Assert(err, IsNil)
	fmt.Println(openinstrument.ProtoText(query))
	c.Check(query.Aggregation[0].Mutation[0].SampleType, Equals, oproto.StreamMutation_RATE)
	c.Check(variable.ProtoToString(query.Aggregation[0].Mutation[0].Variable[0]), Equals, "/test/foo{host=a}")
	c.Check(variable.ProtoToString(query.Aggregation[0].Mutation[0].Variable[1]), Equals, "/test/foo{host=b}")
}
