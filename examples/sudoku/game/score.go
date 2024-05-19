package game

import (
	"fmt"
	"time"
)

type Score struct {
	last, best time.Duration
}

func NewScore(best time.Duration) *Score { return &Score{last: 0, best: best} }

func (s *Score) Last(value time.Duration) {
	s.last = value
	if s.last < s.best || s.best == 0 {
		s.best = value
	}
}

func (s *Score) GetLast() time.Duration { return s.last }
func (s *Score) GetBest() time.Duration { return s.best }
func (s *Score) String() (res string) {
	if s.best == 0 && s.last == 0 {
		res = "Empty Score"
	} else if s.best > 0 && s.last == 0 {
		res = fmt.Sprintf("best(%v)", s.best.Round(time.Millisecond).String())
	} else {
		res = fmt.Sprintf("best(%v) last(%v)", s.best.Round(time.Millisecond).String(), s.last.Round(time.Millisecond).String())
	}
	return res
}
