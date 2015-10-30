package rle

import oproto "github.com/dparrish/openinstrument/proto"

func Encode(input *oproto.ValueStream) *oproto.ValueStream {
	output := &oproto.ValueStream{Variable: input.Variable}
	var last *oproto.Value
	for _, value := range input.Value {
		if last == nil {
			last = value
			continue
		}

		if last.StringValue != "" && value.StringValue != "" {
			if last.StringValue == value.StringValue {
				if value.EndTimestamp > value.Timestamp {
					last.EndTimestamp = value.EndTimestamp
				} else {
					last.EndTimestamp = value.Timestamp
				}
				continue
			}
		} else {
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
		output.Value = append(output.Value, last)
		last = value
	}

	if last != nil {
		if last.EndTimestamp == 0 && last.Timestamp != 0 {
			last.EndTimestamp = last.Timestamp
		}
		output.Value = append(output.Value, last)
	}
	return output
}
