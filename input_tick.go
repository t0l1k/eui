package eui

import (
	"fmt"
	"time"
)

type TickData struct {
	tick     time.Time
	duration time.Duration
}

func NewTickData(t time.Time, d time.Duration) TickData { return TickData{tick: t, duration: d} }

func (t *TickData) Tick() time.Time         { return t.tick }
func (t *TickData) Duration() time.Duration { return t.duration }
func (t *TickData) String() string {
	return fmt.Sprintf("tick:[%v] duration:[%v]", t.tick.Format("15:04:05.000"), t.duration.String())
}

type TickListener struct {
	*Signal[Event]
	duration time.Duration
}

func NewTickListener(fn SlotFunc[Event], i time.Duration) *TickListener {
	t := &TickListener{Signal: NewSignal(func(a, b Event) bool { return false }), duration: i}
	t.Connect(fn)
	return t
}

func (t *TickListener) update(int) {
	value, ok := t.Value().Value.(TickData)
	if !ok {
		t.Emit(NewEvent(EventTick, NewTickData(time.Now(), 0)))
		return
	}
	lastTick := value.tick
	duration := time.Since(lastTick)
	if duration > t.duration {
		t.Emit(NewEvent(EventTick, NewTickData(time.Now(), duration)))
	}
}
