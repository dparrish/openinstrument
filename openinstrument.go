package openinstrument

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/golang/protobuf/proto"
)

func NowMs() uint64 {
	return uint64(time.Now().UnixNano() / 1000000)
}

func ReadDirNames(directory string) ([]string, error) {
	dir, err := os.Open(directory)
	if err != nil {
		return nil, fmt.Errorf("Can't open %s for readdir: %s", directory, err)
	}
	defer dir.Close()
	names, err := dir.Readdirnames(0)
	if err != nil {
		return nil, fmt.Errorf("Can't read file names in %s: %s", directory, err)
	}
	sort.Strings(names)
	return names, nil
}

func ProtoText(msg proto.Message) string {
	buf := new(bytes.Buffer)
	if err := proto.MarshalText(buf, msg); err != nil {
		return ""
	}
	return buf.String()
}
