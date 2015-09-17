package protofile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/protobuf/proto"
	oproto "github.com/dparrish/openinstrument/proto"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestWriteFile(c *C) {
	filename := filepath.Join(c.MkDir(), "protofile_testwrite.dat")
	defer os.Remove(filename)

	{
		// Write a temporary file containing two labels
		writer, err := Write(filename)
		c.Assert(err, IsNil)
		defer writer.Close()
		c.Assert(writer.Tell(), Equals, int64(0))

		msg := &oproto.Label{
			Label: proto.String("greeting"),
			Value: proto.String("Hello world!"),
		}
		i, err := writer.Write(msg)
		c.Assert(err, IsNil)
		c.Assert(i, Equals, int64(32))
		j := i
		c.Assert(writer.Tell(), Equals, j)

		msg = &oproto.Label{
			Label: proto.String("greeting"),
			Value: proto.String("Hola!"),
		}
		i, err = writer.Write(msg)
		j += i
		c.Assert(i, Equals, int64(25))
		c.Assert(err, IsNil)
		c.Assert(writer.Tell(), Equals, j)

		// Write to a specific place in the file
		msg = &oproto.Label{
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

		msg := &oproto.Label{}
		i, err := reader.Read(msg)
		c.Assert(err, IsNil)
		c.Assert(i, Equals, int64(32))
		c.Assert(msg.GetLabel(), Equals, "greeting")
		c.Assert(msg.GetValue(), Equals, "Hello world!")
		j := i
		c.Assert(reader.Tell(), Equals, j)

		msg = &oproto.Label{}
		i, err = reader.Read(msg)
		j += i
		c.Assert(err, IsNil)
		c.Assert(i, Equals, int64(25))
		c.Assert(msg.GetLabel(), Equals, "greeting")
		c.Assert(msg.GetValue(), Equals, "Hola!")
		c.Assert(reader.Tell(), Equals, j)

		// Read the next message, which is after a few random bytes
		msg = &oproto.Label{}
		i, err = reader.Read(msg)
		j = 60 + i
		c.Assert(err, IsNil)
		c.Assert(i, Equals, int64(31))
		c.Assert(msg.GetLabel(), Equals, "greeting")
		c.Assert(msg.GetValue(), Equals, "Far out man")
		c.Assert(reader.Tell(), Equals, j)

		// Read from a specific place in the file
		msg = &oproto.Label{}
		i, err = reader.ReadAt(60, msg)
		c.Assert(err, IsNil)
		c.Assert(i, Equals, int64(31))
		c.Assert(msg.GetLabel(), Equals, "greeting")
		c.Assert(msg.GetValue(), Equals, "Far out man")
		c.Assert(reader.Tell(), Equals, j)

		// An attempt to read past the end of the file should return an error
		msg = &oproto.Label{}
		i, err = reader.Read(msg)
		c.Assert(err, ErrorMatches, "EOF")
		c.Assert(i, Equals, int64(0))
	}
}

func (s *MySuite) TestValueStreamReader(c *C) {
	filename := filepath.Join(c.MkDir(), "protofile_testvsr.dat")
	defer os.Remove(filename)

	{
		// Write a temporary file containing two value streams
		writer, err := Write(filename)
		c.Assert(err, IsNil)
		defer writer.Close()

		vs := &oproto.ValueStream{
			Variable: &oproto.StreamVariable{Name: proto.String("/test/bar")},
			Value: []*oproto.Value{
				{Timestamp: proto.Uint64(uint64(1)), DoubleValue: proto.Float64(1.1)},
				{Timestamp: proto.Uint64(uint64(2)), DoubleValue: proto.Float64(1.2)},
				{Timestamp: proto.Uint64(uint64(3)), DoubleValue: proto.Float64(1.3)},
			},
		}
		writer.Write(vs)

		vs = &oproto.ValueStream{
			Variable: &oproto.StreamVariable{Name: proto.String("/test/foo")},
			Value: []*oproto.Value{
				{Timestamp: proto.Uint64(uint64(1)), DoubleValue: proto.Float64(1.1)},
				{Timestamp: proto.Uint64(uint64(2)), DoubleValue: proto.Float64(1.2)},
				{Timestamp: proto.Uint64(uint64(3)), DoubleValue: proto.Float64(1.3)},
			},
		}
		writer.Write(vs)
	}

	{
		// Read back the contents and check
		file, err := Read(filename)
		c.Assert(err, IsNil)
		defer file.Close()
		reader := file.ValueStreamReader(500)
		vs := <-reader
		c.Check(vs.GetVariable().GetName(), Equals, "/test/bar")
		c.Check(vs.Value[0].GetDoubleValue(), Equals, 1.1)
		c.Check(vs.Value[1].GetDoubleValue(), Equals, 1.2)
		c.Check(vs.Value[2].GetDoubleValue(), Equals, 1.3)

		vs = <-reader
		c.Check(vs.GetVariable().GetName(), Equals, "/test/foo")
		c.Check(vs.Value[0].GetDoubleValue(), Equals, 1.1)
		c.Check(vs.Value[1].GetDoubleValue(), Equals, 1.2)
		c.Check(vs.Value[2].GetDoubleValue(), Equals, 1.3)
	}
}

func (s *MySuite) TestValueStreamWriter(c *C) {
	filename := filepath.Join(c.MkDir(), "protofile_testvsr.dat")
	defer os.Remove(filename)

	{
		// Write a temporary file containing two value streams
		file, err := Write(filename)
		c.Assert(err, IsNil)
		defer file.Close()
		writer, done := file.ValueStreamWriter(10)

		vs := &oproto.ValueStream{
			Variable: &oproto.StreamVariable{Name: proto.String("/test/bar")},
			Value: []*oproto.Value{
				{Timestamp: proto.Uint64(uint64(1)), DoubleValue: proto.Float64(1.1)},
				{Timestamp: proto.Uint64(uint64(2)), DoubleValue: proto.Float64(1.2)},
				{Timestamp: proto.Uint64(uint64(3)), DoubleValue: proto.Float64(1.3)},
			},
		}
		writer <- vs

		vs = &oproto.ValueStream{
			Variable: &oproto.StreamVariable{Name: proto.String("/test/foo")},
			Value: []*oproto.Value{
				{Timestamp: proto.Uint64(uint64(1)), DoubleValue: proto.Float64(1.1)},
				{Timestamp: proto.Uint64(uint64(2)), DoubleValue: proto.Float64(1.2)},
				{Timestamp: proto.Uint64(uint64(3)), DoubleValue: proto.Float64(1.3)},
			},
		}
		writer <- vs
		close(writer)
		<-done
	}

	{
		// Read back the contents and check
		file, err := Read(filename)
		c.Assert(err, IsNil)
		defer file.Close()
		reader := file.ValueStreamReader(500)
		vs := <-reader
		c.Check(vs.GetVariable().GetName(), Equals, "/test/bar")
		c.Check(vs.Value[0].GetDoubleValue(), Equals, 1.1)
		c.Check(vs.Value[1].GetDoubleValue(), Equals, 1.2)
		c.Check(vs.Value[2].GetDoubleValue(), Equals, 1.3)

		vs = <-reader
		c.Check(vs.GetVariable().GetName(), Equals, "/test/foo")
		c.Check(vs.Value[0].GetDoubleValue(), Equals, 1.1)
		c.Check(vs.Value[1].GetDoubleValue(), Equals, 1.2)
		c.Check(vs.Value[2].GetDoubleValue(), Equals, 1.3)
	}
}
