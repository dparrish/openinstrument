package openinstrument

import (
	"fmt"
	"os"
	"sort"
	"time"

	openinstrument_proto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/variable"
)

func NowMs() uint64 {
	return uint64(time.Now().UnixNano() / 1000000)
}

func NewVariableFromString(textvar string) *variable.Variable {
	return variable.NewFromString(textvar)
}

func NewVariableFromProto(p *openinstrument_proto.StreamVariable) *variable.Variable {
	return variable.NewFromProto(p)
}

type Semaphore chan bool

// acquire n resources
func (s Semaphore) P(n int) {
	for i := 0; i < n; i++ {
		s <- true
	}
}

// release n resources
func (s Semaphore) V(n int) {
	for i := 0; i < n; i++ {
		<-s
	}
}

func (s Semaphore) Lock() {
	s.P(1)
}

func (s Semaphore) Unlock() {
	s.V(1)
}

/* signal-wait */

func (s Semaphore) Signal() {
	s.V(1)
}

func (s Semaphore) Wait(n int) {
	s.P(n)
}

// ValueStreamWriter returns a channel that appends values to the supplied ValueStream, performing run-length-encoding.
// No effort is made to ensure that the ValueStream contains sorted Values
func ValueStreamWriter(stream *openinstrument_proto.ValueStream) chan *openinstrument_proto.Value {
	c := make(chan *openinstrument_proto.Value)
	go func() {
		for value := range c {
			if len(stream.Value) > 0 {
				last := stream.Value[len(stream.Value)-1]
				if (last.GetStringValue() != "" && last.GetStringValue() == value.GetStringValue()) ||
					(last.GetDoubleValue() == value.GetDoubleValue()) {
					if value.GetEndTimestamp() > 0 {
						last.EndTimestamp = value.EndTimestamp
					} else {
						last.EndTimestamp = value.Timestamp
					}
				} else {
					stream.Value = append(stream.Value, value)
				}
			} else {
				stream.Value = append(stream.Value, value)
			}
		}
	}()
	return c
}

// ValueStreamReader returns a channel producing Values from the supplied ValueStream
func ValueStreamReader(stream *openinstrument_proto.ValueStream) chan *openinstrument_proto.Value {
	c := make(chan *openinstrument_proto.Value)
	go func() {
		for _, value := range stream.Value {
			c <- value
		}
		close(c)
	}()
	return c
}

// MergeValueStreams merges multiple ValueStreams, returning a channel producing sorted Values.
func MergeValueStreams(streams []*openinstrument_proto.ValueStream) chan *openinstrument_proto.Value {
	c := make(chan *openinstrument_proto.Value)
	n := len(streams)
	go func() {
		indexes := make([]int, n)
		for {
			var minTimestamp uint64
			var minStream *openinstrument_proto.ValueStream
			var minValue *openinstrument_proto.Value
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

func MergeStreamsBy(streams []*openinstrument_proto.ValueStream, by string) chan []*openinstrument_proto.ValueStream {
	c := make(chan []*openinstrument_proto.ValueStream)
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
				var output []*openinstrument_proto.ValueStream
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
					var output []*openinstrument_proto.ValueStream
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

type valueStreamChannelList struct {
	input    chan *openinstrument_proto.Value
	channels []chan *openinstrument_proto.Value
}

func (cl *valueStreamChannelList) Add(c chan *openinstrument_proto.Value) {
	cl.channels = append(cl.channels, c)
}

func (cl *valueStreamChannelList) Last() chan *openinstrument_proto.Value {
	return cl.channels[len(cl.channels)-1]
}

func ValueStreamChannelList(initial chan *openinstrument_proto.Value) *valueStreamChannelList {
	cl := new(valueStreamChannelList)
	cl.channels = make([]chan *openinstrument_proto.Value, 0)
	cl.channels = append(cl.channels, initial)
	return cl
}

func Readdirnames(directory string) ([]string, error) {
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
