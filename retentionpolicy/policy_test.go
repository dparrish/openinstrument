package retentionpolicy

import (
	"log"
	"testing"

	"code.google.com/p/goprotobuf/proto"

	"github.com/dparrish/openinstrument"
	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/variable"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

/*
func (s *MySuite) TestDefaultDropPolicy(c *C) {
	policyTxt := `
		interval: 3600
		policy {
			comment: "Throw everything away"
			variable { name: "*" }
			policy: DROP
		}
	`
	policyProto := &oproto.RetentionPolicy{}
	c.Assert(proto.UnmarshalText(policyTxt, policyProto), IsNil)
	policy := New(policyProto)
	log.Println(policyProto)

	input := make(chan *oproto.Value)
	output := policy.Apply(variable.NewFromString("/test/foo/bar"), input)

	for i := 0; i < 10; i++ {
		input <- &oproto.Value{Timestamp: proto.Uint64(uint64(i)), DoubleValue: proto.Float64(1.1)}
	}
	close(input)

	for value := range output {
		log.Printf("Got output when none was expected: %s", value)
		c.Fail()
	}
}
*/

func (s *MySuite) TestAgeKeepPolicy(c *C) {
	policyTxt := `
		interval: 1
		policy {
			comment: "Throw everything away"
			variable { name: "/test/foo*" }
			min_age: 75
			max_age: 91
			policy: KEEP
		}
		# Implicit DROP
	`
	policyProto := &oproto.RetentionPolicy{}
	c.Assert(proto.UnmarshalText(policyTxt, policyProto), IsNil)
	policy := New(policyProto)

	input := make(chan *oproto.Value)
	output := policy.Apply(variable.NewFromString("/test/foo/bar"), input)
	c.Assert(output, NotNil)

	now := openinstrument.NowMs()
	for i := 1; i <= 10; i++ {
		input <- &oproto.Value{
			Timestamp:    proto.Uint64(now - uint64(98-3*i)),
			EndTimestamp: proto.Uint64(now - uint64(100-3*i)),
			DoubleValue:  proto.Float64(1.1),
		}
	}
	close(input)

	count := 0
	for value := range output {
		age := now - value.GetTimestamp()
		if age < 75 || age > 91 {
			log.Printf("Got value outside expected age (%d)", age)
			c.Fail()
			continue
		}
		count++
	}
	c.Check(count, Equals, 5)
	//c.Fail()
}
