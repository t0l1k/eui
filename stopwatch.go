package eui

import (
	"time"
)

// Умею запустить отсчет времени, умею остановить и возобновить отсчет секундомера
type Stopwatch struct {
	run      bool
	startAt  time.Time
	duration time.Duration
}

func NewStopwatch() *Stopwatch { return &Stopwatch{} }

// Запустить секундомер
func (s *Stopwatch) Start() *Stopwatch {
	s.startAt = time.Now()
	s.run = true
	return s
}

// Остановить или поставить на паузу
func (s *Stopwatch) Stop() *Stopwatch {
	if s.run {
		s.duration += time.Since(s.startAt)
		s.run = false
	}
	return s
}

func (s *Stopwatch) IsRun() bool { return s.run }

// Обнулить секундомер
func (s *Stopwatch) Reset() *Stopwatch {
	s.run = false
	s.duration = 0
	return s
}

// Возвращаю сколько прошло времени от запуска секундомера
func (s *Stopwatch) Duration() (duration time.Duration) {
	if !s.run {
		return s.duration
	}
	duration = time.Since(s.startAt)
	duration += s.duration
	return duration
}
