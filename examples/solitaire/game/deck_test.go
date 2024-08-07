package game

import (
	"fmt"
	"testing"
)

func TestDeckCards52Init(t *testing.T) {
	d := NewDeckCards52().Shuffle()
	got := d.Len()
	want := 52
	if got != want {
		t.Error("Failed Test Deck Len 52", d, len(d.deck))
	}
}

func TestLayout15Init(t *testing.T) {
	d := NewLayout15(NewDeckCards52().Shuffle())
	got := d.Row(14)[0].GetCard()
	if got != nil {
		t.Error("Test Layout15 Game Init Last Row", d, got)
	}
	got = d.Row(0)[0].GetCard()
	if got == nil {
		t.Error("Test Layout15 Game Init First Row", d, got)
	}
}

func TestMove(t *testing.T) {
	d := NewLayout15(NewDeckCards52().Shuffle())
	fmt.Println(d)
	d.MakeMove(0)
	got := d.RowLastCard(13)
	if got == nil {
		t.Error("Test Layout15 Game Make Move from row 0 to row 13", d, got, d.Row(0), d.Row(13))
	}
	d.MakeMove(0)
	got = d.RowLastCard(14)
	if got != nil {
		t.Error("Test Layout15 Game Make Move", d, got, d.Row(0)[2])
	}
}
