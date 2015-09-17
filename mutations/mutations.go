package mutations

import oproto "github.com/dparrish/openinstrument/proto"

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
				firstTimestamp = v.Timestamp
				first = false
			}

			if v.Timestamp-firstTimestamp > duration {
				mean := sum / float64(count)
				output <- &oproto.Value{
					Timestamp:   lastTimestamp,
					DoubleValue: mean,
				}
				firstTimestamp = v.Timestamp
				count = 0
				sum = 0
			}
			sum += v.DoubleValue
			count++
			lastTimestamp = v.Timestamp
		}
		if count > 0 {
			mean := sum / float64(count)
			output <- &oproto.Value{
				Timestamp:   lastTimestamp,
				DoubleValue: mean,
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
			if v.StringValue != "" {
				continue
			}
			if first {
				lastValue = v.DoubleValue
				lastTimestamp = v.Timestamp
				first = false
				continue
			}
			rate := (v.DoubleValue - lastValue) / float64(v.Timestamp-lastTimestamp)
			output <- &oproto.Value{
				Timestamp:   v.Timestamp,
				DoubleValue: rate,
			}
			lastValue = v.DoubleValue
			lastTimestamp = v.Timestamp
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
			if v.DoubleValue >= 0 {
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
				if v.Timestamp%duration == 0 {
					// Value is exactly on a timestamp
					output <- v
					previousValue = v
					timestamp = v.Timestamp + duration
				} else {
					previousValue = v
					previousValue.Timestamp = v.Timestamp / duration * duration
					timestamp = previousValue.Timestamp + duration
				}
				continue
			}
			if v.Timestamp >= timestamp {
				// Fill in any missing values before this one
				rate := float64((v.DoubleValue - previousValue.DoubleValue))
				for ; timestamp <= v.Timestamp; timestamp += duration {
					pct := float64(timestamp-previousValue.Timestamp) / float64(v.Timestamp-previousValue.Timestamp)
					newValue := previousValue.DoubleValue + (rate * pct)
					output <- &oproto.Value{
						Timestamp:   timestamp,
						DoubleValue: newValue,
					}
				}
				if previousValue.Timestamp < v.Timestamp {
					if v.Timestamp%duration == 0 {
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
				lastTimestamp = v.Timestamp
				min = v.DoubleValue
				first = false
				continue
			}
			if v.Timestamp >= lastTimestamp+duration {
				output <- &oproto.Value{
					Timestamp:   v.Timestamp,
					DoubleValue: min,
				}
				lastTimestamp = v.Timestamp
				first = true
			}
			if first || v.DoubleValue < min {
				min = v.DoubleValue
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
				lastTimestamp = v.Timestamp
				max = v.DoubleValue
				first = false
				continue
			}
			if v.Timestamp >= lastTimestamp+duration {
				output <- &oproto.Value{
					Timestamp:   v.Timestamp,
					DoubleValue: max,
				}
				lastTimestamp = v.Timestamp
				first = true
			}
			if first || v.DoubleValue > max {
				max = v.DoubleValue
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
