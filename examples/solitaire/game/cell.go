package game

import "github.com/t0l1k/eui"

type Cell struct {
	eui.SubjectBase
}

func NewCell() *Cell                { return &Cell{} }
func (c *Cell) SetCard(value *Card) { c.SetValue(value) }
func (c *Cell) GetCard() *Card {
	card := c.Value().(*Card)
	if card != nil {
		return card
	}
	return nil
}
