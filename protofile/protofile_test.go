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
		c.Assert(writer.Tell(), Equals, int64(0))

		msg := &openinstrument_proto.Label{
			Label: proto.String("greeting"),
			Value: proto.String("Hello world!"),
		}
		i, err := writer.Write(msg)
		c.Assert(err, IsNil)
		c.Assert(i, Equals, int64(32))
		j := i
		c.Assert(writer.Tell(), Equals, j)

		msg = &openinstrument_proto.Label{
			Label: proto.String("greeting"),
			Value: proto.String("Hola!"),
		}
		i, err = writer.Write(msg)
		j += i
		c.Assert(i, Equals, int64(25))
		c.Assert(err, IsNil)
		c.Assert(writer.Tell(), Equals, j)

		// Write to a specific place in the file
		msg = &openinstrument_proto.Label{
			Label: proto.String("greeting"),
			Value: proto.String("Far out man"),
		}
		i, err = writer.WriteAt(60, msg)
		c.Assert(i, Equals, int64(31))
		c.Assert(writer.Tell(), Equals, i+60)
		c.Assert(err, IsNil)

	}

	{
		// Read back the contents and ensure they are the same
		reader, err := Read(filename)
		c.Assert(err, IsNil)
		defer reader.Close()
		c.Assert(reader.Tell(), Equals, int64(0))

		msg := &openinstrument_proto.Label{}
		i, err := reader.Read(msg)
		c.Assert(err, IsNil)
		c.Assert(i, Equals, int64(32))
		c.Assert(msg.GetLabel(), Equals, "greeting")
		c.Assert(msg.GetValue(), Equals, "Hello world!")
		j := i
		c.Assert(reader.Tell(), Equals, j)

		msg = &openinstrument_proto.Label{}
		i, err = reader.Read(msg)
		j += i
		c.Assert(err, IsNil)
		c.Assert(i, Equals, int64(25))
		c.Assert(msg.GetLabel(), Equals, "greeting")
		c.Assert(msg.GetValue(), Equals, "Hola!")
		c.Assert(reader.Tell(), Equals, j)

		// Read the next message, which is after a few random bytes
		msg = &openinstrument_proto.Label{}
		i, err = reader.Read(msg)
		j = 60 + i
		c.Assert(err, IsNil)
		c.Assert(i, Equals, int64(31))
		c.Assert(msg.GetLabel(), Equals, "greeting")
		c.Assert(msg.GetValue(), Equals, "Far out man")
		c.Assert(reader.Tell(), Equals, j)

		// Read from a specific place in the file
		msg = &openinstrument_proto.Label{}
		i, err = reader.ReadAt(60, msg)
		c.Assert(err, IsNil)
		c.Assert(i, Equals, int64(31))
		c.Assert(msg.GetLabel(), Equals, "greeting")
		c.Assert(msg.GetValue(), Equals, "Far out man")
		c.Assert(reader.Tell(), Equals, j)

		// An attempt to read past the end of the file should return an error
		msg = &openinstrument_proto.Label{}
		i, err = reader.Read(msg)
		c.Assert(err, ErrorMatches, "EOF")
		c.Assert(i, Equals, int64(0))
	}
}
