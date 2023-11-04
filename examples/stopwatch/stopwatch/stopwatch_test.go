package stopwatch

import (
	"fmt"
	"testing"
	"time"

	"github.com/t0l1k/eui"
)

func TestStopwatchStart(t *testing.T) {
	got := int(eui.NewStopwatch().Duration())
	want := 0
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestStopwatchGot3sec(t *testing.T) {
	s := eui.NewStopwatch()
	s.Start()
	time.Sleep(3 * time.Second)
	s.Stop()
	got := fmt.Sprintf("%.2f", s.Duration().Seconds())
	want := "3.00"
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}

	got1 := s.Duration().Milliseconds()
	want1 := 3000
	if got != want {
		t.Errorf("got %v want %v", got1, want1)
	}
}

func TestStopwatchGot3secPauseAnd2sec(t *testing.T) {
	s := eui.NewStopwatch()
	s.Start()
	time.Sleep(3 * time.Second)
	s.Stop()
	time.Sleep(1 * time.Second)
	s.Start()
	time.Sleep(2 * time.Second)
	s.Stop()
	got := fmt.Sprintf("%.2f", s.Duration().Seconds())
	want := "5.00"
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}

	got1 := int64(s.Duration().Truncate(10 * time.Millisecond).Milliseconds())
	want1 := int64(5000)
	if got1 != want1 {
		t.Errorf("got %v want %v", got1, want1)
	}

}
