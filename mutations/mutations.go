package mutations

import (
	"math"

	oproto "github.com/dparrish/openinstrument/proto"
)

func Mean(input *oproto.ValueStream) *oproto.ValueStream {
	output := &oproto.ValueStream{}
	var lastValue float64
	for _, v := range input.Value {
		output.Value = append(output.Value, &oproto.Value{
			Timestamp:   v.Timestamp,
			DoubleValue: (v.DoubleValue + lastValue) / 2,
		})
		lastValue = v.DoubleValue
	}
	return output
}

func MovingAverage(window uint64, input *oproto.ValueStream) *oproto.ValueStream {
	output := &oproto.ValueStream{}
	var items []*oproto.Value
	outItems := func() {
		var total float64
		for _, i := range items {
			total += i.DoubleValue
		}
		output.Value = append(output.Value, &oproto.Value{
			Timestamp:   items[len(items)-1].Timestamp,
			DoubleValue: total / float64(len(items)),
		})
	}
	for _, v := range input.Value {
		if len(items) == 0 || items[0].Timestamp+window >= v.Timestamp {
			// Still within window
			items = append(items, v)
			continue
		}
		// Outside current window
		outItems()
		// Discard items before the new window
		for len(items) > 0 && items[0].Timestamp < v.Timestamp-window {
			items = items[1:]
		}
		items = append(items, v)
	}
	if len(items) > 0 {
		// Last moving average
		outItems()
	}
	return output
}

func SignedRateStream(input *oproto.ValueStream) *oproto.ValueStream {
	var lastValue float64
	var lastTimestamp uint64
	output := &oproto.ValueStream{}
	for _, v := range input.Value {
		if lastTimestamp == 0 {
			lastValue = v.DoubleValue
			lastTimestamp = v.Timestamp
			continue
		}
		rate := (v.DoubleValue - lastValue) / float64(v.Timestamp-lastTimestamp)
		output.Value = append(output.Value, &oproto.Value{
			Timestamp:   v.Timestamp,
			DoubleValue: rate,
		})
		lastValue = v.DoubleValue
		lastTimestamp = v.Timestamp
	}
	return output
}

func SignedRate(input *oproto.ValueStream) *oproto.ValueStream {
	output := &oproto.ValueStream{}
	var lastValue float64
	var lastTimestamp uint64
	first := true
	for _, v := range input.Value {
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
		output.Value = append(output.Value, &oproto.Value{
			Timestamp:   v.Timestamp,
			DoubleValue: rate,
		})
		lastValue = v.DoubleValue
		lastTimestamp = v.Timestamp
	}
	return output
}

func Rate(input *oproto.ValueStream) *oproto.ValueStream {
	output := &oproto.ValueStream{}
	newoutput := SignedRate(input)
	for _, v := range newoutput.Value {
		if v.DoubleValue >= 0 {
			output.Value = append(output.Value, v)
		}
	}
	return output
}

func First(duration uint64, input *oproto.ValueStream) *oproto.ValueStream {
	output := &oproto.ValueStream{}
	var nextTimestamp uint64
	for _, v := range input.Value {
		if len(output.Value) == 0 || v.Timestamp >= nextTimestamp {
			output.Value = append(output.Value, v)
			nextTimestamp += duration
		}
	}
	return output
}

func Last(duration uint64, input *oproto.ValueStream) *oproto.ValueStream {
	output := &oproto.ValueStream{}
	var (
		nextTimestamp uint64
		value         *oproto.Value
	)
	for _, v := range input.Value {
		if value == nil {
			nextTimestamp = v.Timestamp + duration
		}
		if value != nil && v.Timestamp > nextTimestamp {
			output.Value = append(output.Value, value)
			nextTimestamp += duration
		}
		value = v
	}
	if value != nil {
		output.Value = append(output.Value, value)
	}
	return output
}

func Multiply(param float64, input *oproto.ValueStream) *oproto.ValueStream {
	output := &oproto.ValueStream{}
	for _, v := range input.Value {
		newv := &*v
		newv.DoubleValue *= param
		output.Value = append(output.Value, newv)
	}
	return output
}

func Add(param float64, input *oproto.ValueStream) *oproto.ValueStream {
	output := &oproto.ValueStream{}
	for _, v := range input.Value {
		newv := &*v
		newv.DoubleValue += param
		output.Value = append(output.Value, newv)
	}
	return output
}

func Root(param float64, input *oproto.ValueStream) *oproto.ValueStream {
	output := &oproto.ValueStream{}
	for _, v := range input.Value {
		newv := &*v
		if param == 2 {
			newv.DoubleValue = math.Sqrt(v.DoubleValue)
		} else if param == 3 {
			newv.DoubleValue = math.Cbrt(v.DoubleValue)
		} else {
			newv.DoubleValue = math.Pow(v.DoubleValue, 1.0/param)
		}
		output.Value = append(output.Value, newv)
	}
	return output
}

func Power(param float64, input *oproto.ValueStream) *oproto.ValueStream {
	output := &oproto.ValueStream{}
	for _, v := range input.Value {
		newv := &*v
		newv.DoubleValue = math.Pow(v.DoubleValue, param)
		output.Value = append(output.Value, newv)
	}
	return output
}

// Interpolate takes a counter value sampled at variable intervals and aligns each value with a <duration> ms boundary.
// The initial value is not output unless it appears exactly on a ms boundary.
// Runs of two or more missing values will be filled by NaN
func Interpolate(duration uint64, input *oproto.ValueStream) *oproto.ValueStream {
	output := &oproto.ValueStream{}
	var previousValue *oproto.Value
	var timestamp uint64
	for _, v := range input.Value {
		if previousValue == nil {
			if v.Timestamp%duration == 0 {
				// Value is exactly on a timestamp
				output.Value = append(output.Value, v)
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
				output.Value = append(output.Value, &oproto.Value{
					Timestamp:   timestamp,
					DoubleValue: newValue,
				})
			}
			if previousValue.Timestamp < v.Timestamp {
				if v.Timestamp%duration == 0 {
					// Value is exactly on a timestamp
					output.Value = append(output.Value, v)
					previousValue = v
					timestamp += duration
				} else {
					previousValue = v
				}
			}
		}
	}
	return output
}

func Min(duration uint64, input *oproto.ValueStream) *oproto.ValueStream {
	output := &oproto.ValueStream{}
	var (
		min           float64
		lastTimestamp uint64
		minTimestamp  uint64
		first         bool = true
	)
	for _, v := range input.Value {
		if v.StringValue != "" {
			continue
		}
		if first {
			lastTimestamp = v.Timestamp
			minTimestamp = v.Timestamp
			min = v.DoubleValue
			first = false
			continue
		}
		if v.Timestamp >= lastTimestamp+duration {
			output.Value = append(output.Value, &oproto.Value{
				Timestamp:   minTimestamp,
				DoubleValue: min,
			})
			lastTimestamp = v.Timestamp
			minTimestamp = v.Timestamp
			min = v.DoubleValue
			first = true
		}
		if v.DoubleValue < min {
			min = v.DoubleValue
			minTimestamp = v.Timestamp
		}
		first = false
	}
	return output
}

func Max(duration uint64, input *oproto.ValueStream) *oproto.ValueStream {
	output := &oproto.ValueStream{}
	var (
		max           float64
		lastTimestamp uint64
		maxTimestamp  uint64
		first         bool = true
	)
	for _, v := range input.Value {
		if first {
			lastTimestamp = v.Timestamp
			maxTimestamp = v.Timestamp
			max = v.DoubleValue
			first = false
			continue
		}
		if v.DoubleValue >= max {
			max = v.DoubleValue
			maxTimestamp = v.Timestamp
		}
		if v.Timestamp >= lastTimestamp+duration {
			output.Value = append(output.Value, &oproto.Value{
				Timestamp:   maxTimestamp,
				DoubleValue: max,
			})
			lastTimestamp = v.Timestamp
			first = true
		}
		first = false
	}
	return output
}
