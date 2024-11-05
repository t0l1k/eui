package sols

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/solitaire/sols/deck"
)

type Cell struct {
	eui.SubjectBase
}

func NewCell() *Cell                           { return &Cell{} }
func (c *Cell) Reset()                         { c.SubjectBase.Reset() }
func (c *Cell) IsEmpty() bool                  { return c.GetCard() == nil }
func (c *Cell) SetCard(value *deck.Card) *Cell { c.SetValue(value); return c }
func (c *Cell) GetCard() *deck.Card {
	if card := c.Value(); card != nil {
		return card.(*deck.Card)
	}
	return nil
}
