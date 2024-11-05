package deck_test

import (
	"testing"

	"github.com/t0l1k/eui/examples/games/solitaire/sols/deck"
)

func TestDeckCards52Init(t *testing.T) {
	d := deck.NewDeckCards52().Shuffle()
	got := d.Len()
	want := 52
	if got != want {
		t.Error("Failed Test Deck Len 52", d, len(d.Deck52()))
	}
}
