package sols

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/solitaire/sols/deck"
)

type Cell struct{ *eui.Signal }

func NewCell() *Cell                           { return &Cell{Signal: eui.NewSignal()} }
func (c *Cell) Reset()                         { c.Emit(nil) }
func (c *Cell) IsEmpty() bool                  { return c.GetCard() == nil }
func (c *Cell) SetCard(value *deck.Card) *Cell { c.Emit(value); return c }
func (c *Cell) GetCard() *deck.Card {
	if card := c.Value(); card != nil {
		return card.(*deck.Card)
	}
	return nil
}
