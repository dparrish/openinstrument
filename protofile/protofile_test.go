package protofile

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"code.google.com/p/goprotobuf/proto"
	openinstrument_proto "github.com/dparrish/openinstrument/proto"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestWriteFile(c *C) {
	tmpdir, _ := ioutil.TempDir("", "protofile_test")
	defer os.Remove(tmpdir)
	filename := filepath.Join(tmpdir, "protofile_test.dat")
	defer os.Remove(filename)

	{
		// Write a temporary file containing two labels
		writer, err := Write(filename)
		c.Assert(err, IsNil)
		defer writer.Close()

		msg := &openinstrument_proto.Label{
			Label: proto.String("greeting"),
			Value: proto.String("Hello world!"),
		}
		i, err := writer.Write(msg)
		c.Assert(err, IsNil)
		c.Assert(i, Equals, int64(32))

		msg = &openinstrument_proto.Label{
			Label: proto.String("greeting"),
			Value: proto.String("Hola!"),
		}
		i, err = writer.Write(msg)
		c.Assert(i, Equals, int64(25))
		c.Assert(err, IsNil)

		// Write to a specific place in the file
		msg = &openinstrument_proto.Label{
			Label: proto.String("greeting"),
			Value: proto.String("Far out man"),
		}
		i, err = writer.WriteAt(60, msg)
		c.Assert(i, Equals, int64(31))
		c.Assert(err, IsNil)

	}

	{
		// Read back the contents and ensure they are the same
		reader, err := Read(filename)
		c.Assert(err, IsNil)
		defer reader.Close()

		msg := &openinstrument_proto.Label{}
		i, err := reader.Read(msg)
		c.Assert(err, IsNil)
		c.Assert(i, Equals, int64(32))
		c.Assert(msg.GetLabel(), Equals, "greeting")
		c.Assert(msg.GetValue(), Equals, "Hello world!")

		msg = &openinstrument_proto.Label{}
		i, err = reader.Read(msg)
		c.Assert(err, IsNil)
		c.Assert(i, Equals, int64(25))
		c.Assert(msg.GetLabel(), Equals, "greeting")
		c.Assert(msg.GetValue(), Equals, "Hola!")

		// Read the next message, which is after a few random bytes
		msg = &openinstrument_proto.Label{}
		i, err = reader.Read(msg)
		c.Assert(err, IsNil)
		c.Assert(i, Equals, int64(31))
		c.Assert(msg.GetLabel(), Equals, "greeting")
		c.Assert(msg.GetValue(), Equals, "Far out man")

		// Read from a specific place in the file
		msg = &openinstrument_proto.Label{}
		i, err = reader.ReadAt(60, msg)
		c.Assert(err, IsNil)
		c.Assert(i, Equals, int64(31))
		c.Assert(msg.GetLabel(), Equals, "greeting")
		c.Assert(msg.GetValue(), Equals, "Far out man")

		// An attempt to read past the end of the file should return an error
		msg = &openinstrument_proto.Label{}
		i, err = reader.Read(msg)
		c.Assert(err, ErrorMatches, "EOF")
		c.Assert(i, Equals, int64(0))
	}
}
