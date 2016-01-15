package retentionpolicy

import (
	"log"
	"testing"

	"golang.org/x/net/context"

	"github.com/golang/protobuf/proto"

	"github.com/dparrish/openinstrument"
	"github.com/dparrish/openinstrument/datastore"
	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/variable"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestDefaultDropPolicy(c *C) {
	policyTxt := `
		interval: 3600
		policy {
			comment: "Throw everything away"
			policy: DROP
		}
	`
	policyProto := &oproto.RetentionPolicy{}
	c.Assert(proto.UnmarshalText(policyTxt, policyProto), IsNil)
	policy := New(policyProto)

	input := &oproto.ValueStream{
		Variable: variable.NewFromString("/test/foo/bar").AsProto(),
		Value:    []*oproto.Value{},
	}
	for i := 0; i < 10; i++ {
		input.Value = append(input.Value, &oproto.Value{Timestamp: uint64(i), DoubleValue: 1.1})
	}
	output := policy.Apply(input)

	for _, value := range output.Value {
		log.Printf("Got output when none was expected: %s", value)
		c.Fail()
	}
}

func (s *MySuite) TestDropByValue(c *C) {
	policyTxt := `
		interval: 3600
		policy {
			variable: {
				name: "/system/filesystem/*",
				label {
					key: "device"
					value: "/@zfs-auto-snap/"
				}
			}
			policy: DROP
		}
		policy {
			policy: KEEP
		}
	`
	policyProto := &oproto.RetentionPolicy{}
	c.Assert(proto.UnmarshalText(policyTxt, policyProto), IsNil)
	policy := New(policyProto)
	//log.Println(policyProto)

	// Device not matching regexp is kept
	input := &oproto.ValueStream{
		Variable: variable.NewFromString("/system/filesystem/available{device=r2/home}").AsProto(),
		Value:    []*oproto.Value{{Timestamp: uint64(1), DoubleValue: 1.0}},
	}
	c.Check(len(policy.Apply(input).Value), Equals, 1)

	// Device matching regexp is dropped
	input = &oproto.ValueStream{
		Variable: variable.NewFromString("/system/filesystem/used{device=r2/home@zfs-auto-snap_hourly}").AsProto(),
		Value:    []*oproto.Value{{Timestamp: uint64(1), DoubleValue: 1.0}},
	}
	c.Check(len(policy.Apply(input).Value), Equals, 0)
}

func (s *MySuite) TestAgeKeepPolicy(c *C) {
	return
	// TODO(dparrish): Re-enable this when it works
	policyTxt := `
			interval: 1
			policy {
				comment: "Throw everything away"
				variable {
					name: "/test/foo*"
					min_timestamp: -91
					max_timestamp: -75
				}
				policy: KEEP
			}
			# Implicit DROP
		`
	policyProto := &oproto.RetentionPolicy{}
	c.Assert(proto.UnmarshalText(policyTxt, policyProto), IsNil)
	policy := New(policyProto)

	input := &oproto.ValueStream{
		Variable: variable.NewFromString("/test/foo/bar").AsProto(),
		Value:    []*oproto.Value{},
	}
	now := openinstrument.NowMs()
	for i := 1; i <= 10; i++ {
		input.Value = append(input.Value, &oproto.Value{
			Timestamp:    now - uint64(98-3*i),
			EndTimestamp: now - uint64(100-3*i),
			DoubleValue:  1.1,
		})
	}

	output := policy.Apply(input)

	count := 0
	for _, value := range output.Value {
		age := now - value.Timestamp
		if age < 75 || age > 91 {
			log.Printf("Got value outside expected age (%d)", age)
			c.Fail()
			continue
		}
		count++
	}
	c.Check(count, Equals, 5)
}

func (s *MySuite) TestApplyToBlock(c *C) {
	return
	block := datastore.NewBlock(context.Background(), "/system/vmstat/nr_anon_transparent_hugepages{hostname=rage}", "bed18417-dd30-4ab4-6432-0635e0e7a2a7", "/r2/services/openinstrument/task1")

	streams := make([]*oproto.ValueStream, 0)
	originalNumValues := 0
	ch, err := block.GetIndexedStreams(context.Background())
	c.Assert(err, IsNil)
	for stream := range ch {
		streams = append(streams, stream)
		originalNumValues += len(stream.Value)
	}
	log.Printf("Read %d streams containing %d values", len(streams), originalNumValues)

	policyTxt := `
			interval: 1
			policy {
				comment: "Keep everything"
				variable {
					name: "*"
				}
				mutation {
					type: MAX
					sample_frequency: 300000
				}
			}
			policy {
				comment: "Throw everything away"
				policy: DROP
			}
		`
	policyProto := &oproto.RetentionPolicy{}
	c.Assert(proto.UnmarshalText(policyTxt, policyProto), IsNil)
	policy := New(policyProto)

	numValues := 0
	for _, stream := range streams {
		for range policy.Apply(stream).Value {
			numValues++
		}
	}
	log.Printf("After policy application, there are %d values left", numValues)

	c.Check(numValues, Equals, originalNumValues)
}
