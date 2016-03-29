package rle

import oproto "github.com/dparrish/openinstrument/proto"

func getEndTimestamp(v *oproto.Value) uint64 {
	if v.EndTimestamp > v.Timestamp {
		return v.EndTimestamp
	}
	return v.Timestamp
}

func Encode(input *oproto.ValueStream) *oproto.ValueStream {
	output := &oproto.ValueStream{Variable: input.Variable}
	var last *oproto.Value

	for _, v := range input.Value {
		if last == nil {
			last = v
			continue
		}

		switch v.Value.(type) {
		case *oproto.Value_String_:
			if last.GetString_() == v.GetString_() {
				last.EndTimestamp = getEndTimestamp(v)
				continue
			}
		case *oproto.Value_Double:
			if last.GetDouble() == v.GetDouble() {
				last.EndTimestamp = getEndTimestamp(v)
				continue
			}
		case nil:
			// No value field
			continue
		}

		last.EndTimestamp = getEndTimestamp(last)
		output.Value = append(output.Value, last)
		last = v
	}

	if last != nil {
		last.EndTimestamp = getEndTimestamp(last)
		output.Value = append(output.Value, last)
	}
	return output
}
