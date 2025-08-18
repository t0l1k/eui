package eui

import (
	"fmt"
	"strings"
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

func FormatSmartDuration(dur time.Duration, showMS bool) string {
	h := int(dur.Hours())
	m := int(dur.Minutes()) % 60
	s := int(dur.Seconds()) % 60
	ms := int(dur.Milliseconds()) % 100

	var sb strings.Builder

	// Отображаем часы, если они > 0
	if h > 0 {
		fmt.Fprintf(&sb, "%02d:", h)
	}

	// Отображаем минуты, если они > 0 или уже были часы
	if m > 0 || h > 0 {
		fmt.Fprintf(&sb, "%02d:", m)
	}

	// Отображаем секунды, если они > 0 или были минуты/часы
	if s >= 0 || m > 0 || h > 0 {
		fmt.Fprintf(&sb, "%02d", s)
	}

	// Отображаем миллисекунды:
	if showMS {
		if s >= 0 || m > 0 || h > 0 {
			sb.WriteString(",")
		}
		fmt.Fprintf(&sb, "%02d", ms)
	}

	// Если вообще ничего не отображено кроме миллисекунд
	if sb.Len() == 0 {
		if showMS {
			fmt.Fprintf(&sb, "%02d", ms)
		} else {
			fmt.Fprintf(&sb, "00")
		}
	}

	return sb.String()
}
