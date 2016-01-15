package variable

import (
	"testing"

	"github.com/dparrish/openinstrument"
	oproto "github.com/dparrish/openinstrument/proto"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestVariableCreation(c *C) {
	newvar := NewFromString("/openinstrument/test{label1=value1,label2=spaced value}")
	c.Check(newvar.String(), Equals, "/openinstrument/test{label1=value1,label2=spaced value}")

	labels := map[string]string{
		"label1": "value1",
		"label2": "spaced value",
	}
	newvar = NewFromProto(&oproto.StreamVariable{
		Name:  "/openinstrument/test",
		Label: labels,
	})
	c.Check(newvar.String(), Equals, "/openinstrument/test{label1=value1,label2=spaced value}")

	p := newvar.AsProto()
	newvar.ParseFromProto(p)
	c.Check(newvar.String(), Equals, "/openinstrument/test{label1=value1,label2=spaced value}")
}

func (s *MySuite) TestToProto(c *C) {
	newvar := NewFromString("/openinstrument/test{label1=value1,label2=spaced value}")
	p := newvar.AsProto()
	newvar.ParseFromProto(p)
	c.Check(newvar.String(), Equals, "/openinstrument/test{label1=value1,label2=spaced value}")
}

func (s *MySuite) TestMatch(c *C) {
	var1 := NewFromString("/openinstrument/test{label1=value1,label2=spaced value}")
	c.Check(var1.Match(NewFromString("*")), Equals, true)
	c.Check(var1.Match(NewFromString("/openinstrument/*")), Equals, true)
	c.Check(var1.Match(NewFromString("/openinstrx/*")), Equals, false)
	c.Check(var1.Match(NewFromString("/openinstrument/")), Equals, false)
	c.Check(var1.Match(NewFromString("/openinstrument/test")), Equals, true)
	c.Check(var1.Match(NewFromString("/openinstrument/test{label1=value1}")), Equals, true)
	c.Check(var1.Match(NewFromString("/openinstrument/test{label1=value2}")), Equals, false)
	c.Check(var1.Match(NewFromString("/openinstrument/test{label1=*}")), Equals, true)
	c.Check(var1.Match(NewFromString("/openinstrument/test{label1=}")), Equals, false)
	c.Check(var1.Match(NewFromString("/openinstrument/test{label3=*}")), Equals, false)
}

func (s *MySuite) TestMatchRegexp(c *C) {
	v := NewFromString("/openinstrument/test{label1=/^.*foobar$/}")
	c.Assert(v.Labels["label1"], Equals, "/^.*foobar$/")
	c.Check(NewFromString("/openinstrument/test{label1=testing}").Match(v), Equals, false)
	c.Check(NewFromString("/openinstrument/test{label2=testing}").Match(v), Equals, false)
	c.Check(NewFromString("/openinstrument/test{label1=testfoobaring}").Match(v), Equals, false)
	c.Check(NewFromString("/openinstrument/test{label1=testingfoobar}").Match(v), Equals, true)
}

func (s *MySuite) TestInvalidVariable(c *C) {
	c.Check(NewFromProto(nil).String(), Equals, "")
	c.Check(NewFromProto(&oproto.StreamVariable{}).String(), Equals, "")
	v := &oproto.StreamVariable{Label: map[string]string{"foo": "bar"}}
	c.Check(NewFromProto(v).String(), Equals, "{foo=bar}")
}

func (s *MySuite) TestStringWithRange(c *C) {
	// Specify start and end
	v := NewFromString("/test[1200:1500]")
	c.Check(v.MinTimestamp, Equals, int64(1200))
	c.Check(v.MaxTimestamp, Equals, int64(1500))
	c.Check(v.String(), Equals, "/test[1200:1500]")

	// Only specify start
	v = NewFromString("/test[1200]")
	c.Check(v.MinTimestamp, Equals, int64(1200))
	c.Check(v.String(), Equals, "/test[1200]")

	// Relative timstamps
	v = NewFromString("/test[-1500:-1200]")
	c.Check(v.MinTimestamp, Equals, int64(-1500))
	c.Check(v.MaxTimestamp, Equals, int64(-1200))
	c.Check(v.String(), Equals, "/test[-1500:-1200]")

	// Relative start only
	v = NewFromString("/test[-1200]")
	c.Check(v.MinTimestamp, Equals, int64(-1200))
	c.Check(v.String(), Equals, "/test[-1200]")

	// Invalid timestamp order
	v = NewFromString("/test[-1200:-1500]")
	c.Check(v, IsNil)
}

func (s *MySuite) TestVariableInsideRange(c *C) {
	v := NewFromString("/test")
	v.MinTimestamp = 500
	v.MaxTimestamp = 1000
	c.Check(v.TimestampInsideRange(400), Equals, false)
	c.Check(v.TimestampInsideRange(600), Equals, true)
	c.Check(v.TimestampInsideRange(1001), Equals, false)
}

func (s *MySuite) TestVariableInsideRelativeRange(c *C) {
	v := NewFromString("/test")
	// From 40 to 50 seconds ago
	v.MinTimestamp = -50000
	v.MaxTimestamp = -40000
	c.Check(v.TimestampInsideRange(openinstrument.NowMs()-60000), Equals, false)
	c.Check(v.TimestampInsideRange(openinstrument.NowMs()-45000), Equals, true)
	c.Check(v.TimestampInsideRange(openinstrument.NowMs()-20000), Equals, false)
}
