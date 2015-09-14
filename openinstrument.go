package openinstrument

import (
	"fmt"
	"os"
	"sort"
	"time"

	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/variable"
)

func NowMs() uint64 {
	return uint64(time.Now().UnixNano() / 1000000)
}

// ValueStreamWriter returns a channel that appends values to the supplied ValueStream.
// No effort is made to ensure that the ValueStream contains sorted Values
func ValueStreamWriter(stream *oproto.ValueStream) chan<- *oproto.Value {
	c := make(chan *oproto.Value)
	go func() {
		for value := range c {
			stream.Value = append(stream.Value, value)
		}
	}()
	return c
}

// ValueStreamReader returns a channel producing Values from the supplied ValueStream
func ValueStreamReader(stream *oproto.ValueStream) <-chan *oproto.Value {
	c := make(chan *oproto.Value)
	go func() {
		for _, value := range stream.Value {
			c <- value
		}
		close(c)
	}()
	return c
}

// MergeValueStreams merges multiple ValueStreams, returning a channel producing sorted Values.
func MergeValueStreams(streams []*oproto.ValueStream) <-chan *oproto.Value {
	c := make(chan *oproto.Value)
	n := len(streams)
	go func() {
		indexes := make([]int, n)
		for {
			var minTimestamp uint64
			var minStream *oproto.ValueStream
			var minValue *oproto.Value
			for i := 0; i < n; i++ {
				if indexes[i] >= len(streams[i].Value) {
					continue
				}
				v := streams[i].Value[indexes[i]]
				if minStream == nil || v.GetTimestamp() < minTimestamp {
					minTimestamp = v.GetTimestamp()
					minStream = streams[i]
					minValue = v
					indexes[i]++
				}
			}
			if minValue == nil {
				break
			}
			c <- minValue
		}
		close(c)
	}()
	return c
}

func MergeStreamsBy(streams []*oproto.ValueStream, by string) <-chan []*oproto.ValueStream {
	c := make(chan []*oproto.ValueStream)
	go func() {
		uniqueVars := make(map[string]bool)
		uniqueLabels := make(map[string]bool)
		for _, stream := range streams {
			v := variable.NewFromProto(stream.Variable)
			uniqueVars[v.Variable] = true
			labelValue, ok := v.Labels[by]
			if !ok {
				uniqueLabels[""] = true
			} else {
				uniqueLabels[labelValue] = true
			}
		}
		for varname := range uniqueVars {
			v := variable.NewFromString(varname)
			if by == "" {
				var output []*oproto.ValueStream
				for _, stream := range streams {
					testvar := variable.NewFromProto(stream.Variable)
					if testvar.Variable != v.Variable {
						continue
					}
					output = append(output, stream)
				}
				if len(output) > 0 {
					c <- output
				}
			} else {
				for labelvalue := range uniqueLabels {
					var output []*oproto.ValueStream
					for _, stream := range streams {
						testvar := variable.NewFromProto(stream.Variable)
						if testvar.Variable != v.Variable {
							continue
						}
						value, ok := testvar.Labels[by]
						if !ok {
							continue
						}
						if value != labelvalue {
							continue
						}
						output = append(output, stream)
					}
					if len(output) > 0 {
						c <- output
					}
				}
			}
		}
		close(c)
	}()
	return c
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
