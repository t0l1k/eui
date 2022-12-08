package app

import "testing"

func TestCountInit(t *testing.T) {
	c := NewCounter()
	c.Inc() //1
	got := c.Get()
	want := 1
	if got != Counter(want) {
		t.Errorf("Error: got %v want %v", got, want)
	}
	c.Inc() //2
	c.Inc() //3
	c.Dec() //2
	got = c.Get()
	want = 2
	if got != Counter(want) {
		t.Errorf("Error: got %v want %v", got, want)
	}
	c.Dec() //1
	c.Dec() //0
	c.Dec() //0
	got = c.Get()
	want = 0
	if got != Counter(want) {
		t.Errorf("Error: got %v want %v", got, want)
	}
}
