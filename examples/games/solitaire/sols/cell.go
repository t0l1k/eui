package sols

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/solitaire/sols/deck"
)

type Cell struct {
	eui.SubjectBase
}

func NewCell() *Cell                     { return &Cell{} }
func (c *Cell) SetCard(value *deck.Card) { c.SetValue(value) }
func (c *Cell) GetCard() *deck.Card {
	card := c.Value().(*deck.Card)
	if card != nil {
		return card
	}
	return nil
}
