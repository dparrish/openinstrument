package rle

import oproto "github.com/dparrish/openinstrument/proto"

func Encode(input <-chan *oproto.Value) <-chan *oproto.Value {
	output := make(chan *oproto.Value)
	go func() {
		var last *oproto.Value
		for value := range input {
			if last == nil {
				last = value
				continue
			}

			if last.StringValue != "" && value.StringValue != "" {
				//log.Printf("Last at %d is %s, this is %s", last.Timestamp, last.StringValue, value.StringValue)
				if last.StringValue == value.StringValue {
					if value.EndTimestamp > value.Timestamp {
						last.EndTimestamp = value.EndTimestamp
					} else {
						last.EndTimestamp = value.Timestamp
					}
					continue
				}
			} else {
				//log.Printf("Last is at %d %f, this is %f", last.Timestamp, last.DoubleValue, value.DoubleValue)
				if last.DoubleValue == value.DoubleValue {
					if value.EndTimestamp > value.Timestamp {
						last.EndTimestamp = value.EndTimestamp
					} else {
						last.EndTimestamp = value.Timestamp
					}
					continue
				}
			}

			if last.EndTimestamp == 0 && last.Timestamp != 0 {
				last.EndTimestamp = last.Timestamp
			}
			output <- last
			last = value
		}

		if last != nil {
			if last.EndTimestamp == 0 && last.Timestamp != 0 {
				last.EndTimestamp = last.Timestamp
			}
			output <- last
		}
		close(output)
	}()
	return output
}
