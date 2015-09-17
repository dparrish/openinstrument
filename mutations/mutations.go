package mutations

import (
	"github.com/golang/protobuf/proto"
	oproto "github.com/dparrish/openinstrument/proto"
)

type MutateFunc func(duration uint64, input chan *oproto.Value) <-chan *oproto.Value

func Copy(duration uint64, input <-chan *oproto.Value) <-chan *oproto.Value {
	output := make(chan *oproto.Value, 1000)
	go func() {
		for v := range input {
			output <- v
		}
		close(output)
	}()
	return output
}

func Mean(duration uint64, input <-chan *oproto.Value) <-chan *oproto.Value {
	output := make(chan *oproto.Value, 1000)
	go func() {
		var (
			sum            float64
			count          uint64
			firstTimestamp uint64
			lastTimestamp  uint64
		)
		first := true
		for v := range input {
			if duration == 0 {
				output <- v
				continue
			}

			if first {
				firstTimestamp = v.GetTimestamp()
				first = false
			}

			if v.GetTimestamp()-firstTimestamp > duration {
				mean := sum / float64(count)
				output <- &oproto.Value{
					Timestamp:   proto.Uint64(lastTimestamp),
					DoubleValue: proto.Float64(mean),
				}
				firstTimestamp = v.GetTimestamp()
				count = 0
				sum = 0
			}
			sum += v.GetDoubleValue()
			count++
			lastTimestamp = v.GetTimestamp()
		}
		if count > 0 {
			mean := sum / float64(count)
			output <- &oproto.Value{
				Timestamp:   proto.Uint64(lastTimestamp),
				DoubleValue: proto.Float64(mean),
			}
		}
		close(output)
	}()
	return output
}

func SignedRate(duration uint64, input <-chan *oproto.Value) <-chan *oproto.Value {
	output := make(chan *oproto.Value, 1000)
	go func() {
		var lastValue float64
		var lastTimestamp uint64
		first := true
		for v := range input {
			if v.DoubleValue == nil {
				continue
			}
			if first {
				lastValue = v.GetDoubleValue()
				lastTimestamp = v.GetTimestamp()
				first = false
				continue
			}
			rate := (v.GetDoubleValue() - lastValue) / float64(v.GetTimestamp()-lastTimestamp)
			output <- &oproto.Value{
				Timestamp:   v.Timestamp,
				DoubleValue: proto.Float64(rate),
			}
			lastValue = v.GetDoubleValue()
			lastTimestamp = v.GetTimestamp()
		}
		close(output)
	}()
	return output
}

func Rate(duration uint64, input <-chan *oproto.Value) <-chan *oproto.Value {
	output := make(chan *oproto.Value, 1000)
	go func() {
		newoutput := SignedRate(duration, input)
		for v := range newoutput {
			if v.GetDoubleValue() >= 0 {
				output <- v
			}
		}
		close(output)
	}()
	return output
}

// Interpolate takes a counter value sampled at variable intervals and aligns each value with a <duration> ms boundary.
// The initial value is not output unless it appears exactly on a ms boundary.
// Runs of two or more missing values will be filled by NaN
func Interpolate(duration uint64, input <-chan *oproto.Value) <-chan *oproto.Value {
	output := make(chan *oproto.Value, 1000)
	go func() {
		var previousValue *oproto.Value
		var timestamp uint64
		for v := range input {
			if previousValue == nil {
				if v.GetTimestamp()%duration == 0 {
					// Value is exactly on a timestamp
					output <- v
					previousValue = v
					timestamp = v.GetTimestamp() + duration
				} else {
					previousValue = v
					previousValue.Timestamp = proto.Uint64(v.GetTimestamp() / duration * duration)
					timestamp = previousValue.GetTimestamp() + duration
				}
				continue
			}
			if v.GetTimestamp() >= timestamp {
				// Fill in any missing values before this one
				rate := float64((v.GetDoubleValue() - previousValue.GetDoubleValue()))
				for ; timestamp <= v.GetTimestamp(); timestamp += duration {
					pct := float64(timestamp-previousValue.GetTimestamp()) / float64(v.GetTimestamp()-previousValue.GetTimestamp())
					newValue := previousValue.GetDoubleValue() + (rate * pct)
					output <- &oproto.Value{
						Timestamp:   proto.Uint64(timestamp),
						DoubleValue: proto.Float64(newValue),
					}
				}
				if previousValue.GetTimestamp() < v.GetTimestamp() {
					if v.GetTimestamp()%duration == 0 {
						// Value is exactly on a timestamp
						output <- v
						previousValue = v
						timestamp += duration
					} else {
						previousValue = v
					}
				}
			}
		}
		close(output)
	}()
	return output
}

func Min(duration uint64, input <-chan *oproto.Value) <-chan *oproto.Value {
	output := make(chan *oproto.Value, 1000)
	go func() {
		var min float64
		var lastTimestamp uint64
		first := true
		for v := range input {
			if first {
				lastTimestamp = v.GetTimestamp()
				min = v.GetDoubleValue()
				first = false
				continue
			}
			if v.GetTimestamp() >= lastTimestamp+duration {
				output <- &oproto.Value{
					Timestamp:   proto.Uint64(v.GetTimestamp()),
					DoubleValue: proto.Float64(min),
				}
				lastTimestamp = v.GetTimestamp()
				first = true
			}
			if first || v.GetDoubleValue() < min {
				min = v.GetDoubleValue()
			}
			first = false
		}
		close(output)
	}()
	return output
}

func Max(duration uint64, input <-chan *oproto.Value) <-chan *oproto.Value {
	output := make(chan *oproto.Value, 1000)
	go func() {
		var max float64
		var lastTimestamp uint64
		first := true
		for v := range input {
			if first {
				lastTimestamp = v.GetTimestamp()
				max = v.GetDoubleValue()
				first = false
				continue
			}
			if v.GetTimestamp() >= lastTimestamp+duration {
				output <- &oproto.Value{
					Timestamp:   proto.Uint64(v.GetTimestamp()),
					DoubleValue: proto.Float64(max),
				}
				lastTimestamp = v.GetTimestamp()
				first = true
			}
			if first || v.GetDoubleValue() > max {
				max = v.GetDoubleValue()
			}
			first = false
		}
		close(output)
	}()
	return output
}

// MutateValues performs a mutation on an input channel of Values.
// The mutation will be performed on values every <duration> ms apart, with the output timestamp set to the last
// timestamp of the sequence.
func MutateValues(duration uint64, input <-chan *oproto.Value, f MutateFunc) <-chan *oproto.Value {
	output := make(chan *oproto.Value, 1000)
	go func() {
		mutateInput := make(chan *oproto.Value, 1000)
		c := f(duration, mutateInput)
		for value := range input {
			mutateInput <- value
		}
		close(mutateInput)
		for value := range c {
			output <- value
		}
	}()
	return output
}
