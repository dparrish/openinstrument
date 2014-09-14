package exported

import (
  "code.google.com/p/goprotobuf/proto"
  "github.com/dparrish/openinstrument/variable"
  oproto "github.com/dparrish/openinstrument/proto"
  "time"
  "log"
)

type VariableExporter struct {
  shutdown bool
  completed chan bool
}

func NewVariableExporter(address string, interval int32) *VariableExporter {
  this := new(VariableExporter)
  this.completed = make(chan bool, 1)
  go func(this *VariableExporter) {
    tick := time.Tick(time.Duration(interval) * time.Millisecond)
    for !this.shutdown {
      <-tick
      // Flush
      log.Printf("Flushing exported variables")
    }
    this.completed <- true
  }(this)
  return this
}

func (this *VariableExporter) Shutdown() {
  this.shutdown = true
  <-this.completed
}

type Exportable interface {
  Export() chan *oproto.ValueStream;
  Variable() *variable.Variable;
  SetVariable(v *variable.Variable);
}

//////////// ExportedInteger ///////////////////////////////
type ExportedInteger struct {
  value uint64
  v *variable.Variable
}

func NewExportedInteger(v string) *ExportedInteger {
  return &ExportedInteger{
    value: 0,
    v: variable.NewFromString(v),
  }
}

func (this *ExportedInteger) Variable() *variable.Variable {
  return this.v
}

func (this *ExportedInteger) SetVariable(v *variable.Variable) {
  this.v = v
}

func (this *ExportedInteger) Export() chan *oproto.ValueStream {
  c := make(chan *oproto.ValueStream)
  go func() {
    stream := new(oproto.ValueStream)
    stream.Variable = this.v.AsProto()
    stream.Value = append(stream.Value, &oproto.Value{
      Timestamp: proto.Uint64(uint64(time.Now().UnixNano() / 1000000)),
      DoubleValue: proto.Float64(float64(this.value)),
    })
    c <- stream
    close(c)
  }()
  return c
}

func (this *ExportedInteger) Get() uint64 {
  return this.value
}

func (this *ExportedInteger) Set(value uint64) {
  this.value = value
}

func (this *ExportedInteger) Add(value uint64) {
  this.value += value
}

func (this *ExportedInteger) Sub(value uint64) {
  this.value -= value
}

//////////// ExportedFloat ///////////////////////////////
type ExportedFloat struct {
  value float64
  v *variable.Variable
}

func NewExportedFloat(v string) *ExportedFloat {
  return &ExportedFloat{
    value: 0,
    v: variable.NewFromString(v),
  }
}

func (this *ExportedFloat) Variable() *variable.Variable {
  return this.v
}

func (this *ExportedFloat) SetVariable(v *variable.Variable) {
  this.v = v
}

func (this *ExportedFloat) Export() chan *oproto.ValueStream {
  c := make(chan *oproto.ValueStream)
  go func() {
    stream := new(oproto.ValueStream)
    stream.Variable = this.v.AsProto()
    stream.Value = append(stream.Value, &oproto.Value{
      Timestamp: proto.Uint64(uint64(time.Now().UnixNano() / 1000000)),
      DoubleValue: proto.Float64(this.value),
    })
    c <- stream
    close(c)
  }()
  return c
}

func (this *ExportedFloat) Get() float64 {
  return this.value
}

func (this *ExportedFloat) Set(value float64) {
  this.value = value
}

func (this *ExportedFloat) Add(value float64) {
  this.value += value
}

func (this *ExportedFloat) Sub(value float64) {
  this.value -= value
}

//////////// ExportedRatio ///////////////////////////////
type ExportedRatio struct {
  success *ExportedInteger
  failure *ExportedInteger
  total *ExportedInteger
}

func NewExportedRatio(v string) *ExportedRatio {
  success_var := variable.NewFromString(v)
  success_var.Variable += "-success"
  failure_var := variable.NewFromString(v)
  failure_var.Variable += "-failure"
  total_var := variable.NewFromString(v)
  total_var.Variable += "-total"
  return &ExportedRatio{
    success: NewExportedInteger(success_var.String()),
    failure: NewExportedInteger(failure_var.String()),
    total: NewExportedInteger(total_var.String()),
  }
}

func (this *ExportedRatio) Export() chan *oproto.ValueStream {
  c := make(chan *oproto.ValueStream)
  go func() {
    for stream := range this.success.Export() {
      c <- stream
    }
    for stream := range this.failure.Export() {
      c <- stream
    }
    for stream := range this.total.Export() {
      c <- stream
    }
    close(c)
  }()
  return c
}

func (this *ExportedRatio) Success() {
  this.success.Add(1)
  this.total.Add(1)
}

func (this *ExportedRatio) Failure() {
  this.failure.Add(1)
  this.total.Add(1)
}

//////////// ExportedAverage ///////////////////////////////
type ExportedAverage struct {
  overallSum *ExportedFloat
  totalCount *ExportedInteger
}

func NewExportedAverage(v string) *ExportedAverage {
  sum_var := variable.NewFromString(v)
  sum_var.Variable += "-total-count"
  total_var := variable.NewFromString(v)
  sum_var.Variable += "-overall-sum"
  return &ExportedAverage{
    overallSum: NewExportedFloat(sum_var.String()),
    totalCount: NewExportedInteger(total_var.String()),
  }
}

func (this *ExportedAverage) Export() chan *oproto.ValueStream {
  c := make(chan *oproto.ValueStream)
  go func() {
    for stream := range this.overallSum.Export() {
      c <- stream
    }
    for stream := range this.totalCount.Export() {
      c <- stream
    }
    close(c)
  }()
  return c
}

func (this *ExportedAverage) Update(sum float64, count uint64) {
  this.overallSum.Add(sum)
  this.totalCount.Add(count)
}


//////////// ExportedTimer ///////////////////////////////
type ExportedTimer struct {
  average *ExportedAverage
  start_time time.Time
}

func NewExportedTimer(v string) *ExportedTimer {
  return &ExportedTimer{
    average: NewExportedAverage(v),
  }
}

func (this *ExportedTimer) Export() chan *oproto.ValueStream {
  return this.average.Export()
}

func (this *ExportedTimer) Start() {
  this.start_time = time.Now()
}

func (this *ExportedTimer) Stop() {
  duration := time.Since(this.start_time)
  this.average.Update(float64(duration.Nanoseconds() / 1000000), 1)
}

//////////// ExportedString ///////////////////////////////
type ExportedString struct {
  value string
  v *variable.Variable
}

func NewExportedString(v string) *ExportedString {
  return &ExportedString{
    value: "",
    v: variable.NewFromString(v),
  }
}

func (this *ExportedString) Variable() *variable.Variable {
  return this.v
}

func (this *ExportedString) SetVariable(v *variable.Variable) {
  this.v = v
}

func (this *ExportedString) Export() chan *oproto.ValueStream {
  c := make(chan *oproto.ValueStream)
  go func() {
    stream := new(oproto.ValueStream)
    stream.Variable = this.v.AsProto()
    stream.Value = append(stream.Value, &oproto.Value{
      Timestamp: proto.Uint64(uint64(time.Now().UnixNano() / 1000000)),
      StringValue: proto.String(this.value),
    })
    c <- stream
    close(c)
  }()
  return c
}

func (this *ExportedString) Get() string {
  return this.value
}

func (this *ExportedString) Set(value string) {
  this.value = value
}
