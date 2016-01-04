package exported

import (
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/dparrish/openinstrument/client"
	oproto "github.com/dparrish/openinstrument/proto"
	"github.com/dparrish/openinstrument/variable"
)

type VariableExporter struct {
	ctx  context.Context
	wg   *sync.WaitGroup
	vars []Exportable
}

func NewExporter(ctx context.Context, c client.Client, interval time.Duration) *VariableExporter {
	e := &VariableExporter{
		ctx: ctx,
		wg:  &sync.WaitGroup{},
	}
	go func() {
		e.wg.Add(1)
		defer e.wg.Done()
		tick := time.Tick(interval)
		for {
			select {
			case <-ctx.Done():
				break
			case <-tick:
				e.GetStreams()
			}
		}
	}()
	return e
}

func (e *VariableExporter) GetStreams() []*oproto.ValueStream {
	// Flush
	streams := []*oproto.ValueStream{}
	for _, v := range e.vars {
		for stream := range v.Export() {
			streams = append(streams, stream)
		}
	}
	return streams
}

func (e *VariableExporter) Shutdown() {
	e.wg.Wait()
	e.vars = []Exportable{}
}

func (e *VariableExporter) addExportable(v Exportable) {
	e.vars = append(e.vars, v)
}

type Exportable interface {
	Export() chan *oproto.ValueStream
}

//////////// Integer ///////////////////////////////
type Integer struct {
	value uint64
	v     *variable.Variable
}

func NewInteger(exporter *VariableExporter, v *variable.Variable) *Integer {
	i := &Integer{
		value: 0,
		v:     v,
	}
	exporter.addExportable(i)
	return i
}

func (ei *Integer) Export() chan *oproto.ValueStream {
	c := make(chan *oproto.ValueStream)
	go func() {
		defer close(c)
		c <- &oproto.ValueStream{
			Variable: ei.v.AsProto(),
			Value: []*oproto.Value{
				{
					Timestamp:   uint64(time.Now().UnixNano() / 1000000),
					DoubleValue: float64(ei.value),
				},
			},
		}
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

func NewFloat(exporter *VariableExporter, v *variable.Variable) *Float {
	f := &Float{
		value: 0,
		v:     v,
	}
	exporter.addExportable(f)
	return f
}

func (ef *Float) Export() chan *oproto.ValueStream {
	c := make(chan *oproto.ValueStream)
	go func() {
		defer close(c)
		c <- &oproto.ValueStream{
			Variable: ef.v.AsProto(),
			Value: []*oproto.Value{
				{
					Timestamp:   uint64(time.Now().UnixNano() / 1000000),
					DoubleValue: ef.value,
				},
			},
		}
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

func NewRatio(exporter *VariableExporter, v *variable.Variable) *Ratio {
	successVar := variable.NewFromProto(v.AsProto())
	successVar.Variable += "-success"
	failureVar := variable.NewFromProto(v.AsProto())
	failureVar.Variable += "-failure"
	totalVar := variable.NewFromProto(v.AsProto())
	totalVar.Variable += "-total"
	return &Ratio{
		success: NewInteger(exporter, successVar),
		failure: NewInteger(exporter, failureVar),
		total:   NewInteger(exporter, totalVar),
	}
}

func (er *Ratio) Export() chan *oproto.ValueStream {
	c := make(chan *oproto.ValueStream)
	go func() {
		defer close(c)
		for stream := range er.success.Export() {
			c <- stream
		}
		for stream := range er.failure.Export() {
			c <- stream
		}
		for stream := range er.total.Export() {
			c <- stream
		}
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

func NewAverage(exporter *VariableExporter, v *variable.Variable) *Average {
	sumVar := variable.NewFromProto(v.AsProto())
	sumVar.Variable += "-total-count"
	totalVar := variable.NewFromProto(v.AsProto())
	totalVar.Variable += "-overall-sum"
	return &Average{
		overallSum: NewFloat(exporter, sumVar),
		totalCount: NewInteger(exporter, totalVar),
	}
}

func (ea *Average) Export() chan *oproto.ValueStream {
	c := make(chan *oproto.ValueStream)
	go func() {
		defer close(c)
		for stream := range ea.overallSum.Export() {
			c <- stream
		}
		for stream := range ea.totalCount.Export() {
			c <- stream
		}
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

func NewTimer(exporter *VariableExporter, v *variable.Variable) *Timer {
	return &Timer{
		average: NewAverage(exporter, v),
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

func NewString(exporter *VariableExporter, v *variable.Variable) *String {
	s := &String{
		value: "",
		v:     v,
	}
	exporter.addExportable(s)
	return s
}

func (es *String) Export() chan *oproto.ValueStream {
	c := make(chan *oproto.ValueStream)
	go func() {
		defer close(c)
		c <- &oproto.ValueStream{
			Variable: es.v.AsProto(),
			Value: []*oproto.Value{
				{
					Timestamp:   uint64(time.Now().UnixNano() / 1000000),
					StringValue: es.value,
				},
			},
		}
	}()
	return c
}

func (es *String) Get() string {
	return es.value
}

func (es *String) Set(value string) {
	es.value = value
}
