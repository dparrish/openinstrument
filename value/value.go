package value

import (
	"sort"

	oproto "github.com/dparrish/openinstrument/proto"
)

type By func(p1, p2 *oproto.Value) bool

func (by By) Sort(values []*oproto.Value) {
	sfs := &valuesSorter{
		values: values,
		by:     by,
	}
	sort.Sort(sfs)
}

type valuesSorter struct {
	values []*oproto.Value
	by     By
}

func (vs *valuesSorter) Len() int {
	return len(vs.values)
}

func (vs *valuesSorter) Swap(i, j int) {
	vs.values[i], vs.values[j] = vs.values[j], vs.values[i]
}

func (vs *valuesSorter) Less(i, j int) bool {
	return vs.by(vs.values[i], vs.values[j])
}

func NewDouble(timestamp uint64, val float64) *oproto.Value {
	return &oproto.Value{
		Value:     &oproto.Value_Double{Double: val},
		Timestamp: timestamp,
	}
}

func NewString(timestamp uint64, val string) *oproto.Value {
	return &oproto.Value{
		Value:     &oproto.Value_String_{String_: val},
		Timestamp: timestamp,
	}
}
