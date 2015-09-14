package protofile

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"code.google.com/p/goprotobuf/proto"
	openinstrument_proto "github.com/dparrish/openinstrument/proto"
	"github.com/joaojeronimo/go-crc16"
)

var (
	checkProtofileCrc = flag.Bool("check_protofile_crc", true, "Check CRC of each protofile entry.")
)

type protoFileReader struct {
	filename string
	file     *os.File
	stat     os.FileInfo
}

func Read(filename string) (*protoFileReader, error) {
	reader := new(protoFileReader)
	reader.filename = filename
	var err error
	if reader.file, err = os.Open(filename); err != nil {
		log.Printf("Couldn't open %s: %s", filename, err)
		return nil, err
	}
	if reader.stat, err = reader.file.Stat(); err != nil {
		log.Printf("Couldn't stat %s: %s", filename, err)
		reader.file.Close()
		return nil, err
	}
	return reader, nil
}

func (pfr *protoFileReader) Close() error {
	return pfr.file.Close()
}

func (pfr *protoFileReader) Tell() int64 {
	pos, err := pfr.file.Seek(0, os.SEEK_CUR)
	if err != nil {
		log.Printf("Erorr getting file position: %s", err)
	}
	return pos
}

func (pfr *protoFileReader) Seek(pos int64) (int64, error) {
	npos, err := pfr.file.Seek(pos, os.SEEK_SET)
	if err != nil {
		log.Printf("Erorr getting seeking to %d: %s", pos, err)
		return npos, err
	}
	return npos, nil
}

func (pfr *protoFileReader) Stat() (os.FileInfo, error) {
	return pfr.file.Stat()
}

func (pfr *protoFileReader) ReadAt(pos int64, message proto.Message) (int64, error) {
	if _, err := pfr.Seek(pos); err != nil {
		return 0, err
	}
	return pfr.Read(message)
}

func (pfr *protoFileReader) ValueStreamReader(chanSize int) chan *openinstrument_proto.ValueStream {
	c := make(chan *openinstrument_proto.ValueStream, chanSize)
	go func() {
		for {
			value := new(openinstrument_proto.ValueStream)
			_, err := pfr.Read(value)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println(err)
				break
			}
			c <- value
		}
		close(c)
	}()
	return c
}

func (pfr *protoFileReader) ValueStreamReaderUntil(maxPos uint64, chanSize int) chan *openinstrument_proto.ValueStream {
	c := make(chan *openinstrument_proto.ValueStream, chanSize)
	go func() {
		for uint64(pfr.Tell()) < maxPos {
			value := new(openinstrument_proto.ValueStream)
			_, err := pfr.Read(value)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println(err)
				break
			}
			c <- value
		}
		close(c)
	}()
	return c
}

func (pfr *protoFileReader) Read(message proto.Message) (int64, error) {
	for {
		pos := pfr.Tell()
		type header struct {
			Magic  uint16
			Length uint32
		}
		var h header

		err := binary.Read(pfr.file, binary.LittleEndian, &h)
		if err != nil {
			if err == io.EOF {
				return 0, io.EOF
			}
			log.Printf("Error reading record header from recordlog: %s", err)
			return 0, err
		}

		// Read Magic header
		if h.Magic != protoMagic {
			log.Printf("Protobuf delimeter at %s:%x does not match %#x", pfr.filename, pos, protoMagic)
			pfr.Seek(pos + 1)
			continue
		}
		if int64(h.Length) >= pfr.stat.Size() {
			log.Printf("Chunk length %d at %s:%x is greater than file size %d", h.Length, pfr.filename, pos, pfr.stat.Size())
			continue
		}

		// Read Proto
		buf := make([]byte, h.Length)
		n, err := pfr.file.Read(buf)
		if err != nil || uint32(n) != h.Length {
			log.Printf("Could not read %d bytes from file: %s", h.Length, err)
			return 0, io.EOF
		}

		// Read CRC
		var crc uint16
		err = binary.Read(pfr.file, binary.LittleEndian, &crc)
		if err != nil {
			log.Printf("Error reading CRC from recordlog: %s", err)
			continue
		}
		if *checkProtofileCrc {
			checkcrc := crc16.Crc16(buf)
			if checkcrc != crc {
				log.Printf("CRC %x does not match %x", crc, checkcrc)
				continue
			}
		}

		// Decode and add proto
		if err = proto.Unmarshal(buf, message); err != nil {
			return 0, fmt.Errorf("Error decoding protobuf at %s:%x: %s", pfr.filename, pos, err)
		}

		return pfr.Tell() - pos, nil
	}
}
