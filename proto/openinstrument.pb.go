// Code generated by protoc-gen-go.
// source: openinstrument.proto
// DO NOT EDIT!

/*
Package openinstrument_proto is a generated protocol buffer package.

It is generated from these files:
	openinstrument.proto

It has these top-level messages:
	LogMessage
	Label
	StreamVariable
	StreamMutation
	StreamAggregation
	Value
	ValueStream
	Query
	GetRequest
	GetResponse
	AddRequest
	AddResponse
	ListRequest
	ListResponse
	StoreFileHeaderIndex
	StoreFileHeader
	RetentionPolicyItem
	RetentionPolicy
	StoreServerStatus
	StoreConfig
	Block
	LookupBlockRequest
	LookupBlockResponse
	SplitBlockRequest
	SplitBlockResponse
	JoinBlockRequest
	JoinBlockResponse
	CompactBlockRequest
	CompactBlockResponse
*/
package openinstrument_proto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type StreamVariable_ValueType int32

const (
	StreamVariable_UNKNOWN StreamVariable_ValueType = 0
	StreamVariable_GAUGE   StreamVariable_ValueType = 1
	StreamVariable_RATE    StreamVariable_ValueType = 2
)

var StreamVariable_ValueType_name = map[int32]string{
	0: "UNKNOWN",
	1: "GAUGE",
	2: "RATE",
}
var StreamVariable_ValueType_value = map[string]int32{
	"UNKNOWN": 0,
	"GAUGE":   1,
	"RATE":    2,
}

func (x StreamVariable_ValueType) String() string {
	return proto.EnumName(StreamVariable_ValueType_name, int32(x))
}

type StreamMutation_SampleType int32

const (
	// Do not modify the data at all
	StreamMutation_NONE StreamMutation_SampleType = 0
	// Average out values between each sampling point
	StreamMutation_AVERAGE StreamMutation_SampleType = 1
	// Get the maximum value between each sampling point
	StreamMutation_MAX StreamMutation_SampleType = 2
	// Get the minimum value between each sampling point
	StreamMutation_MIN StreamMutation_SampleType = 3
	// Get the non-negative rate between each sampling point
	StreamMutation_RATE StreamMutation_SampleType = 4
	// Get the (possibly) negative rate between each sampling point
	StreamMutation_RATE_SIGNED StreamMutation_SampleType = 5
	// Difference between each point
	StreamMutation_DELTA StreamMutation_SampleType = 6
	// Latest point
	StreamMutation_LATEST StreamMutation_SampleType = 7
	// Multiply all values by a constant
	StreamMutation_MULTIPLY StreamMutation_SampleType = 8
	// Add a constant to each value
	StreamMutation_ADD StreamMutation_SampleType = 9
)

var StreamMutation_SampleType_name = map[int32]string{
	0: "NONE",
	1: "AVERAGE",
	2: "MAX",
	3: "MIN",
	4: "RATE",
	5: "RATE_SIGNED",
	6: "DELTA",
	7: "LATEST",
	8: "MULTIPLY",
	9: "ADD",
}
var StreamMutation_SampleType_value = map[string]int32{
	"NONE":        0,
	"AVERAGE":     1,
	"MAX":         2,
	"MIN":         3,
	"RATE":        4,
	"RATE_SIGNED": 5,
	"DELTA":       6,
	"LATEST":      7,
	"MULTIPLY":    8,
	"ADD":         9,
}

func (x StreamMutation_SampleType) String() string {
	return proto.EnumName(StreamMutation_SampleType_name, int32(x))
}

type StreamAggregation_AggregateType int32

const (
	// Do not aggregate between streams
	StreamAggregation_NONE StreamAggregation_AggregateType = 0
	// Get an average of all values at a common time across streams
	StreamAggregation_AVERAGE StreamAggregation_AggregateType = 1
	// Get the maximum of all values at a common time across streams
	StreamAggregation_MAX StreamAggregation_AggregateType = 2
	// Get the minimum of all values at a common time across streams
	StreamAggregation_MIN StreamAggregation_AggregateType = 3
	// Get the median of all values at a common time across streams
	StreamAggregation_MEDIAN StreamAggregation_AggregateType = 4
	// Get the total of all values at a common time across streams
	StreamAggregation_SUM StreamAggregation_AggregateType = 5
	// Get the standard deviation of all values at a common time across streams
	StreamAggregation_STDDEV StreamAggregation_AggregateType = 6
	// Get a particular percentile (provided in percentile) of all values at a common time across streams
	StreamAggregation_PERCENTILE StreamAggregation_AggregateType = 7
)

var StreamAggregation_AggregateType_name = map[int32]string{
	0: "NONE",
	1: "AVERAGE",
	2: "MAX",
	3: "MIN",
	4: "MEDIAN",
	5: "SUM",
	6: "STDDEV",
	7: "PERCENTILE",
}
var StreamAggregation_AggregateType_value = map[string]int32{
	"NONE":       0,
	"AVERAGE":    1,
	"MAX":        2,
	"MIN":        3,
	"MEDIAN":     4,
	"SUM":        5,
	"STDDEV":     6,
	"PERCENTILE": 7,
}

func (x StreamAggregation_AggregateType) String() string {
	return proto.EnumName(StreamAggregation_AggregateType_name, int32(x))
}

type RetentionPolicyItem_Target int32

const (
	RetentionPolicyItem_UNKNOWN RetentionPolicyItem_Target = 0
	RetentionPolicyItem_KEEP    RetentionPolicyItem_Target = 1
	RetentionPolicyItem_DROP    RetentionPolicyItem_Target = 2
)

var RetentionPolicyItem_Target_name = map[int32]string{
	0: "UNKNOWN",
	1: "KEEP",
	2: "DROP",
}
var RetentionPolicyItem_Target_value = map[string]int32{
	"UNKNOWN": 0,
	"KEEP":    1,
	"DROP":    2,
}

func (x RetentionPolicyItem_Target) String() string {
	return proto.EnumName(RetentionPolicyItem_Target_name, int32(x))
}

type StoreServerStatus_State int32

const (
	StoreServerStatus_UNKNOWN  StoreServerStatus_State = 0
	StoreServerStatus_LOAD     StoreServerStatus_State = 1
	StoreServerStatus_RUN      StoreServerStatus_State = 2
	StoreServerStatus_DRAIN    StoreServerStatus_State = 3
	StoreServerStatus_READONLY StoreServerStatus_State = 4
	StoreServerStatus_SHUTDOWN StoreServerStatus_State = 5
)

var StoreServerStatus_State_name = map[int32]string{
	0: "UNKNOWN",
	1: "LOAD",
	2: "RUN",
	3: "DRAIN",
	4: "READONLY",
	5: "SHUTDOWN",
}
var StoreServerStatus_State_value = map[string]int32{
	"UNKNOWN":  0,
	"LOAD":     1,
	"RUN":      2,
	"DRAIN":    3,
	"READONLY": 4,
	"SHUTDOWN": 5,
}

func (x StoreServerStatus_State) String() string {
	return proto.EnumName(StoreServerStatus_State_name, int32(x))
}

type LogMessage struct {
	Timestamp uint64 `protobuf:"varint,1,opt,name=timestamp" json:"timestamp,omitempty"`
	Message   string `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
}

func (m *LogMessage) Reset()         { *m = LogMessage{} }
func (m *LogMessage) String() string { return proto.CompactTextString(m) }
func (*LogMessage) ProtoMessage()    {}

type Label struct {
	Label string `protobuf:"bytes,1,opt,name=label" json:"label,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
}

func (m *Label) Reset()         { *m = Label{} }
func (m *Label) String() string { return proto.CompactTextString(m) }
func (*Label) ProtoMessage()    {}

type StreamVariable struct {
	Name  string                   `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Label []*Label                 `protobuf:"bytes,2,rep,name=label" json:"label,omitempty"`
	Type  StreamVariable_ValueType `protobuf:"varint,3,opt,name=type,enum=openinstrument.proto.StreamVariable_ValueType" json:"type,omitempty"`
}

func (m *StreamVariable) Reset()         { *m = StreamVariable{} }
func (m *StreamVariable) String() string { return proto.CompactTextString(m) }
func (*StreamVariable) ProtoMessage()    {}

func (m *StreamVariable) GetLabel() []*Label {
	if m != nil {
		return m.Label
	}
	return nil
}

type StreamMutation struct {
	SampleType StreamMutation_SampleType `protobuf:"varint,1,opt,name=sample_type,enum=openinstrument.proto.StreamMutation_SampleType" json:"sample_type,omitempty"`
	// Stretch or compress the stream so that there is a value every <sample_frequency> ms.
	// Extra values between each sample will be aggregated according to <sample_type>.
	// Gaps in the stream less than <max_gap_interpolate> samples will be filled with interpolated values between the
	// closest real values.
	// Gaps in the stream of more than <max_gap_interpolate> samples will not be filled.
	SampleFrequency   uint32 `protobuf:"varint,2,opt,name=sample_frequency" json:"sample_frequency,omitempty"`
	MaxGapInterpolate uint32 `protobuf:"varint,3,opt,name=max_gap_interpolate" json:"max_gap_interpolate,omitempty"`
}

func (m *StreamMutation) Reset()         { *m = StreamMutation{} }
func (m *StreamMutation) String() string { return proto.CompactTextString(m) }
func (*StreamMutation) ProtoMessage()    {}

type StreamAggregation struct {
	Type StreamAggregation_AggregateType `protobuf:"varint,1,opt,name=type,enum=openinstrument.proto.StreamAggregation_AggregateType" json:"type,omitempty"`
	// Labels to aggregate by on the input streams. If no labels are specified, aggregation will be done on the variable
	// only.
	Label []string `protobuf:"bytes,2,rep,name=label" json:"label,omitempty"`
	// Points will be aggregated if they are less than sample_interval ms apart. Default is 30 seconds.
	SampleInterval uint32 `protobuf:"varint,3,opt,name=sample_interval" json:"sample_interval,omitempty"`
}

func (m *StreamAggregation) Reset()         { *m = StreamAggregation{} }
func (m *StreamAggregation) String() string { return proto.CompactTextString(m) }
func (*StreamAggregation) ProtoMessage()    {}

type Value struct {
	// Milliseconds since epoch
	Timestamp   uint64  `protobuf:"varint,1,opt,name=timestamp" json:"timestamp,omitempty"`
	DoubleValue float64 `protobuf:"fixed64,2,opt,name=double_value" json:"double_value,omitempty"`
	StringValue string  `protobuf:"bytes,3,opt,name=string_value" json:"string_value,omitempty"`
	// Used for run-length encoding
	EndTimestamp uint64 `protobuf:"varint,4,opt,name=end_timestamp" json:"end_timestamp,omitempty"`
}

func (m *Value) Reset()         { *m = Value{} }
func (m *Value) String() string { return proto.CompactTextString(m) }
func (*Value) ProtoMessage()    {}

type ValueStream struct {
	Variable *StreamVariable `protobuf:"bytes,2,opt,name=variable" json:"variable,omitempty"`
	Value    []*Value        `protobuf:"bytes,4,rep,name=value" json:"value,omitempty"`
	// An optional set of mutations that have been applied to this stream.
	// If the stream contains the raw data, this should be empty.
	Mutation []*StreamMutation `protobuf:"bytes,5,rep,name=mutation" json:"mutation,omitempty"`
	// Used internally to save cycles regenerating the name
	VariableName string `protobuf:"bytes,6,opt,name=variable_name" json:"variable_name,omitempty"`
}

func (m *ValueStream) Reset()         { *m = ValueStream{} }
func (m *ValueStream) String() string { return proto.CompactTextString(m) }
func (*ValueStream) ProtoMessage()    {}

func (m *ValueStream) GetVariable() *StreamVariable {
	if m != nil {
		return m.Variable
	}
	return nil
}

func (m *ValueStream) GetValue() []*Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *ValueStream) GetMutation() []*StreamMutation {
	if m != nil {
		return m.Mutation
	}
	return nil
}

type Query struct {
	// At least one of these must be set for each subquery
	Subquery []*Query          `protobuf:"bytes,1,rep,name=subquery" json:"subquery,omitempty"`
	Variable []*StreamVariable `protobuf:"bytes,2,rep,name=variable" json:"variable,omitempty"`
	Constant []float64         `protobuf:"fixed64,3,rep,name=constant" json:"constant,omitempty"`
	// Milliseconds since epoch
	MinTimestamp uint64               `protobuf:"varint,4,opt,name=min_timestamp" json:"min_timestamp,omitempty"`
	MaxTimestamp uint64               `protobuf:"varint,5,opt,name=max_timestamp" json:"max_timestamp,omitempty"`
	Mutation     []*StreamMutation    `protobuf:"bytes,6,rep,name=mutation" json:"mutation,omitempty"`
	Aggregation  []*StreamAggregation `protobuf:"bytes,7,rep,name=aggregation" json:"aggregation,omitempty"`
}

func (m *Query) Reset()         { *m = Query{} }
func (m *Query) String() string { return proto.CompactTextString(m) }
func (*Query) ProtoMessage()    {}

func (m *Query) GetSubquery() []*Query {
	if m != nil {
		return m.Subquery
	}
	return nil
}

func (m *Query) GetVariable() []*StreamVariable {
	if m != nil {
		return m.Variable
	}
	return nil
}

func (m *Query) GetMutation() []*StreamMutation {
	if m != nil {
		return m.Mutation
	}
	return nil
}

func (m *Query) GetAggregation() []*StreamAggregation {
	if m != nil {
		return m.Aggregation
	}
	return nil
}

type GetRequest struct {
	Variable *StreamVariable `protobuf:"bytes,9,opt,name=variable" json:"variable,omitempty"`
	// Milliseconds since epoch
	MinTimestamp uint64 `protobuf:"varint,2,opt,name=min_timestamp" json:"min_timestamp,omitempty"`
	MaxTimestamp uint64 `protobuf:"varint,3,opt,name=max_timestamp" json:"max_timestamp,omitempty"`
	// If mutations or aggregations are supplied, the GetResponse will have a separate stream for each change requested.
	Mutation    []*StreamMutation    `protobuf:"bytes,6,rep,name=mutation" json:"mutation,omitempty"`
	Aggregation []*StreamAggregation `protobuf:"bytes,7,rep,name=aggregation" json:"aggregation,omitempty"`
	// Limit to how many variables can be returned in a standard GetRequest. This can be overridden if required.
	MaxVariables uint32 `protobuf:"varint,8,opt,name=max_variables" json:"max_variables,omitempty"`
	// This request has been forwarded by another store server and should not be forwarded again.
	// This shouldn't happen but is here as a failsafe.
	Forwarded bool `protobuf:"varint,10,opt,name=forwarded" json:"forwarded,omitempty"`
	// Limit the number of values that can be returned for each variable.
	MaxValues uint32 `protobuf:"varint,11,opt,name=max_values" json:"max_values,omitempty"`
}

func (m *GetRequest) Reset()         { *m = GetRequest{} }
func (m *GetRequest) String() string { return proto.CompactTextString(m) }
func (*GetRequest) ProtoMessage()    {}

func (m *GetRequest) GetVariable() *StreamVariable {
	if m != nil {
		return m.Variable
	}
	return nil
}

func (m *GetRequest) GetMutation() []*StreamMutation {
	if m != nil {
		return m.Mutation
	}
	return nil
}

func (m *GetRequest) GetAggregation() []*StreamAggregation {
	if m != nil {
		return m.Aggregation
	}
	return nil
}

type GetResponse struct {
	Success      bool   `protobuf:"varint,1,opt,name=success" json:"success,omitempty"`
	Errormessage string `protobuf:"bytes,2,opt,name=errormessage" json:"errormessage,omitempty"`
	// Contains a separate stream for every mutation requested in GetRequest, in the same order. If no mutations are
	// requested, the response will contain a single stream of the raw data.
	Stream []*ValueStream `protobuf:"bytes,3,rep,name=stream" json:"stream,omitempty"`
	Timer  []*LogMessage  `protobuf:"bytes,4,rep,name=timer" json:"timer,omitempty"`
}

func (m *GetResponse) Reset()         { *m = GetResponse{} }
func (m *GetResponse) String() string { return proto.CompactTextString(m) }
func (*GetResponse) ProtoMessage()    {}

func (m *GetResponse) GetStream() []*ValueStream {
	if m != nil {
		return m.Stream
	}
	return nil
}

func (m *GetResponse) GetTimer() []*LogMessage {
	if m != nil {
		return m.Timer
	}
	return nil
}

type AddRequest struct {
	Stream []*ValueStream `protobuf:"bytes,1,rep,name=stream" json:"stream,omitempty"`
	// This request has been forwarded by another store server and should not be forwarded again.
	// This shouldn't happen but is here as a failsafe.
	Forwarded bool `protobuf:"varint,2,opt,name=forwarded" json:"forwarded,omitempty"`
}

func (m *AddRequest) Reset()         { *m = AddRequest{} }
func (m *AddRequest) String() string { return proto.CompactTextString(m) }
func (*AddRequest) ProtoMessage()    {}

func (m *AddRequest) GetStream() []*ValueStream {
	if m != nil {
		return m.Stream
	}
	return nil
}

type AddResponse struct {
	Success      bool          `protobuf:"varint,1,opt,name=success" json:"success,omitempty"`
	Errormessage string        `protobuf:"bytes,2,opt,name=errormessage" json:"errormessage,omitempty"`
	Timer        []*LogMessage `protobuf:"bytes,3,rep,name=timer" json:"timer,omitempty"`
}

func (m *AddResponse) Reset()         { *m = AddResponse{} }
func (m *AddResponse) String() string { return proto.CompactTextString(m) }
func (*AddResponse) ProtoMessage()    {}

func (m *AddResponse) GetTimer() []*LogMessage {
	if m != nil {
		return m.Timer
	}
	return nil
}

type ListRequest struct {
	Prefix *StreamVariable `protobuf:"bytes,3,opt,name=prefix" json:"prefix,omitempty"`
	// Limit to how many variables can be returned in a standard ListRequest. This can be overridden if required.
	MaxVariables uint32 `protobuf:"varint,2,opt,name=max_variables" json:"max_variables,omitempty"`
	// Maximum age of variable.
	// This controls how far back the search will go for variables that were used in the past but not currently.
	MaxAge uint64 `protobuf:"varint,4,opt,name=max_age" json:"max_age,omitempty"`
}

func (m *ListRequest) Reset()         { *m = ListRequest{} }
func (m *ListRequest) String() string { return proto.CompactTextString(m) }
func (*ListRequest) ProtoMessage()    {}

func (m *ListRequest) GetPrefix() *StreamVariable {
	if m != nil {
		return m.Prefix
	}
	return nil
}

type ListResponse struct {
	Success      bool   `protobuf:"varint,1,opt,name=success" json:"success,omitempty"`
	Errormessage string `protobuf:"bytes,2,opt,name=errormessage" json:"errormessage,omitempty"`
	// Although this uses the ValueStream for returning a list of variables, the expectation is that the value field is
	// empty.
	Stream   []*ValueStream    `protobuf:"bytes,3,rep,name=stream" json:"stream,omitempty"`
	Variable []*StreamVariable `protobuf:"bytes,4,rep,name=variable" json:"variable,omitempty"`
	Timer    []*LogMessage     `protobuf:"bytes,5,rep,name=timer" json:"timer,omitempty"`
}

func (m *ListResponse) Reset()         { *m = ListResponse{} }
func (m *ListResponse) String() string { return proto.CompactTextString(m) }
func (*ListResponse) ProtoMessage()    {}

func (m *ListResponse) GetStream() []*ValueStream {
	if m != nil {
		return m.Stream
	}
	return nil
}

func (m *ListResponse) GetVariable() []*StreamVariable {
	if m != nil {
		return m.Variable
	}
	return nil
}

func (m *ListResponse) GetTimer() []*LogMessage {
	if m != nil {
		return m.Timer
	}
	return nil
}

type StoreFileHeaderIndex struct {
	Variable     *StreamVariable `protobuf:"bytes,1,opt,name=variable" json:"variable,omitempty"`
	Offset       uint64          `protobuf:"fixed64,2,opt,name=offset" json:"offset,omitempty"`
	NumValues    uint32          `protobuf:"fixed32,3,opt,name=num_values" json:"num_values,omitempty"`
	MinTimestamp uint64          `protobuf:"fixed64,4,opt,name=min_timestamp" json:"min_timestamp,omitempty"`
	MaxTimestamp uint64          `protobuf:"fixed64,5,opt,name=max_timestamp" json:"max_timestamp,omitempty"`
}

func (m *StoreFileHeaderIndex) Reset()         { *m = StoreFileHeaderIndex{} }
func (m *StoreFileHeaderIndex) String() string { return proto.CompactTextString(m) }
func (*StoreFileHeaderIndex) ProtoMessage()    {}

func (m *StoreFileHeaderIndex) GetVariable() *StreamVariable {
	if m != nil {
		return m.Variable
	}
	return nil
}

type StoreFileHeader struct {
	Version        uint32                  `protobuf:"varint,6,opt,name=version" json:"version,omitempty"`
	StartTimestamp uint64                  `protobuf:"varint,1,opt,name=start_timestamp" json:"start_timestamp,omitempty"`
	EndTimestamp   uint64                  `protobuf:"varint,2,opt,name=end_timestamp" json:"end_timestamp,omitempty"`
	Variable       []*StreamVariable       `protobuf:"bytes,4,rep,name=variable" json:"variable,omitempty"`
	Index          []*StoreFileHeaderIndex `protobuf:"bytes,5,rep,name=index" json:"index,omitempty"`
	EndKey         string                  `protobuf:"bytes,7,opt,name=end_key" json:"end_key,omitempty"`
}

func (m *StoreFileHeader) Reset()         { *m = StoreFileHeader{} }
func (m *StoreFileHeader) String() string { return proto.CompactTextString(m) }
func (*StoreFileHeader) ProtoMessage()    {}

func (m *StoreFileHeader) GetVariable() []*StreamVariable {
	if m != nil {
		return m.Variable
	}
	return nil
}

func (m *StoreFileHeader) GetIndex() []*StoreFileHeaderIndex {
	if m != nil {
		return m.Index
	}
	return nil
}

type RetentionPolicyItem struct {
	Variable []*StreamVariable          `protobuf:"bytes,1,rep,name=variable" json:"variable,omitempty"`
	Comment  []string                   `protobuf:"bytes,2,rep,name=comment" json:"comment,omitempty"`
	Policy   RetentionPolicyItem_Target `protobuf:"varint,3,opt,name=policy,enum=openinstrument.proto.RetentionPolicyItem_Target" json:"policy,omitempty"`
	Mutation []*StreamMutation          `protobuf:"bytes,4,rep,name=mutation" json:"mutation,omitempty"`
	MinAge   uint64                     `protobuf:"varint,5,opt,name=min_age" json:"min_age,omitempty"`
	MaxAge   uint64                     `protobuf:"varint,6,opt,name=max_age" json:"max_age,omitempty"`
}

func (m *RetentionPolicyItem) Reset()         { *m = RetentionPolicyItem{} }
func (m *RetentionPolicyItem) String() string { return proto.CompactTextString(m) }
func (*RetentionPolicyItem) ProtoMessage()    {}

func (m *RetentionPolicyItem) GetVariable() []*StreamVariable {
	if m != nil {
		return m.Variable
	}
	return nil
}

func (m *RetentionPolicyItem) GetMutation() []*StreamMutation {
	if m != nil {
		return m.Mutation
	}
	return nil
}

type RetentionPolicy struct {
	Policy   []*RetentionPolicyItem `protobuf:"bytes,1,rep,name=policy" json:"policy,omitempty"`
	Interval uint32                 `protobuf:"varint,2,opt,name=interval" json:"interval,omitempty"`
}

func (m *RetentionPolicy) Reset()         { *m = RetentionPolicy{} }
func (m *RetentionPolicy) String() string { return proto.CompactTextString(m) }
func (*RetentionPolicy) ProtoMessage()    {}

func (m *RetentionPolicy) GetPolicy() []*RetentionPolicyItem {
	if m != nil {
		return m.Policy
	}
	return nil
}

type StoreServerStatus struct {
	Address     string                  `protobuf:"bytes,1,opt,name=address" json:"address,omitempty"`
	State       StoreServerStatus_State `protobuf:"varint,2,opt,name=state,enum=openinstrument.proto.StoreServerStatus_State" json:"state,omitempty"`
	LastUpdated uint64                  `protobuf:"varint,3,opt,name=last_updated" json:"last_updated,omitempty"`
	// Optional name for use in the hash ring. This should never change once the
	// server has been added but can be used to replace a server with another one
	// on a different host.
	// If this is not set, the address is used.
	Name string `protobuf:"bytes,5,opt,name=name" json:"name,omitempty"`
	// Desired size in bytes for indexed datastore files
	TargetIndexedFileSize uint64 `protobuf:"varint,4,opt,name=target_indexed_file_size" json:"target_indexed_file_size,omitempty"`
}

func (m *StoreServerStatus) Reset()         { *m = StoreServerStatus{} }
func (m *StoreServerStatus) String() string { return proto.CompactTextString(m) }
func (*StoreServerStatus) ProtoMessage()    {}

type StoreConfig struct {
	Server          []*StoreServerStatus `protobuf:"bytes,1,rep,name=server" json:"server,omitempty"`
	RetentionPolicy *RetentionPolicy     `protobuf:"bytes,2,opt,name=retention_policy" json:"retention_policy,omitempty"`
}

func (m *StoreConfig) Reset()         { *m = StoreConfig{} }
func (m *StoreConfig) String() string { return proto.CompactTextString(m) }
func (*StoreConfig) ProtoMessage()    {}

func (m *StoreConfig) GetServer() []*StoreServerStatus {
	if m != nil {
		return m.Server
	}
	return nil
}

func (m *StoreConfig) GetRetentionPolicy() *RetentionPolicy {
	if m != nil {
		return m.RetentionPolicy
	}
	return nil
}

type Block struct {
	Id              string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	EndKey          string `protobuf:"bytes,2,opt,name=end_key" json:"end_key,omitempty"`
	IndexedStreams  uint32 `protobuf:"varint,3,opt,name=indexed_streams" json:"indexed_streams,omitempty"`
	IndexedValues   uint32 `protobuf:"varint,4,opt,name=indexed_values" json:"indexed_values,omitempty"`
	LoggedStreams   uint32 `protobuf:"varint,5,opt,name=logged_streams" json:"logged_streams,omitempty"`
	LoggedValues    uint32 `protobuf:"varint,6,opt,name=logged_values" json:"logged_values,omitempty"`
	UnloggedStreams uint32 `protobuf:"varint,7,opt,name=unlogged_streams" json:"unlogged_streams,omitempty"`
	UnloggedValues  uint32 `protobuf:"varint,8,opt,name=unlogged_values" json:"unlogged_values,omitempty"`
	IsCompacting    bool   `protobuf:"varint,9,opt,name=is_compacting" json:"is_compacting,omitempty"`
	CompactDuration string `protobuf:"bytes,10,opt,name=compact_duration" json:"compact_duration,omitempty"`
}

func (m *Block) Reset()         { *m = Block{} }
func (m *Block) String() string { return proto.CompactTextString(m) }
func (*Block) ProtoMessage()    {}

type LookupBlockRequest struct {
	Block *Block `protobuf:"bytes,1,opt,name=block" json:"block,omitempty"`
}

func (m *LookupBlockRequest) Reset()         { *m = LookupBlockRequest{} }
func (m *LookupBlockRequest) String() string { return proto.CompactTextString(m) }
func (*LookupBlockRequest) ProtoMessage()    {}

func (m *LookupBlockRequest) GetBlock() *Block {
	if m != nil {
		return m.Block
	}
	return nil
}

type LookupBlockResponse struct {
	Block *Block `protobuf:"bytes,1,opt,name=block" json:"block,omitempty"`
}

func (m *LookupBlockResponse) Reset()         { *m = LookupBlockResponse{} }
func (m *LookupBlockResponse) String() string { return proto.CompactTextString(m) }
func (*LookupBlockResponse) ProtoMessage()    {}

func (m *LookupBlockResponse) GetBlock() *Block {
	if m != nil {
		return m.Block
	}
	return nil
}

type SplitBlockRequest struct {
	Block *Block `protobuf:"bytes,1,opt,name=block" json:"block,omitempty"`
}

func (m *SplitBlockRequest) Reset()         { *m = SplitBlockRequest{} }
func (m *SplitBlockRequest) String() string { return proto.CompactTextString(m) }
func (*SplitBlockRequest) ProtoMessage()    {}

func (m *SplitBlockRequest) GetBlock() *Block {
	if m != nil {
		return m.Block
	}
	return nil
}

type SplitBlockResponse struct {
	Block []*Block `protobuf:"bytes,1,rep,name=block" json:"block,omitempty"`
}

func (m *SplitBlockResponse) Reset()         { *m = SplitBlockResponse{} }
func (m *SplitBlockResponse) String() string { return proto.CompactTextString(m) }
func (*SplitBlockResponse) ProtoMessage()    {}

func (m *SplitBlockResponse) GetBlock() []*Block {
	if m != nil {
		return m.Block
	}
	return nil
}

type JoinBlockRequest struct {
	// This the last block in the sequence. The block with the previous end key
	// will be joined with this one, so the output block will have the same end
	// key, but a different ID.
	Block *Block `protobuf:"bytes,1,opt,name=block" json:"block,omitempty"`
}

func (m *JoinBlockRequest) Reset()         { *m = JoinBlockRequest{} }
func (m *JoinBlockRequest) String() string { return proto.CompactTextString(m) }
func (*JoinBlockRequest) ProtoMessage()    {}

func (m *JoinBlockRequest) GetBlock() *Block {
	if m != nil {
		return m.Block
	}
	return nil
}

type JoinBlockResponse struct {
	Block *Block `protobuf:"bytes,1,opt,name=block" json:"block,omitempty"`
}

func (m *JoinBlockResponse) Reset()         { *m = JoinBlockResponse{} }
func (m *JoinBlockResponse) String() string { return proto.CompactTextString(m) }
func (*JoinBlockResponse) ProtoMessage()    {}

func (m *JoinBlockResponse) GetBlock() *Block {
	if m != nil {
		return m.Block
	}
	return nil
}

type CompactBlockRequest struct {
	Block *Block `protobuf:"bytes,1,opt,name=block" json:"block,omitempty"`
}

func (m *CompactBlockRequest) Reset()         { *m = CompactBlockRequest{} }
func (m *CompactBlockRequest) String() string { return proto.CompactTextString(m) }
func (*CompactBlockRequest) ProtoMessage()    {}

func (m *CompactBlockRequest) GetBlock() *Block {
	if m != nil {
		return m.Block
	}
	return nil
}

type CompactBlockResponse struct {
	Block *Block `protobuf:"bytes,1,opt,name=block" json:"block,omitempty"`
}

func (m *CompactBlockResponse) Reset()         { *m = CompactBlockResponse{} }
func (m *CompactBlockResponse) String() string { return proto.CompactTextString(m) }
func (*CompactBlockResponse) ProtoMessage()    {}

func (m *CompactBlockResponse) GetBlock() *Block {
	if m != nil {
		return m.Block
	}
	return nil
}

func init() {
	proto.RegisterEnum("openinstrument.proto.StreamVariable_ValueType", StreamVariable_ValueType_name, StreamVariable_ValueType_value)
	proto.RegisterEnum("openinstrument.proto.StreamMutation_SampleType", StreamMutation_SampleType_name, StreamMutation_SampleType_value)
	proto.RegisterEnum("openinstrument.proto.StreamAggregation_AggregateType", StreamAggregation_AggregateType_name, StreamAggregation_AggregateType_value)
	proto.RegisterEnum("openinstrument.proto.RetentionPolicyItem_Target", RetentionPolicyItem_Target_name, RetentionPolicyItem_Target_value)
	proto.RegisterEnum("openinstrument.proto.StoreServerStatus_State", StoreServerStatus_State_name, StoreServerStatus_State_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Client API for Store service

type StoreClient interface {
	List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error)
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (Store_GetClient, error)
	Add(ctx context.Context, opts ...grpc.CallOption) (Store_AddClient, error)
	LookupBlock(ctx context.Context, in *LookupBlockRequest, opts ...grpc.CallOption) (*LookupBlockResponse, error)
	SplitBlock(ctx context.Context, in *SplitBlockRequest, opts ...grpc.CallOption) (*SplitBlockResponse, error)
	JoinBlock(ctx context.Context, in *JoinBlockRequest, opts ...grpc.CallOption) (*JoinBlockResponse, error)
	CompactBlock(ctx context.Context, in *CompactBlockRequest, opts ...grpc.CallOption) (*CompactBlockResponse, error)
}

type storeClient struct {
	cc *grpc.ClientConn
}

func NewStoreClient(cc *grpc.ClientConn) StoreClient {
	return &storeClient{cc}
}

func (c *storeClient) List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error) {
	out := new(ListResponse)
	err := grpc.Invoke(ctx, "/openinstrument.proto.Store/List", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storeClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (Store_GetClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Store_serviceDesc.Streams[0], c.cc, "/openinstrument.proto.Store/Get", opts...)
	if err != nil {
		return nil, err
	}
	x := &storeGetClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Store_GetClient interface {
	Recv() (*GetResponse, error)
	grpc.ClientStream
}

type storeGetClient struct {
	grpc.ClientStream
}

func (x *storeGetClient) Recv() (*GetResponse, error) {
	m := new(GetResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *storeClient) Add(ctx context.Context, opts ...grpc.CallOption) (Store_AddClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Store_serviceDesc.Streams[1], c.cc, "/openinstrument.proto.Store/Add", opts...)
	if err != nil {
		return nil, err
	}
	x := &storeAddClient{stream}
	return x, nil
}

type Store_AddClient interface {
	Send(*AddRequest) error
	Recv() (*AddResponse, error)
	grpc.ClientStream
}

type storeAddClient struct {
	grpc.ClientStream
}

func (x *storeAddClient) Send(m *AddRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *storeAddClient) Recv() (*AddResponse, error) {
	m := new(AddResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *storeClient) LookupBlock(ctx context.Context, in *LookupBlockRequest, opts ...grpc.CallOption) (*LookupBlockResponse, error) {
	out := new(LookupBlockResponse)
	err := grpc.Invoke(ctx, "/openinstrument.proto.Store/LookupBlock", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storeClient) SplitBlock(ctx context.Context, in *SplitBlockRequest, opts ...grpc.CallOption) (*SplitBlockResponse, error) {
	out := new(SplitBlockResponse)
	err := grpc.Invoke(ctx, "/openinstrument.proto.Store/SplitBlock", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storeClient) JoinBlock(ctx context.Context, in *JoinBlockRequest, opts ...grpc.CallOption) (*JoinBlockResponse, error) {
	out := new(JoinBlockResponse)
	err := grpc.Invoke(ctx, "/openinstrument.proto.Store/JoinBlock", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storeClient) CompactBlock(ctx context.Context, in *CompactBlockRequest, opts ...grpc.CallOption) (*CompactBlockResponse, error) {
	out := new(CompactBlockResponse)
	err := grpc.Invoke(ctx, "/openinstrument.proto.Store/CompactBlock", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Store service

type StoreServer interface {
	List(context.Context, *ListRequest) (*ListResponse, error)
	Get(*GetRequest, Store_GetServer) error
	Add(Store_AddServer) error
	LookupBlock(context.Context, *LookupBlockRequest) (*LookupBlockResponse, error)
	SplitBlock(context.Context, *SplitBlockRequest) (*SplitBlockResponse, error)
	JoinBlock(context.Context, *JoinBlockRequest) (*JoinBlockResponse, error)
	CompactBlock(context.Context, *CompactBlockRequest) (*CompactBlockResponse, error)
}

func RegisterStoreServer(s *grpc.Server, srv StoreServer) {
	s.RegisterService(&_Store_serviceDesc, srv)
}

func _Store_List_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(ListRequest)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(StoreServer).List(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _Store_Get_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StoreServer).Get(m, &storeGetServer{stream})
}

type Store_GetServer interface {
	Send(*GetResponse) error
	grpc.ServerStream
}

type storeGetServer struct {
	grpc.ServerStream
}

func (x *storeGetServer) Send(m *GetResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Store_Add_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StoreServer).Add(&storeAddServer{stream})
}

type Store_AddServer interface {
	Send(*AddResponse) error
	Recv() (*AddRequest, error)
	grpc.ServerStream
}

type storeAddServer struct {
	grpc.ServerStream
}

func (x *storeAddServer) Send(m *AddResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *storeAddServer) Recv() (*AddRequest, error) {
	m := new(AddRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Store_LookupBlock_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(LookupBlockRequest)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(StoreServer).LookupBlock(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _Store_SplitBlock_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(SplitBlockRequest)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(StoreServer).SplitBlock(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _Store_JoinBlock_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(JoinBlockRequest)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(StoreServer).JoinBlock(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func _Store_CompactBlock_Handler(srv interface{}, ctx context.Context, codec grpc.Codec, buf []byte) (interface{}, error) {
	in := new(CompactBlockRequest)
	if err := codec.Unmarshal(buf, in); err != nil {
		return nil, err
	}
	out, err := srv.(StoreServer).CompactBlock(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _Store_serviceDesc = grpc.ServiceDesc{
	ServiceName: "openinstrument.proto.Store",
	HandlerType: (*StoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "List",
			Handler:    _Store_List_Handler,
		},
		{
			MethodName: "LookupBlock",
			Handler:    _Store_LookupBlock_Handler,
		},
		{
			MethodName: "SplitBlock",
			Handler:    _Store_SplitBlock_Handler,
		},
		{
			MethodName: "JoinBlock",
			Handler:    _Store_JoinBlock_Handler,
		},
		{
			MethodName: "CompactBlock",
			Handler:    _Store_CompactBlock_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Get",
			Handler:       _Store_Get_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Add",
			Handler:       _Store_Add_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
}
