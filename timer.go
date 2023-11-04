package eui

import (
	"strconv"
	"time"
)

// Умею определеленную паузу отсчитывать обратно в миллисекундах. Умею отсчитывать регулярно обновляя таймер через метод update, или после разовой проверки методом endCheck
type Timer struct {
	startAt       time.Time
	timer         int
	pauseDuration int
	run           bool
}

func NewTimer(duration int) *Timer { return &Timer{pauseDuration: duration} }

func (t *Timer) Reset() { t.timer = t.pauseDuration }

func (t *Timer) IsOn() bool { return t.run }

func (t *Timer) On() {
	t.Reset()
	t.run = true
	t.startAt = time.Now()
}

func (t *Timer) IsDone() bool { return !(t.timer > 0) }
func (t *Timer) Off()         { t.run = false }

func (t *Timer) EndCheck() {
	dt := time.Since(t.startAt).Milliseconds()
	if dt > int64(t.pauseDuration) {
		t.Off()
	}
}

func (t *Timer) Update(dt int) {
	if t.IsOn() {
		t.timer -= dt
	}
}

func (t *Timer) String() string {
	s := "timer "
	if t.run {
		s += "on:"
	} else {
		s += "off:"
	}
	s += strconv.Itoa(t.timer)
	return s
}
