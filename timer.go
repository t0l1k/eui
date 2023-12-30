package eui

import "strconv"

// Умею определеленную длительность времени в миллисекундах отсчитывать обратно. Умею отсчитывать регулярно обновляя таймер через метод update
type Timer struct {
	timer    int
	duration int
	run      bool
}

func NewTimer(duration int) *Timer { return &Timer{duration: duration} }

func (t *Timer) SetDuration(value int) { t.duration = value }
func (t *Timer) Reset()                { t.timer = t.duration }
func (t *Timer) IsOn() bool            { return t.run }
func (t *Timer) IsOff() bool           { return !t.run }
func (t *Timer) IsDone() bool          { return !(t.timer > 0) }

func (t *Timer) On() {
	t.Reset()
	t.run = true
}
func (t *Timer) Off() { t.run = false }

func (t *Timer) Update(dt int) {
	if t.IsOn() {
		t.timer -= dt
	}
}

func (t *Timer) TimePassed() int { return t.duration - t.timer }
func (t *Timer) TimeLeft() int   { return t.timer }

func (t *Timer) String() string {
	return strconv.Itoa(t.timer/1000 + 1)
}
