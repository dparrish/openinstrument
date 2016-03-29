package aggregations

import (
	"fmt"
	"math"
	"strings"

	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/value"
	"github.com/dparrish/openinstrument/variable"
)

func CartesianProduct(input [][]string) [][]string {
	output := [][]string{}
	if len(input) == 0 {
		return output
	}
	for _, v := range input[0] {
		if len(input) > 1 {
			for _, c := range CartesianProduct(input[1:]) {
				o := []string{v}
				o = append(o, c...)
				output = append(output, o)
			}
		} else {
			output = append(output, []string{v})
		}
	}
	return output
}

func CartesianProductToVariable(input [][]string) []*variable.Variable {
	output := []*variable.Variable{}
	for _, l := range input {
		v := &variable.Variable{
			Labels: make(map[string]string),
		}
		for _, m := range l {
			parts := strings.SplitN(m, "=", 2)
			v.Labels[parts[0]] = parts[1]
		}
		output = append(output, v)
	}
	return output
}

func splitStreamsByLabels(by []string, input []*oproto.ValueStream) map[string]map[string][]*oproto.ValueStream {
	//log.Println("Getting streams by", by)
	output := map[string]map[string][]*oproto.ValueStream{}
	for _, b := range by {
		output[b] = map[string][]*oproto.ValueStream{}
		for _, stream := range input {
			//log.Println("Checking", variable.ProtoToString(stream.Variable))
			v, ok := stream.Variable.Label[b]
			if !ok {
				v = ""
			}
			//log.Printf("  %s = %s", b, v)
			output[b][v] = append(output[b][v], stream)
		}
	}
	//log.Println(output)
	return output
}

func By(by []string, input []*oproto.ValueStream) []*variable.Variable {
	if len(by) == 0 {
		return []*variable.Variable{variable.NewFromString("*")}
	}
	//log.Println("Getting streams by", by)
	bys := splitStreamsByLabels(by, input)
	byMatches := [][]string{}
	for b, x := range bys {
		byMatch := []string{}
		for v, _ := range x {
			byMatch = append(byMatch, fmt.Sprintf("%s=%s", b, v))
		}
		byMatches = append(byMatches, byMatch)
	}
	//log.Println("byMatches:", byMatches)
	cp := CartesianProduct(byMatches)
	//log.Println("Cartesian product:", cp)
	return CartesianProductToVariable(cp)
}

// Get the mean of all values at a common time across streams
func Mean(by []string, input []*oproto.ValueStream) []*oproto.ValueStream {
	output := []*oproto.ValueStream{}
	for _, inputVar := range By(by, input) {
		streams := []*oproto.ValueStream{}
		for _, stream := range input {
			matchVar := variable.NewFromProto(stream.Variable)
			if matchVar.Match(inputVar) {
				streams = append(streams, stream)
			}
		}
		stream := &oproto.ValueStream{Variable: inputVar.AsProto()}
		iPos := make([]int, len(streams))
		for {
			values := []float64{}
			timestamps := []uint64{}
			for i := 0; i < len(streams); i++ {
				if iPos[i] >= len(streams[i].Value) {
					continue
				}
				if streams[i] != nil {
					values = append(values, streams[i].Value[iPos[i]].GetDouble())
					timestamps = append(timestamps, streams[i].Value[iPos[i]].Timestamp)
				}
				iPos[i]++
			}
			if len(values) == 0 {
				break
			}
			var total float64
			for _, i := range values {
				total += i
			}
			var tsTotal uint64
			for _, i := range timestamps {
				tsTotal += i
			}
			stream.Value = append(stream.Value, value.NewDouble(tsTotal/uint64(len(timestamps)), total/float64(len(values))))
		}
		if len(stream.Value) > 0 {
			output = append(output, stream)
		}
	}
	return output
}

// Get the minimum of all values at a common time across streams
func Min(by []string, input []*oproto.ValueStream) []*oproto.ValueStream {
	output := []*oproto.ValueStream{{Variable: input[0].Variable}}
	iPos := make([]int, len(input))
	for {
		var min *float64
		var ts *uint64
		for i := 0; i < len(input); i++ {
			if iPos[i] >= len(input[i].Value) {
				continue
			}
			if input[i] != nil {
				if min == nil || input[i].Value[iPos[i]].GetDouble() < *min {
					min = &input[i].Value[iPos[i]].Value.(*oproto.Value_Double).Double
					ts = &input[i].Value[iPos[i]].Timestamp
				}
			}
			iPos[i]++
		}
		if min == nil {
			break
		}
		output[0].Value = append(output[0].Value, value.NewDouble(*ts, *min))
	}
	return output
}

// Get the maximum of all values at a common time across streams
func Max(by []string, input []*oproto.ValueStream) []*oproto.ValueStream {
	output := []*oproto.ValueStream{{Variable: input[0].Variable}}
	iPos := make([]int, len(input))
	for {
		var max *float64
		var ts *uint64
		for i := 0; i < len(input); i++ {
			if iPos[i] >= len(input[i].Value) {
				continue
			}
			if input[i] != nil {
				if max == nil || input[i].Value[iPos[i]].GetDouble() > *max {
					max = &input[i].Value[iPos[i]].Value.(*oproto.Value_Double).Double
					ts = &input[i].Value[iPos[i]].Timestamp
				}
			}
			iPos[i]++
		}
		if max == nil {
			break
		}
		output[0].Value = append(output[0].Value, value.NewDouble(*ts, *max))
	}
	return output
}

// Get the median of all values at a common time across streams
func Median(by []string, input []*oproto.ValueStream) []*oproto.ValueStream {
	output := []*oproto.ValueStream{{Variable: input[0].Variable}}
	return output
}

// Get the total of all values at a common time across streams
func Sum(by []string, input []*oproto.ValueStream) []*oproto.ValueStream {
	output := []*oproto.ValueStream{{Variable: input[0].Variable}}
	iPos := make([]int, len(input))
	for {
		values := []float64{}
		timestamps := []uint64{}
		for i := 0; i < len(input); i++ {
			if iPos[i] >= len(input[i].Value) {
				continue
			}
			if input[i] != nil {
				values = append(values, input[i].Value[iPos[i]].GetDouble())
				timestamps = append(timestamps, input[i].Value[iPos[i]].Timestamp)
			}
			iPos[i]++
		}
		if len(values) == 0 {
			break
		}
		var total float64
		for _, i := range values {
			total += i
		}
		var tsTotal uint64
		for _, i := range timestamps {
			tsTotal += i
		}
		output[0].Value = append(output[0].Value, value.NewDouble(tsTotal/uint64(len(timestamps)), total))
	}
	return output
}

// Get the standard deviation of all values at a common time across streams
func StdDev(by []string, input []*oproto.ValueStream) []*oproto.ValueStream {
	output := []*oproto.ValueStream{{Variable: input[0].Variable}}
	iPos := make([]int, len(input))
	for {
		values := []float64{}
		timestamps := []uint64{}
		for i := 0; i < len(input); i++ {
			if iPos[i] >= len(input[i].Value) {
				continue
			}
			if input[i] != nil {
				values = append(values, input[i].Value[iPos[i]].GetDouble())
				timestamps = append(timestamps, input[i].Value[iPos[i]].Timestamp)
			}
			iPos[i]++
		}
		if len(values) == 0 {
			break
		}
		var total float64
		for _, i := range values {
			total += i
		}
		var variances float64
		for _, i := range values {
			variances += math.Pow(((total / float64(len(values))) - i), 2)
		}

		var tsTotal uint64
		for _, i := range timestamps {
			tsTotal += i
		}
		output[0].Value = append(output[0].Value, value.NewDouble(tsTotal/uint64(len(timestamps)),
			math.Sqrt(variances/float64(len(values)))))
	}
	return output
}

// Get a particular percentile (provided in percentile) of all values at a common time across streams
func Percentile(by []string, percentile float64, input []*oproto.ValueStream) []*oproto.ValueStream {
	output := []*oproto.ValueStream{{Variable: input[0].Variable}}
	return output
}
