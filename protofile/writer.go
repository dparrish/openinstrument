package protofile

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/joaojeronimo/go-crc16"

	"code.google.com/p/goprotobuf/proto"
	oproto "github.com/dparrish/openinstrument/proto"
)

// Write creates a file handle for writing a protofile, returning a ProtoFile.
// After calling this, the file is opened for writing and the file position it the end of the file, ready for appending.
func Write(filename string) (ReaderWriter, error) {
	writer := &ProtoFile{}
	writer.filename = filename
	var err error
	writer.file, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0664)
	if err != nil {
		return nil, err
	}
	writer.file.Seek(0, os.SEEK_END)
	return writer, nil
}

func (pf *ProtoFile) WriteAt(pos int64, message proto.Message) (int64, error) {
	var err error
	pf.pos, err = pf.file.Seek(pos, os.SEEK_SET)
	if err != nil {
		return 0, err
	}
	return pf.Write(message)
}

func (pf *ProtoFile) Write(message proto.Message) (int64, error) {
	data, err := proto.Marshal(message)
	if err != nil {
		return 0, fmt.Errorf("Marshaling error: %s", err)
	}
	var buf = []interface{}{
		uint16(protoMagic),
		uint32(len(data)),
		data,
		crc16.Crc16(data),
	}
	var bytes int64
	for _, v := range buf {
		err = binary.Write(pf.file, binary.LittleEndian, v)
		if err != nil {
			return 0, fmt.Errorf("Error writing entry to protofile: %s", err)
		}
		bytes += int64(binary.Size(v))
	}
	pf.pos += bytes
	return bytes, nil
}

func (pf *ProtoFile) ValueStreamWriter(chanSize int) (chan<- *oproto.ValueStream, <-chan interface{}) {
	c := make(chan *oproto.ValueStream, chanSize)
	done := make(chan interface{}, 1)
	go func() {
		defer close(done)
		for value := range c {
			_, err := pf.Write(value)
			if err == io.EOF {
				return
			}
			if err != nil {
				log.Println(err)
				return
			}
		}
	}()
	return c, done
}
