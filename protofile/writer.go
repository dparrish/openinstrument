package protofile

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/joaojeronimo/go-crc16"

	"code.google.com/p/goprotobuf/proto"
)

type Writer struct {
	filename string
	file     *os.File
	pos      int64
}

func Write(filename string) (*Writer, error) {
	writer := new(Writer)
	writer.filename = filename
	var err error
	writer.file, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0664)
	if err != nil {
		return nil, err
	}
	writer.file.Seek(0, os.SEEK_END)
	return writer, nil
}

func (pfw *Writer) Close() error {
	pfw.pos = 0
	return pfw.file.Close()
}

func (pfw *Writer) Tell() int64 {
	return pfw.pos
}

func (pfw *Writer) Stat() (os.FileInfo, error) {
	return pfw.file.Stat()
}

func (pfw *Writer) Sync() error {
	return pfw.file.Sync()
}

func (pfw *Writer) WriteAt(pos int64, message proto.Message) (int64, error) {
	var err error
	pfw.pos, err = pfw.file.Seek(pos, os.SEEK_SET)
	if err != nil {
		return 0, err
	}
	return pfw.Write(message)
}

func (pfw *Writer) Write(message proto.Message) (int64, error) {
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
		err = binary.Write(pfw.file, binary.LittleEndian, v)
		if err != nil {
			return 0, fmt.Errorf("Error writing entry to protofile: %s", err)
		}
		bytes += int64(binary.Size(v))
	}
	pfw.pos += bytes
	return bytes, nil
}
