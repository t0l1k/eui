package mines

import (
	"fmt"
)

type Timer struct {
	tm    int
	pause bool
}

func NewTimer() *Timer {
	return &Timer{pause: true}
}

func (t *Timer) Reset() {
	t.tm = 0
	t.pause = true
}

func (t *Timer) Start() {
	t.tm = 0
	t.pause = false
}

func (t *Timer) Update(tm int) {
	if t.pause {
		return
	}
	t.tm += tm
}

func (t *Timer) Pause() {
	t.pause = true
}

func (t *Timer) Resume() {
	t.pause = false
}

func (t *Timer) Stop() {
	t.pause = true
}

func (t *Timer) String() string {
	seconds := t.tm / 1000 % 60
	minutes := t.tm / 1000 / 60
	return fmt.Sprintf("%02v:%02v", minutes, seconds)
}
