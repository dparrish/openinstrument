package valuestream

import (
	"sort"

	oproto "github.com/dparrish/openinstrument/proto"
)

type By func(p1, p2 *oproto.ValueStream) bool

func (by By) Sort(values []*oproto.ValueStream) {
	sfs := &valuesSorter{
		values: values,
		by:     by,
	}
	sort.Sort(sfs)
}

type valuesSorter struct {
	values []*oproto.ValueStream
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
