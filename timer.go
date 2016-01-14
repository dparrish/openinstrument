package openinstrument

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/net/context"

	oproto "github.com/dparrish/openinstrument/proto"
)

const contextLogKey string = "openinstrument_log"

type Timer struct {
	ctx       context.Context
	startTime time.Time
	message   string
}

func GetLog(ctx context.Context) *oproto.OperationLog {
	_, l := getContextWithLog(ctx)
	return l
}

func LogContext(ctx context.Context) context.Context {
	l, _ := getContextWithLog(ctx)
	return l
}

func getContextWithLog(ctx context.Context) (context.Context, *oproto.OperationLog) {
	v := ctx.Value("openinstrument_log")
	if v == nil {
		log := &oproto.OperationLog{}
		return context.WithValue(ctx, "openinstrument_log", log), log
	}
	return ctx, v.(*oproto.OperationLog)
}

func Logf(ctx context.Context, format string, args ...interface{}) context.Context {
	ctx, l := getContextWithLog(ctx)
	log.Output(2, fmt.Sprintf(format, args...))
	l.Log = append(l.Log, &oproto.LogMessage{
		Message:   fmt.Sprintf(format, args...),
		Timestamp: uint64(time.Now().UnixNano()),
	})
	return ctx
}

func StringLog(ctx context.Context) string {
	out := ""
	l := GetLog(ctx).Log
	if len(l) == 0 {
		return out
	}

	st := time.Unix(int64(l[0].Timestamp/1000000000), int64(l[0].Timestamp%1000000000))
	out += fmt.Sprintf("%15s\n", st)
	for i, entry := range l {
		if i == 0 {
			out += fmt.Sprintf("%12d %s\n", 0, entry.Message)
		} else {
			t := time.Unix(int64(entry.Timestamp/1000000000), int64(entry.Timestamp%1000000000))
			d := t.Sub(st)
			ds := fmt.Sprintf("%0.3f", float64(d.Nanoseconds())/1000000.0)
			for x := 0; x < 12-len(ds); x++ {
				out += "."
			}
			out += fmt.Sprintf("%s %s\n", ds, entry.Message)
			st = t
		}
	}
	return out
}

func NewTimer(ctx context.Context, format string, args ...interface{}) *Timer {
	return &Timer{
		startTime: time.Now(),
		ctx:       ctx,
		message:   fmt.Sprintf(format, args...),
	}
}

func (t *Timer) Stop() uint64 {
	duration := uint64(time.Since(t.startTime).Nanoseconds())
	if t.ctx != nil {
		_, l := getContextWithLog(t.ctx)
		l.Log = append(l.Log, &oproto.LogMessage{
			Message:      fmt.Sprintf("%s: %s", t.message, time.Since(t.startTime).String()),
			Timestamp:    uint64(t.startTime.UnixNano()),
			EndTimestamp: uint64(time.Now().UnixNano()),
		})
		log.Output(2, fmt.Sprintf("%s: %s", t.message, time.Since(t.startTime).String()))
	}
	return duration
}
