package rle

import (
	"code.google.com/p/goprotobuf/proto"
	oproto "github.com/dparrish/openinstrument/proto"
)

func Encode(input <-chan *oproto.Value, outputBufSize int) <-chan *oproto.Value {
	output := make(chan *oproto.Value, outputBufSize)
	go func() {
		var last *oproto.Value
		for value := range input {
			if last == nil {
				last = value
				continue
			}

			if last.StringValue != nil && value.StringValue != nil {
				//log.Printf("Last at %d is %s, this is %s", last.GetTimestamp(), last.GetStringValue(), value.GetStringValue())
				if last.GetStringValue() == value.GetStringValue() {
					if value.GetEndTimestamp() > value.GetTimestamp() {
						last.EndTimestamp = proto.Uint64(value.GetEndTimestamp())
					} else {
						last.EndTimestamp = proto.Uint64(value.GetTimestamp())
					}
					continue
				}
			}

			if last.DoubleValue != nil && value.DoubleValue != nil {
				//log.Printf("Last is at %d %f, this is %f", last.GetTimestamp(), last.GetDoubleValue(), value.GetDoubleValue())
				if last.GetDoubleValue() == value.GetDoubleValue() {
					if value.GetEndTimestamp() > value.GetTimestamp() {
						last.EndTimestamp = proto.Uint64(value.GetEndTimestamp())
					} else {
						last.EndTimestamp = proto.Uint64(value.GetTimestamp())
					}
					continue
				}
			}

			if last.EndTimestamp == nil {
				last.EndTimestamp = proto.Uint64(last.GetTimestamp())
			}
			output <- last
			last = value
		}

		if last != nil {
			if last.EndTimestamp == nil {
				last.EndTimestamp = proto.Uint64(last.GetTimestamp())
			}
			output <- last
		}
		close(output)
	}()
	return output
}
