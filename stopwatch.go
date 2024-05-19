package eui

import (
	"fmt"
	"log"
	"time"
)

// Умею запустить отсчет времени, умею остановить и возобновить отсчет секундомера
type Stopwatch struct {
	run      bool
	startAt  time.Time
	duration time.Duration
}

func NewStopwatch() *Stopwatch {
	s := &Stopwatch{}
	return s
}

// Запустить секундомер
func (s *Stopwatch) Start() {
	s.startAt = timeNow()
	s.run = true
	log.Println("Stopwatch start")
}

// Остановить или поставить на паузу
func (s *Stopwatch) Stop() {
	if s.run {
		s.duration += timeNow().Sub(s.startAt)
		s.run = false
		log.Println("Stopwatch stop")
	}
}

func (s *Stopwatch) IsRun() bool {
	return s.run
}

// Обнулить секундомер
func (s *Stopwatch) Reset() {
	s.run = false
	s.duration = 0
	log.Println("Stopwatch reset")
}

// Возвращаю сколько прошло времени от запуска секундомера
func (s *Stopwatch) Duration() (duration time.Duration) {
	if !s.run {
		return s.duration
	}
	duration = timeNow().Sub(s.startAt)
	duration += s.duration
	return duration
}

var timeNow = func() time.Time {
	return time.Now()
}

// Возвращаю строку в формате часов включая миллисекунды, секунды или минуты, или часы
func (s *Stopwatch) String() string {
	var (
		duration time.Duration
		str      string
	)
	duration = s.Duration()
	mSec := int(duration.Milliseconds()) % 1000
	sec := int(duration.Seconds())
	seconds := sec % 60
	minutes := sec / 60
	hour := minutes / 60
	if hour > 0 {
		str = fmt.Sprintf("%vh%02vm%02vs", hour, minutes, seconds)
	} else if minutes > 0 {
		str = fmt.Sprintf("%vm%02v.%02vs", minutes, seconds, mSec)
	} else {
		str = fmt.Sprintf("%v.%02vs", seconds, mSec)
	}
	return str
}

// Возвращаю строку в формате часов включая секунды, минуты или часы
func (s *Stopwatch) StringShort() string {
	var (
		duration time.Duration
		str      string
	)
	duration = s.Duration()
	sec := int(duration.Seconds())
	seconds := sec % 60
	minutes := sec / 60
	hour := minutes / 60
	if hour > 0 {
		str = fmt.Sprintf("%v:%02v:%02v", hour, minutes, seconds)
	} else {
		str = fmt.Sprintf("%v:%02v", minutes, seconds)
	}
	return str
}
