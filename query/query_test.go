package query

import (
	"log"
	"testing"

	"fmt"
	"time"

	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/variable"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestNewQuery(c *C) {
	q := New()
	q.AddVariableString("/test/foobar{host=rage}")
	c.Check("/test/foobar{host=rage}", Equals, variable.NewFromProto(q.AsProto().GetVariable()[0]).String())
	q.SetMinTimestamp(time.Now().Add(-30 * time.Minute)).SetMaxTimestamp(time.Now())
	//log.Println(q)

	txt := `
  variable: <
    name: "/test/foobar"
    label: <
      label: "host"
      value: "rage"
    >
  >
  `
	q = NewFromString(txt)
	//log.Println(q)
}

func (s *MySuite) TestDeepQuery(c *C) {
	// Sum /test/stat1 and /test/stat2 on each host
	q := New()
	for _, host := range []string{"host1", "host2", "host3", "host4"} {
		cluster := "a"
		if host == "host3" || host == "host4" {
			cluster = "b"
		}
		sq := New()
		sq.AddVariableString(fmt.Sprintf("/test/stat1{cluster=%s, hostname=%s}", cluster, host))
		sq.AddVariableString(fmt.Sprintf("/test/stat2{cluster=%s, hostname=%s}", cluster, host))
		sq.AddAggregation(oproto.StreamAggregation_SUM, nil)
		q.AddSubquery(sq.AsProto())
	}

	// Then get the average across each cluster
	q.AddAggregation(oproto.StreamAggregation_AVERAGE, []string{"cluster"})

	// Turn it into a 5 minute rate
	q.AddMutation(oproto.StreamMutation_RATE, 5*60*1000, 0)

	// Multiply the result by 8 (bps)
	finalQ := New()
	finalQ.AddSubquery(q.AsProto())
	finalQ.AddConstant(8)
	finalQ.AddMutation(oproto.StreamMutation_MULTIPLY, 0, 0)

	log.Println(finalQ)
}
