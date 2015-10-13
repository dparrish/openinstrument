package protofile

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/variable"
	"github.com/golang/protobuf/proto"
	"github.com/joaojeronimo/go-crc16"
)

var (
	checkProtofileCrc = flag.Bool("check_protofile_crc", true, "Check CRC of each protofile entry.")
)

func Read(filename string) (ReaderWriter, error) {
	reader := new(ProtoFile)
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

func (pf *ProtoFile) ReadAt(pos int64, message proto.Message) (int64, error) {
	var err error
	if pf.pos, err = pf.Seek(pos); err != nil {
		return 0, err
	}
	return pf.Read(message)
}

func (pf *ProtoFile) ValueStreamReader(chanSize int) <-chan *oproto.ValueStream {
	c := make(chan *oproto.ValueStream, 10000)
	go func() {
		for {
			value := &oproto.ValueStream{}
			_, err := pf.Read(value)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println(err)
				break
			}
			value.VariableName = variable.ProtoToString(value.Variable)
			c <- value
		}
		close(c)
	}()
	return c
}

func (pf *ProtoFile) Read(message proto.Message) (int64, error) {
	for {
		h := blockHeader{}
		var bytesRead int64

		err := binary.Read(pf.file, binary.LittleEndian, &h)
		if err != nil {
			if err == io.EOF {
				return 0, io.EOF
			}
			log.Printf("Error reading record header from recordlog: %s", err)
			return 0, err
		}
		bytesRead += 6

		// Read Magic header
		if h.Magic != protoMagic {
			//log.Printf("Protobuf delimeter at %s:%x does not match %#x", pf.filename, pf.pos, protoMagic)
			pf.Seek(pf.pos + 1)
			continue
		}
		if int64(h.Length) >= pf.stat.Size() {
			log.Printf("Chunk length %d at %s:%x is greater than file size %d", h.Length, pf.filename, pf.pos, pf.stat.Size())
			continue
		}

		// Read Proto
		buf := make([]byte, h.Length)
		n, err := pf.file.Read(buf)
		bytesRead += int64(n)
		if err != nil || uint32(n) != h.Length {
			log.Printf("Could not read %d bytes from file: %s", h.Length, err)
			return 0, io.EOF
		}

		// Read CRC
		var crc uint16
		err = binary.Read(pf.file, binary.LittleEndian, &crc)
		bytesRead += int64(2)
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
			return 0, fmt.Errorf("Error decoding protobuf at %s:%x: %s", pf.filename, pf.pos, err)
		}

		pf.pos += bytesRead
		return bytesRead, nil
	}
}
