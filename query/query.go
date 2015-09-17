package query

import (
	"time"

	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/variable"
	"github.com/golang/protobuf/proto"
)

type Query struct {
	q *oproto.Query
}

func (query *Query) AsProto() *oproto.Query {
	return query.q
}

func (query *Query) String() string {
	return proto.MarshalTextString(query.q)
}

func New() *Query {
	query := &Query{
		q: &oproto.Query{},
	}
	return query
}

func NewFromString(str string) *Query {
	query := &Query{
		q: &oproto.Query{},
	}
	proto.UnmarshalText(str, query.q)
	return query
}

func (query *Query) AddVariable(v *variable.Variable) *Query {
	if query.q.Variable == nil {
		query.q.Variable = make([]*oproto.StreamVariable, 0)
	}
	query.q.Variable = append(query.q.Variable, v.AsProto())
	return query
}

func (query *Query) AddVariableString(v string) *Query {
	return query.AddVariable(variable.NewFromString(v))
}

func (query *Query) AddSubquery(q *oproto.Query) *Query {
	if query.q.Subquery == nil {
		query.q.Subquery = make([]*oproto.Query, 0)
	}
	query.q.Subquery = append(query.q.Subquery, q)
	return query
}

func (query *Query) AddConstant(constant float64) *Query {
	if query.q.Constant == nil {
		query.q.Constant = make([]float64, 0)
	}
	query.q.Constant = append(query.q.Constant, constant)
	return query
}

func (query *Query) SetMinTimestamp(minTimestamp time.Time) *Query {
	query.q.MinTimestamp = uint64(minTimestamp.UnixNano() / 1000000)
	return query
}

func (query *Query) SetMaxTimestamp(maxTimestamp time.Time) *Query {
	query.q.MaxTimestamp = uint64(maxTimestamp.UnixNano() / 1000000)
	return query
}

func (query *Query) AddAggregation(t oproto.StreamAggregation_AggregateType, labels []string) *Query {
	if query.q.Aggregation == nil {
		query.q.Aggregation = make([]*oproto.StreamAggregation, 0)
	}
	agg := &oproto.StreamAggregation{
		Type: t,
	}
	agg.Label = labels
	query.q.Aggregation = append(query.q.Aggregation, agg)
	return query
}

func (query *Query) AddMutation(t oproto.StreamMutation_SampleType, sampleFrequency, maxGapInterpolate uint32) *Query {
	if query.q.Mutation == nil {
		query.q.Mutation = make([]*oproto.StreamMutation, 0)
	}
	agg := &oproto.StreamMutation{
		SampleType: t,
	}
	if sampleFrequency > 0 {
		agg.SampleFrequency = sampleFrequency
	}
	if maxGapInterpolate > 0 {
		agg.MaxGapInterpolate = maxGapInterpolate
	}
	query.q.Mutation = append(query.q.Mutation, agg)
	return query
}
