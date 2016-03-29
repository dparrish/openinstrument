package mutations_test

// Different package name to resolve cyclic dependencies

import (
	"math"
	"testing"

	"golang.org/x/net/context"

	"github.com/dparrish/openinstrument/datastore"
	. "github.com/dparrish/openinstrument/mutations"
	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/value"
	"github.com/dparrish/openinstrument/variable"
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

func checkValue(c *C, value *oproto.Value, expectedTimestamp uint64, expectedValue float64) {
	c.Assert(value, Not(IsNil))
	c.Check(value.GetDouble(), Equals, expectedValue)
	c.Check(int(value.Timestamp), Equals, int(expectedTimestamp))
}

func (s *MySuite) TestMean(c *C) {
	for input := range s.store.Reader(context.Background(), variable.NewFromString("/test")) {
		output := Mean(input)
		checkValue(c, output.Value[0], 60*0, float64((20*0+20*1)/2))
		checkValue(c, output.Value[1], 60*1, float64((20*1+20*2)/2))
		checkValue(c, output.Value[2], 60*2, float64((20*2+20*3)/2))
		checkValue(c, output.Value[3], 60*3, float64((20*3+20*4)/2))
		checkValue(c, output.Value[4], 60*4, float64((20*4+20*5)/2))
		checkValue(c, output.Value[5], 60*5, float64((20*5+20*6)/2))
		checkValue(c, output.Value[6], 60*6, float64((20*6+20*7)/2))
		checkValue(c, output.Value[7], 60*7, float64((20*7+20*8)/2))
		checkValue(c, output.Value[8], 60*8, float64((20*8+20*9)/2))
		checkValue(c, output.Value[9], 60*9, float64((20*9+20*10)/2))
		checkValue(c, output.Value[10], 60*10, float64((20*10+20*11)/2))
	}
}

func (s *MySuite) TestMovingAverage(c *C) {
	for input := range s.store.Reader(context.Background(), variable.NewFromString("/test")) {
		output := MovingAverage(120, input)
		// Length should be 2 less than the input, as there is no moving average across the first two elements
		c.Assert(len(output.Value), Equals, len(input.Value)-2)
		checkValue(c, output.Value[0], 60*2, 40)
		checkValue(c, output.Value[1], 60*3, 60)
		checkValue(c, output.Value[2], 60*4, 80)
		checkValue(c, output.Value[3], 60*5, 100)
		checkValue(c, output.Value[4], 60*6, 120)
		checkValue(c, output.Value[5], 60*7, 140)
		checkValue(c, output.Value[6], 60*8, 160)
		checkValue(c, output.Value[7], 60*9, 180)
		checkValue(c, output.Value[8], 60*10, 200)
	}
}

func (s *MySuite) TestMin(c *C) {
	for input := range s.store.Reader(context.Background(), variable.NewFromString("/test")) {
		output := Min(120, input)
		checkValue(c, output.Value[0], 0, 20*1)
		checkValue(c, output.Value[1], 120, 20*3)
		checkValue(c, output.Value[2], 240, 20*5)
		checkValue(c, output.Value[3], 360, 20*7)
		checkValue(c, output.Value[4], 480, 20*9)
	}
}

func (s *MySuite) TestMax(c *C) {
	for input := range s.store.Reader(context.Background(), variable.NewFromString("/test")) {
		output := Max(120, input)
		checkValue(c, output.Value[0], 60*2, 20*3)
		checkValue(c, output.Value[1], 60*4, 20*5)
		checkValue(c, output.Value[2], 60*6, 20*7)
		checkValue(c, output.Value[3], 60*8, 20*9)
		checkValue(c, output.Value[4], 60*10, 20*11)
	}
}

func (s *MySuite) TestInterpolate(c *C) {
	for input := range s.store.Reader(context.Background(), variable.NewFromString("/test/offset")) {
		output := Interpolate(300, input)
		checkValue(c, output.Value[0], 0, float64(20))
		checkValue(c, output.Value[1], 300, float64(121.81818181818181))
		checkValue(c, output.Value[2], 600, float64(191.86046511627907))
		checkValue(c, output.Value[3], 900, float64(258.37209302325584))
	}
}

func (s *MySuite) TestRate(c *C) {
	for input := range s.store.Reader(context.Background(), variable.NewFromString("/test")) {
		output := Rate(input)
		for _, v := range output.Value {
			checkValue(c, v, v.Timestamp, float64(1)/float64(3))
		}
	}
}

func (s *MySuite) TestSignedRate(c *C) {
	for input := range s.store.Reader(context.Background(), variable.NewFromString("/test")) {
		output := SignedRate(input)
		for _, v := range output.Value {
			checkValue(c, v, v.Timestamp, float64(1)/float64(3))
		}
	}
}

func (s *MySuite) TestFirst(c *C) {
	for input := range s.store.Reader(context.Background(), variable.NewFromString("/test")) {
		output := First(120, input)
		c.Assert(len(output.Value), Equals, 6)
		checkValue(c, output.Value[0], 60*0, 20*1)
		checkValue(c, output.Value[1], 60*2, 20*3)
		checkValue(c, output.Value[2], 60*4, 20*5)
		checkValue(c, output.Value[3], 60*6, 20*7)
		checkValue(c, output.Value[4], 60*8, 20*9)
		checkValue(c, output.Value[5], 60*10, 20*11)
	}
}

func (s *MySuite) TestLast(c *C) {
	for input := range s.store.Reader(context.Background(), variable.NewFromString("/test")) {
		output := Last(120, input)
		c.Assert(len(output.Value), Equals, 5)
		checkValue(c, output.Value[0], 60*2, 20*3)
		checkValue(c, output.Value[1], 60*4, 20*5)
		checkValue(c, output.Value[2], 60*6, 20*7)
		checkValue(c, output.Value[3], 60*8, 20*9)
		checkValue(c, output.Value[4], 60*10, 20*11)
	}
}

func (s *MySuite) TestMultiply(c *C) {
	for input := range s.store.Reader(context.Background(), variable.NewFromString("/test")) {
		output := Multiply(2, input)
		c.Assert(len(output.Value), Equals, len(input.Value))
		checkValue(c, output.Value[0], 60*0, 20*1*2)
		checkValue(c, output.Value[1], 60*1, 20*2*2)
		checkValue(c, output.Value[2], 60*2, 20*3*2)
		checkValue(c, output.Value[3], 60*3, 20*4*2)
		checkValue(c, output.Value[4], 60*4, 20*5*2)
		checkValue(c, output.Value[5], 60*5, 20*6*2)
		checkValue(c, output.Value[6], 60*6, 20*7*2)
		checkValue(c, output.Value[7], 60*7, 20*8*2)
		checkValue(c, output.Value[8], 60*8, 20*9*2)
		checkValue(c, output.Value[9], 60*9, 20*10*2)
		checkValue(c, output.Value[10], 60*10, 20*11*2)
	}
}

func (s *MySuite) TestDivide(c *C) {
	for input := range s.store.Reader(context.Background(), variable.NewFromString("/test")) {
		output := Multiply(1.0/2.0, input)
		c.Assert(len(output.Value), Equals, len(input.Value))
		checkValue(c, output.Value[0], 60*0, (20*1)/2.0)
		checkValue(c, output.Value[1], 60*1, (20*2)/2.0)
		checkValue(c, output.Value[2], 60*2, (20*3)/2.0)
		checkValue(c, output.Value[3], 60*3, (20*4)/2.0)
		checkValue(c, output.Value[4], 60*4, (20*5)/2.0)
		checkValue(c, output.Value[5], 60*5, (20*6)/2.0)
		checkValue(c, output.Value[6], 60*6, (20*7)/2.0)
		checkValue(c, output.Value[7], 60*7, (20*8)/2.0)
		checkValue(c, output.Value[8], 60*8, (20*9)/2.0)
		checkValue(c, output.Value[9], 60*9, (20*10)/2.0)
		checkValue(c, output.Value[10], 60*10, (20*11)/2.0)
	}
}

func (s *MySuite) TestAdd(c *C) {
	for input := range s.store.Reader(context.Background(), variable.NewFromString("/test")) {
		output := Add(5, input)
		c.Assert(len(output.Value), Equals, len(input.Value))
		checkValue(c, output.Value[0], 60*0, 20*1+5)
		checkValue(c, output.Value[1], 60*1, 20*2+5)
		checkValue(c, output.Value[2], 60*2, 20*3+5)
		checkValue(c, output.Value[3], 60*3, 20*4+5)
		checkValue(c, output.Value[4], 60*4, 20*5+5)
		checkValue(c, output.Value[5], 60*5, 20*6+5)
		checkValue(c, output.Value[6], 60*6, 20*7+5)
		checkValue(c, output.Value[7], 60*7, 20*8+5)
		checkValue(c, output.Value[8], 60*8, 20*9+5)
		checkValue(c, output.Value[9], 60*9, 20*10+5)
		checkValue(c, output.Value[10], 60*10, 20*11+5)
	}
}

func (s *MySuite) TestSubtract(c *C) {
	for input := range s.store.Reader(context.Background(), variable.NewFromString("/test")) {
		output := Add(-5, input)
		c.Assert(len(output.Value), Equals, len(input.Value))
		checkValue(c, output.Value[0], 60*0, 20*1-5)
		checkValue(c, output.Value[1], 60*1, 20*2-5)
		checkValue(c, output.Value[2], 60*2, 20*3-5)
		checkValue(c, output.Value[3], 60*3, 20*4-5)
		checkValue(c, output.Value[4], 60*4, 20*5-5)
		checkValue(c, output.Value[5], 60*5, 20*6-5)
		checkValue(c, output.Value[6], 60*6, 20*7-5)
		checkValue(c, output.Value[7], 60*7, 20*8-5)
		checkValue(c, output.Value[8], 60*8, 20*9-5)
		checkValue(c, output.Value[9], 60*9, 20*10-5)
		checkValue(c, output.Value[10], 60*10, 20*11-5)
	}
}

func (s *MySuite) TestRoot(c *C) {
	for input := range s.store.Reader(context.Background(), variable.NewFromString("/test")) {
		output := Root(2, input)
		c.Assert(len(output.Value), Equals, len(input.Value))
		checkValue(c, output.Value[0], 60*0, math.Sqrt(20*1))
		checkValue(c, output.Value[1], 60*1, math.Sqrt(20*2))
		checkValue(c, output.Value[2], 60*2, math.Sqrt(20*3))
		checkValue(c, output.Value[3], 60*3, math.Sqrt(20*4))
		checkValue(c, output.Value[4], 60*4, math.Sqrt(20*5))
		checkValue(c, output.Value[5], 60*5, math.Sqrt(20*6))
		checkValue(c, output.Value[6], 60*6, math.Sqrt(20*7))
		checkValue(c, output.Value[7], 60*7, math.Sqrt(20*8))
		checkValue(c, output.Value[8], 60*8, math.Sqrt(20*9))
		checkValue(c, output.Value[9], 60*9, math.Sqrt(20*10))
		checkValue(c, output.Value[10], 60*10, math.Sqrt(20*11))
	}
}

func (s *MySuite) TestPower(c *C) {
	for input := range s.store.Reader(context.Background(), variable.NewFromString("/test")) {
		output := Power(2, input)
		c.Assert(len(output.Value), Equals, len(input.Value))
		checkValue(c, output.Value[0], 60*0, math.Pow(20*1, 2))
		checkValue(c, output.Value[1], 60*1, math.Pow(20*2, 2))
		checkValue(c, output.Value[2], 60*2, math.Pow(20*3, 2))
		checkValue(c, output.Value[3], 60*3, math.Pow(20*4, 2))
		checkValue(c, output.Value[4], 60*4, math.Pow(20*5, 2))
		checkValue(c, output.Value[5], 60*5, math.Pow(20*6, 2))
		checkValue(c, output.Value[6], 60*6, math.Pow(20*7, 2))
		checkValue(c, output.Value[7], 60*7, math.Pow(20*8, 2))
		checkValue(c, output.Value[8], 60*8, math.Pow(20*9, 2))
		checkValue(c, output.Value[9], 60*9, math.Pow(20*10, 2))
		checkValue(c, output.Value[10], 60*10, math.Pow(20*11, 2))
	}
}

type FakeReadableStore struct{}

func (s *FakeReadableStore) Reader(ctx context.Context, v *variable.Variable) <-chan *oproto.ValueStream {
	c := make(chan *oproto.ValueStream, 100)
	go func() {
		defer close(c)
		switch v.String() {
		case "/test/offset":
			c <- &oproto.ValueStream{
				Variable: v.AsProto(),
				Value: []*oproto.Value{
					value.NewDouble(uint64(0), float64(20)),
					value.NewDouble(uint64(61), float64(40)),
					value.NewDouble(uint64(122), float64(60)),
					value.NewDouble(uint64(185), float64(80)),
					value.NewDouble(uint64(241), float64(100)),
					value.NewDouble(uint64(299), float64(122)),
					value.NewDouble(uint64(330), float64(132)),
					value.NewDouble(uint64(359), float64(140)),
					value.NewDouble(uint64(421), float64(160)),
					value.NewDouble(uint64(488), float64(180)),
					value.NewDouble(uint64(540), float64(200)),
					value.NewDouble(uint64(975), float64(275)),
				},
			}
		default:
			c <- &oproto.ValueStream{
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
	}()
	return c
}
