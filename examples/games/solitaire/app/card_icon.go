package app

import (
	"image/color"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/solitaire/sols"
	"github.com/t0l1k/eui/examples/games/solitaire/sols/deck"
	"golang.org/x/image/colornames"
)

type CardIcon struct {
	*eui.Container
	cell *sols.Cell
	btn  *eui.Button
	show bool
	f    func(b *eui.Button)
}

func NewCardIcon(cell *sols.Cell, f func(b *eui.Button)) *CardIcon {
	c := &CardIcon{Container: eui.NewContainer(eui.NewAbsoluteLayout())}
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
	c.btn.SetBg(colornames.White)
	c.btn.SetFg(col)
	c.Add(c.btn)
	// c.Hide()
	return c
}

func (c *CardIcon) UpdateData(value *deck.Card) {
	if value == nil {
		c.btn.SetText(" ")
		c.btn.SetFg(colornames.Blue)
	} else {
		c.btn.SetText(value.StringShort())
		c.btn.SetFg(c.cell.GetCard().Color())
	}
	c.MarkDirty()
}

func (d *CardIcon) Update() {
	if d.IsHidden() {
		return
	}
	d.btn.Update()
}

func (c *CardIcon) IsHidden() bool       { return c.show }
func (c *CardIcon) SetHidden(value bool) { c.show = value; c.MarkDirty() }

func (c *CardIcon) SetRect(rect eui.Rect[int]) {
	c.Container.SetRect(rect)
	c.btn.SetRect(rect)
	c.ImageReset()
}
