package app

import (
	"image/color"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/solitaire/sols"
	"github.com/t0l1k/eui/examples/games/solitaire/sols/deck"
	"golang.org/x/image/colornames"
)

type CardIcon struct {
	eui.DrawableBase
	cell *sols.Cell
	btn  *eui.Button
	show bool
	f    func(b *eui.Button)
}

func NewCardIcon(cell *sols.Cell, f func(b *eui.Button)) *CardIcon {
	c := &CardIcon{}
	c.cell = cell
	c.f = f
	str := ""
	var col color.Color
	if c.cell.IsEmpty() {
		col = colornames.Blue
		str = ""
	} else {
		col = c.cell.GetCard().Color()
		str = c.cell.GetCard().StringShort()
	}
	c.btn = eui.NewButton(str, f)
	c.btn.Bg(colornames.White)
	c.btn.Fg(col)
	c.Add(c.btn)
	c.Visible(true)
	return c
}

func (c *CardIcon) UpdateData(value *deck.Card) {
	if value == nil {
		c.btn.SetText(" ")
		c.btn.Fg(colornames.Blue)
	} else {
		c.btn.SetText(value.StringShort())
		c.btn.Fg(c.cell.GetCard().Color())
	}
	c.Dirty = true
}

func (d *CardIcon) Update(dt int) {
	if !d.IsVisible() {
		return
	}
	d.btn.Update(dt)
}

func (c *CardIcon) IsVisible() bool    { return c.show }
func (c *CardIcon) Visible(value bool) { c.show = value }

func (c *CardIcon) Resize(rect []int) {
	c.Rect(eui.NewRect(rect))
	c.SpriteBase.Resize(rect)
	c.btn.Resize(rect)
	c.ImageReset()
}
