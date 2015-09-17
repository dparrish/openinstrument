package exported

import (
	"log"
	"sync"
	"time"

	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/variable"
)

type VariableExporter struct {
	shutdown  bool
	completed sync.Mutex
}

func NewVariableExporter(address string, interval int32) *VariableExporter {
	ve := new(VariableExporter)
	go func(ve *VariableExporter) {
		ve.completed.Lock()
		defer ve.completed.Unlock()
		tick := time.Tick(time.Duration(interval) * time.Millisecond)
		for !ve.shutdown {
			<-tick
			// Flush
			log.Printf("Flushing exported variables")
		}
	}(ve)
	return ve
}

func (ve *VariableExporter) Shutdown() {
	ve.shutdown = true
	ve.completed.Lock()
	ve.completed.Unlock()
}

type Exportable interface {
	Export() chan *oproto.ValueStream
	Variable() *variable.Variable
	SetVariable(v *variable.Variable)
}

//////////// Integer ///////////////////////////////
type Integer struct {
	value uint64
	v     *variable.Variable
}

func NewInteger(v string) *Integer {
	return &Integer{
		value: 0,
		v:     variable.NewFromString(v),
	}
}

func (ei *Integer) Variable() *variable.Variable {
	return ei.v
}

func (ei *Integer) SetVariable(v *variable.Variable) {
	ei.v = v
}

func (ei *Integer) Export() chan *oproto.ValueStream {
	c := make(chan *oproto.ValueStream)
	go func() {
		stream := new(oproto.ValueStream)
		stream.Variable = ei.v.AsProto()
		stream.Value = append(stream.Value, &oproto.Value{
			Timestamp:   uint64(time.Now().UnixNano() / 1000000),
			DoubleValue: float64(ei.value),
		})
		c <- stream
		close(c)
	}()
	return c
}

func (ei *Integer) Get() uint64 {
	return ei.value
}

func (ei *Integer) Set(value uint64) {
	ei.value = value
}

func (ei *Integer) Add(value uint64) {
	ei.value += value
}

func (ei *Integer) Sub(value uint64) {
	ei.value -= value
}

//////////// Float ///////////////////////////////
type Float struct {
	value float64
	v     *variable.Variable
}

func NewFloat(v string) *Float {
	return &Float{
		value: 0,
		v:     variable.NewFromString(v),
	}
}

func (ef *Float) Variable() *variable.Variable {
	return ef.v
}

func (ef *Float) SetVariable(v *variable.Variable) {
	ef.v = v
}

func (ef *Float) Export() chan *oproto.ValueStream {
	c := make(chan *oproto.ValueStream)
	go func() {
		stream := new(oproto.ValueStream)
		stream.Variable = ef.v.AsProto()
		stream.Value = append(stream.Value, &oproto.Value{
			Timestamp:   uint64(time.Now().UnixNano() / 1000000),
			DoubleValue: ef.value,
		})
		c <- stream
		close(c)
	}()
	return c
}

func (ef *Float) Get() float64 {
	return ef.value
}

func (ef *Float) Set(value float64) {
	ef.value = value
}

func (ef *Float) Add(value float64) {
	ef.value += value
}

func (ef *Float) Sub(value float64) {
	ef.value -= value
}

//////////// Ratio ///////////////////////////////
type Ratio struct {
	success *Integer
	failure *Integer
	total   *Integer
}

func NewRatio(v string) *Ratio {
	successVar := variable.NewFromString(v)
	successVar.Variable += "-success"
	failureVar := variable.NewFromString(v)
	failureVar.Variable += "-failure"
	totalVar := variable.NewFromString(v)
	totalVar.Variable += "-total"
	return &Ratio{
		success: NewInteger(successVar.String()),
		failure: NewInteger(failureVar.String()),
		total:   NewInteger(totalVar.String()),
	}
}

func (er *Ratio) Export() chan *oproto.ValueStream {
	c := make(chan *oproto.ValueStream)
	go func() {
		for stream := range er.success.Export() {
			c <- stream
		}
		for stream := range er.failure.Export() {
			c <- stream
		}
		for stream := range er.total.Export() {
			c <- stream
		}
		close(c)
	}()
	return c
}

func (er *Ratio) Success() {
	er.success.Add(1)
	er.total.Add(1)
}

func (er *Ratio) Failure() {
	er.failure.Add(1)
	er.total.Add(1)
}

//////////// Average ///////////////////////////////
type Average struct {
	overallSum *Float
	totalCount *Integer
}

func NewAverage(v string) *Average {
	sumVar := variable.NewFromString(v)
	sumVar.Variable += "-total-count"
	totalVar := variable.NewFromString(v)
	sumVar.Variable += "-overall-sum"
	return &Average{
		overallSum: NewFloat(sumVar.String()),
		totalCount: NewInteger(totalVar.String()),
	}
}

func (ea *Average) Export() chan *oproto.ValueStream {
	c := make(chan *oproto.ValueStream)
	go func() {
		for stream := range ea.overallSum.Export() {
			c <- stream
		}
		for stream := range ea.totalCount.Export() {
			c <- stream
		}
		close(c)
	}()
	return c
}

func (ea *Average) Update(sum float64, count uint64) {
	ea.overallSum.Add(sum)
	ea.totalCount.Add(count)
}

//////////// Timer ///////////////////////////////
type Timer struct {
	average   *Average
	startTime time.Time
}

func NewTimer(v string) *Timer {
	return &Timer{
		average: NewAverage(v),
	}
}

func (et *Timer) Export() chan *oproto.ValueStream {
	return et.average.Export()
}

func (et *Timer) Start() {
	et.startTime = time.Now()
}

func (et *Timer) Stop() {
	duration := time.Since(et.startTime)
	et.average.Update(float64(duration.Nanoseconds()/1000000), 1)
}

//////////// String ///////////////////////////////
type String struct {
	value string
	v     *variable.Variable
}

func NewString(v string) *String {
	return &String{
		value: "",
		v:     variable.NewFromString(v),
	}
}

func (es *String) Variable() *variable.Variable {
	return es.v
}

func (es *String) SetVariable(v *variable.Variable) {
	es.v = v
}

func (es *String) Export() chan *oproto.ValueStream {
	c := make(chan *oproto.ValueStream)
	go func() {
		stream := new(oproto.ValueStream)
		stream.Variable = es.v.AsProto()
		stream.Value = append(stream.Value, &oproto.Value{
			Timestamp:   uint64(time.Now().UnixNano() / 1000000),
			StringValue: es.value,
		})
		c <- stream
		close(c)
	}()
	return c
}

func (es *String) Get() string {
	return es.value
}

func (es *String) Set(value string) {
	es.value = value
}
