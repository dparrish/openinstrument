package openinstrument

import (
	"fmt"
	"time"

	"github.com/dparrish/openinstrument/proto"
)

type Timer struct {
	t         *openinstrument_proto.LogMessage
	startTime time.Time
	message   string
}

func NewTimer(message string, t *openinstrument_proto.LogMessage) *Timer {
	return &Timer{
		startTime: time.Now(),
		t:         t,
		message:   message,
	}
}

func (t *Timer) Stop() uint64 {
	duration := time.Since(t.startTime)
	if t.t != nil {
		t.t.Timestamp = uint64(duration.Nanoseconds() / 1000000)
		if t.message != "" {
			t.t.Message = t.message
		}
	}
	return uint64(duration.Nanoseconds() / 1000000)
}

type DurationTimer struct {
	name      string
	startTime time.Time
	endTime   time.Time
	totalTime time.Duration
	running   bool
}

func NewNamedDurationTimer(name string) *DurationTimer {
	return &DurationTimer{
		name: name,
	}
}

func NewDurationTimer() *DurationTimer {
	return &DurationTimer{}
}

func (t *DurationTimer) Start() {
	if !t.running {
		t.startTime = time.Now()
		t.running = true
	}
}

func (t *DurationTimer) Stop() {
	if t.running {
		t.endTime = time.Now()
		d := t.endTime.Sub(t.startTime)
		t.totalTime = time.Duration(d.Nanoseconds() + t.totalTime.Nanoseconds())
	}
}

func (t *DurationTimer) Duration() time.Duration {
	if t.running {
		t.Stop()
	}
	return t.totalTime
}

func (t *DurationTimer) String() string {
	if t.name != "" {
		return fmt.Sprintf("%s: %s", t.name, t.Duration())
	}
	return t.Duration().String()
}
