package query

import (
  //"github.com/dparrish/openinstrument/query/ast"
  //"github.com/dparrish/openinstrument/lexer"
  //"github.com/dparrish/openinstrument/parser"
  oproto "github.com/dparrish/openinstrument/proto"
  "code.google.com/p/goprotobuf/proto"
  "github.com/dparrish/openinstrument/variable"
  "time"
)

type Query struct {
  q *oproto.Query
}

func (this *Query) AsProto() *oproto.Query {
  return this.q
}

func (this *Query) String() string {
  return proto.MarshalTextString(this.q)
}

func New() *Query {
  this := &Query{
    q : &oproto.Query{},
  }
  return this
}

func NewFromString(str string) *Query {
  this := &Query{
    q : &oproto.Query{},
  }
  proto.UnmarshalText(str, this.q)
  return this
}

func (this *Query) AddVariable(v *variable.Variable) *Query {
  if this.q.Variable == nil {
    this.q.Variable = make([]*oproto.StreamVariable, 0)
  }
  this.q.Variable = append(this.q.Variable, v.AsProto())
  return this
}

func (this *Query) AddVariableString(v string) *Query {
  return this.AddVariable(variable.NewFromString(v))
}

func (this *Query) AddSubquery(q *oproto.Query) *Query {
  if this.q.Subquery == nil {
    this.q.Subquery = make([]*oproto.Query, 0)
  }
  this.q.Subquery = append(this.q.Subquery, q)
  return this
}

func (this *Query) AddConstant(constant float64) *Query {
  if this.q.Constant == nil {
    this.q.Constant = make([]float64, 0)
  }
  this.q.Constant = append(this.q.Constant, constant)
  return this
}

func (this *Query) SetMinTimestamp(min_timestamp time.Time) *Query {
  this.q.MinTimestamp = proto.Uint64(uint64(min_timestamp.UnixNano() / 1000000))
  return this
}

func (this *Query) SetMaxTimestamp(max_timestamp time.Time) *Query {
  this.q.MaxTimestamp = proto.Uint64(uint64(max_timestamp.UnixNano() / 1000000))
  return this
}

func (this *Query) AddAggregation(t oproto.StreamAggregation_AggregateType, labels []string) *Query {
  if this.q.Aggregation == nil {
    this.q.Aggregation = make([]*oproto.StreamAggregation, 0)
  }
  agg := &oproto.StreamAggregation{
    Type: oproto.StreamAggregation_AggregateType.Enum(t),
  }
  agg.Label = labels
  this.q.Aggregation = append(this.q.Aggregation, agg)
  return this
}

func (this *Query) AddMutation(t oproto.StreamMutation_SampleType, sample_frequency, max_gap_interpolate uint32) *Query {
  if this.q.Mutation == nil {
    this.q.Mutation = make([]*oproto.StreamMutation, 0)
  }
  agg := &oproto.StreamMutation{
    SampleType: oproto.StreamMutation_SampleType.Enum(t),
  }
  if sample_frequency > 0 {
    agg.SampleFrequency = proto.Uint32(sample_frequency)
  }
  if max_gap_interpolate > 0 {
    agg.MaxGapInterpolate = proto.Uint32(max_gap_interpolate)
  }
  this.q.Mutation = append(this.q.Mutation, agg)
  return this
}


