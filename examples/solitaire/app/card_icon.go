package app

import (
	"image/color"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/solitaire/game"
)

type CardIcon struct {
	eui.DrawableBase
	cell *game.Cell
	btn  *eui.Button
	show bool
	f    func(b *eui.Button)
}

func NewCardIcon(cell *game.Cell, f func(b *eui.Button)) *CardIcon {
	c := &CardIcon{}
	c.cell = cell
	c.f = f
	str := ""
	var col color.Color
	if c.cell.GetCard() == nil {
		col = eui.Blue
		str = ""
	} else {
		col = c.cell.GetCard().Color()
		str = c.cell.GetCard().StringShort()
	}
	c.btn = eui.NewButton(str, f)
	c.btn.Bg(eui.White)
	c.btn.Fg(col)
	c.Add(c.btn)
	c.Visible(true)
	return c
}

func (c *CardIcon) UpdateData(value interface{}) {
	switch v := value.(type) {
	case *game.Card:
		if v == nil {
			c.btn.SetText(" ")
			c.btn.Fg(eui.Blue)
		} else {
			c.btn.SetText(v.StringShort())
			c.btn.Fg(c.cell.GetCard().Color())
		}
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
