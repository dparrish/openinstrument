package valuestream

import (
	"testing"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestSort(c *C) {
}

func (s *MySuite) TestMerge(c *C) {
}

func (s *MySuite) TestMergeBy(c *C) {
}
