package eui

import (
	"time"
)

type Timer struct {
	duration, left time.Duration
	run            bool
	onDone         func()
	id             int64
}

func NewTimer(duration time.Duration, fn func()) *Timer {
	t := &Timer{duration: duration, onDone: fn}
	t.id = GetUi().tickListener.Connect(func(ev Event) {
		if !t.run {
			return
		}
		dt := ev.Value.(TickData)
		t.left -= dt.Duration()
		if t.left <= 0 {
			t.run = false
			t.left = t.duration
			if t.onDone != nil {
				t.onDone()
			}
		}
	})
	return t
}
func (t *Timer) SetOnDoneFunc(value func())      { t.onDone = value }
func (t *Timer) SetDuration(value time.Duration) { t.duration = value }
func (t *Timer) Reset()                          { t.left = t.duration }
func (t *Timer) IsOn() bool                      { return t.run }
func (t *Timer) IsOff() bool                     { return !t.run }
func (t *Timer) IsDone() bool                    { return !(t.left > 0) }
func (t *Timer) On() *Timer {
	t.Reset()
	t.run = true
	return t
}
func (t *Timer) Off()                      { t.run = false }
func (t *Timer) TimePassed() time.Duration { return t.duration - t.left }
func (t *Timer) TimeLeft() time.Duration   { return t.left }
func (t *Timer) String() string            { return t.left.String() }
func (t *Timer) Close()                    { GetUi().tickListener.Disconnect(t.id) }
