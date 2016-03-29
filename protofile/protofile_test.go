package protofile

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/net/context"

	"github.com/dparrish/openinstrument"
	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/value"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestWriteFile(c *C) {
	filename := filepath.Join(c.MkDir(), "protofile_testwrite.dat")

	{
		// Write a temporary file containing two labels
		writer, err := Write(filename)
		c.Assert(err, IsNil)
		defer writer.Close()
		c.Check(writer.Tell(), Equals, int64(0))

		msg := &oproto.LogMessage{
			Timestamp: 1,
			Message:   "Hello world!",
		}
		i, err := writer.Write(msg)
		c.Assert(err, IsNil)
		c.Check(i, Equals, int64(24))
		j := i
		c.Check(writer.Tell(), Equals, j)

		msg = &oproto.LogMessage{
			Timestamp: 2,
			Message:   "Hola!",
		}
		i, err = writer.Write(msg)
		j += i
		c.Assert(err, IsNil)
		c.Check(i, Equals, int64(17))
		c.Check(writer.Tell(), Equals, j)

		// Write to a specific place in the file
		msg = &oproto.LogMessage{
			Timestamp: 3,
			Message:   "Far out man",
		}
		i, err = writer.WriteAt(60, msg)
		c.Assert(err, IsNil)
		c.Check(i, Equals, int64(23))
		c.Check(writer.Tell(), Equals, i+60)

	}

	{
		// Read back the contents and ensure they are the same
		reader, err := Read(filename)
		c.Assert(err, IsNil)
		defer reader.Close()
		c.Check(reader.Tell(), Equals, int64(0))

		msg := &oproto.LogMessage{}
		i, err := reader.Read(msg)
		c.Assert(err, IsNil)
		c.Check(i, Equals, int64(24))
		c.Check(msg.Timestamp, Equals, uint64(1))
		c.Check(msg.Message, Equals, "Hello world!")
		j := i
		c.Check(reader.Tell(), Equals, j)

		msg = &oproto.LogMessage{}
		i, err = reader.Read(msg)
		j += i
		c.Assert(err, IsNil)
		c.Check(i, Equals, int64(17))
		c.Check(msg.Timestamp, Equals, uint64(2))
		c.Check(msg.Message, Equals, "Hola!")
		c.Check(reader.Tell(), Equals, j)

		// Read the next message, which is after a few random bytes
		msg = &oproto.LogMessage{}
		i, err = reader.Read(msg)
		j = 60 + i
		c.Assert(err, IsNil)
		c.Check(i, Equals, int64(23))
		c.Check(msg.Timestamp, Equals, uint64(3))
		c.Check(msg.Message, Equals, "Far out man")
		c.Check(reader.Tell(), Equals, j)

		// Read from a specific place in the file
		msg = &oproto.LogMessage{}
		i, err = reader.ReadAt(60, msg)
		c.Assert(err, IsNil)
		c.Check(i, Equals, int64(23))
		c.Check(msg.Timestamp, Equals, uint64(3))
		c.Check(msg.Message, Equals, "Far out man")
		c.Check(reader.Tell(), Equals, j)

		// An attempt to read past the end of the file should return an error
		msg = &oproto.LogMessage{}
		i, err = reader.Read(msg)
		c.Assert(err, ErrorMatches, "EOF")
		c.Check(i, Equals, int64(0))
	}
}

func (s *MySuite) TestValueStreamReadWrite(c *C) {
	filename := filepath.Join(c.MkDir(), "protofile_testvar.dat")

	{
		// Write a temporary file containing two value streams
		file, err := Write(filename)
		c.Assert(err, IsNil)
		defer file.Close()
		writer, done := file.ValueStreamWriter(10)

		vs := &oproto.ValueStream{
			Variable: &oproto.StreamVariable{
				Name: "/test/bar",
				Label: map[string]string{
					"test": "bar",
				},
			},
			Value: []*oproto.Value{
				value.NewDouble(1, 1.1),
				value.NewDouble(2, 1.2),
				value.NewDouble(3, 1.3),
			},
		}
		writer <- vs

		vs = &oproto.ValueStream{
			Variable: &oproto.StreamVariable{Name: "/test/foo"},
			Value: []*oproto.Value{
				value.NewDouble(1, 1.1),
				value.NewDouble(2, 1.2),
				value.NewDouble(3, 1.3),
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
		reader := file.ValueStreamReader(context.Background(), 500)
		vs := <-reader
		c.Check(vs.Variable.Name, Equals, "/test/bar")
		c.Check(vs.Value[0].GetDouble(), Equals, 1.1)
		c.Check(vs.Value[1].GetDouble(), Equals, 1.2)
		c.Check(vs.Value[2].GetDouble(), Equals, 1.3)

		vs = <-reader
		c.Check(vs.Variable.Name, Equals, "/test/foo")
		c.Check(vs.Value[0].GetDouble(), Equals, 1.1)
		c.Check(vs.Value[1].GetDouble(), Equals, 1.2)
		c.Check(vs.Value[2].GetDouble(), Equals, 1.3)

		for range reader {
			log.Printf("Got unexpected value")
			c.Fail()
		}
	}
}

func (s *MySuite) TestValueStreamWriterMemoryLeak(c *C) {
	filename := filepath.Join(c.MkDir(), "protofile_testvar.dat")

	vs := &oproto.ValueStream{
		Variable: &oproto.StreamVariable{Name: "/test/bar"},
		Value:    []*oproto.Value{},
	}
	for i := 0; i < 10; i++ {
		vs.Value = append(vs.Value, value.NewDouble(openinstrument.NowMs(), 1.1))
	}
	// Write a temporary file containing two value streams
	file, err := Write(filename)
	c.Assert(err, IsNil)
	defer file.Close()
	writer, done := file.ValueStreamWriter(10)

	wroteStreams := 0
	for j := 0; j < 1000; j++ {
		writer <- vs
		wroteStreams++
	}

	close(writer)
	<-done

	file, err = Read(filename)
	c.Assert(err, IsNil)
	defer file.Close()

	reader := file.ValueStreamReader(context.Background(), 10)
	readStreams := 0
	for range reader {
		readStreams++
	}

	c.Assert(readStreams, Equals, wroteStreams)
}

func (s *MySuite) BenchmarkReader(c *C) {
	filename := filepath.Join(c.MkDir(), "protofile_testvar.dat")

	{
		// Write a temporary file containing lots of data
		file, err := Write(filename)
		c.Assert(err, IsNil)
		writer, done := file.ValueStreamWriter(10)

		for i := 0; i < 10000; i++ {
			vs := &oproto.ValueStream{
				Variable: &oproto.StreamVariable{Name: "/test/bar"},
				Value:    []*oproto.Value{value.NewDouble(uint64(i), float64(i))},
			}
			writer <- vs
		}
		close(writer)
		<-done
		file.Close()
	}

	for run := 0; run < c.N; run++ {
		// Read back the contents
		file, err := Read(filename)
		c.Assert(err, IsNil)
		defer file.Close()
		reader := file.ValueStreamReader(context.Background(), 500)
		for range reader {
		}
	}
}

// Write 100000 ValueStreams containing 1 value each
func (s *MySuite) BenchmarkWriterManyStreams(c *C) {
	filename := filepath.Join(c.MkDir(), "protofile_testvar.dat")

	vs := &oproto.ValueStream{
		Variable: &oproto.StreamVariable{Name: "/test/bar"},
		Value:    []*oproto.Value{value.NewDouble(1, float64(1.1))},
	}

	for run := 0; run < c.N; run++ {
		os.Remove(filename)
		file, err := Write(filename)
		c.Assert(err, IsNil)
		defer file.Close()
		writer, done := file.ValueStreamWriter(100000)

		for i := 0; i < 100000; i++ {
			writer <- vs
		}
		close(writer)
		<-done
	}

	file, _ := Read(filename)
	defer file.Close()
	stat, _ := file.Stat()
	log.Printf("BenchmarkWriterManyStreams wrote %d kB", stat.Size()/1024)
}

// Write 1 ValueStream with 100000 values
func (s *MySuite) BenchmarkWriterManyValues(c *C) {
	filename := filepath.Join(c.MkDir(), "protofile_testvar.dat")

	vs := &oproto.ValueStream{
		Variable: &oproto.StreamVariable{Name: "/test/bar"},
		Value:    []*oproto.Value{},
	}
	for j := 0; j < 100000; j++ {
		vs.Value = append(vs.Value, value.NewDouble(uint64(j), 1.1))
	}

	for run := 0; run < c.N; run++ {
		os.Remove(filename)
		file, err := Write(filename)
		c.Assert(err, IsNil)
		defer file.Close()
		writer, done := file.ValueStreamWriter(10)
		writer <- vs
		close(writer)
		<-done
	}

	file, _ := Read(filename)
	defer file.Close()
	stat, _ := file.Stat()
	log.Printf("BenchmarkWriterManyValues wrote %d kB", stat.Size()/1024)
}
