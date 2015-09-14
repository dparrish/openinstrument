package protofile

import (
	"log"
	"os"
)

const protoMagic uint16 = 0xDEAD

// ProtoFile contains information for reading and writing to a protofile.
// It is returned from either Read() or Write().
type ProtoFile struct {
	filename string
	file     *os.File
	stat     os.FileInfo
	pos      int64
}

// Close closes a ProtoFile opened by Read() or Write().
func (pf *ProtoFile) Close() error {
	pf.pos = 0
	return pf.file.Close()
}

func (pf *ProtoFile) Stat() (os.FileInfo, error) {
	return pf.file.Stat()
}

func (pf *ProtoFile) Tell() int64 {
	return pf.pos
}

func (pf *ProtoFile) Sync() error {
	return pf.file.Sync()
}

func (pf *ProtoFile) Seek(pos int64) (int64, error) {
	npos, err := pf.file.Seek(pos, os.SEEK_SET)
	if err != nil {
		log.Printf("Erorr seeking to %d: %s", pos, err)
		return npos, err
	}
	pf.pos = npos
	return npos, nil
}
