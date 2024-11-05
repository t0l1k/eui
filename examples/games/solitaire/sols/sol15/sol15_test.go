package sol15

import (
	"fmt"
	"testing"

	"github.com/t0l1k/eui/examples/games/solitaire/sols/deck"
)

func TestLayout15Init(t *testing.T) {
	d := NewLayout15(deck.NewDeckCards52().Shuffle())
	got := d.Column(14)[0].GetCard()
	if got != nil {
		t.Error("Test Layout15 Game Init Last Row", d, got)
	}
	got = d.Column(0)[0].GetCard()
	if got == nil {
		t.Error("Test Layout15 Game Init First Row", d, got)
	}
}

func TestMove(t *testing.T) {
	d := NewLayout15(deck.NewDeckCards52().Shuffle())
	fmt.Println(d)
	d.MakeMove(0)
	got := d.columnLastCard(13)
	if got == nil {
		t.Error("Test Layout15 Game Make Move from row 0 to row 13", d, got, d.Column(0), d.Column(13))
	}
	d.MakeMove(0)
	got = d.columnLastCard(14)
	if got == nil {
		t.Error("Test Layout15 Game Make Move", d, got, d.Column(0)[2])
	}
}
