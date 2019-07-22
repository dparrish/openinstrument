package openinstrument

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
	return proto.MarshalTextString(msg)
}

func SafeWriteFile(filename string, msg proto.Message) error {
	f, err := ioutil.TempFile(filepath.Dir(filename), "tmpfile")
	if err != nil {
		return fmt.Errorf("can't create temporary file: %s", err)
	}
	w := bufio.NewWriter(f)
	if err := proto.MarshalText(w, msg); err != nil {
		f.Close()
		os.Remove(f.Name())
		return fmt.Errorf("marshal %s failed: %s", f.Name(), err)
	}
	w.Flush()
	f.Sync()
	f.Close()
	f.Sync()
	if err := os.Rename(f.Name(), filename); err != nil {
		os.Remove(f.Name())
		return fmt.Errorf("rename failed: %s", err)
	}
	f.Sync()
	return nil
}
