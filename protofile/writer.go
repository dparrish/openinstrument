package protofile

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/joaojeronimo/go-crc16"

	"github.com/golang/protobuf/proto"
)

// Write creates a file handle for writing a protofile, returning a ProtoFile.
// After calling this, the file is opened for writing and the file position it the end of the file, ready for appending.
func Write(filename string) (ReaderWriter, error) {
	pf := &ProtoFile{filename: filename}
	var err error
	if pf.file, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0664); err != nil {
		return nil, err
	}
	pf.file.Seek(0, os.SEEK_END)
	return pf, nil
}

func (pf *ProtoFile) WriteAt(pos int64, message proto.Message) (int64, error) {
	if npos, err := pf.Seek(pos); err != nil || npos != pos {
		return 0, err
	}
	return pf.Write(message)
}

func (pf *ProtoFile) Write(message proto.Message) (int64, error) {
	data, err := proto.Marshal(message)
	if err != nil {
		return 0, fmt.Errorf("Marshaling error: %s", err)
	}
	if binary.Write(pf.file, binary.LittleEndian, protoMagic) != nil {
		return 0, fmt.Errorf("Error writing entry to protofile: %s", err)
	}
	if binary.Write(pf.file, binary.LittleEndian, uint32(len(data))) != nil {
		return 0, fmt.Errorf("Error writing entry to protofile: %s", err)
	}
	if binary.Write(pf.file, binary.LittleEndian, data) != nil {
		return 0, fmt.Errorf("Error writing entry to protofile: %s", err)
	}
	if binary.Write(pf.file, binary.LittleEndian, crc16.Crc16(data)) != nil {
		return 0, fmt.Errorf("Error writing entry to protofile: %s", err)
	}
	bytes := int64(4 + 2 + len(data) + 2)
	pf.pos += bytes
	return bytes, nil
}
