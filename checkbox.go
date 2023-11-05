package eui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Checkbox struct {
	View
	b       *Button
	t       *Text
	checked bool
}

func NewCheckbox(text string, bg, fg color.Color, f func(c *Checkbox)) *Checkbox {
	c := &Checkbox{}
	c.SetupCheckbox(text, f)
	return c
}

func (c *Checkbox) SetupCheckbox(text string, f func(c *Checkbox)) {
	c.b = NewButton("", func(b *Button) {
		c.checked = !c.checked
		if c.checked {
			c.b.SetText("*")
			f(c)
		} else {
			c.b.SetText(" ")
		}
	})
	c.t = NewText(text)
}

func (c *Checkbox) IsChecked() bool { return c.checked }
func (c *Checkbox) SetChecked(value bool) {
	if c.checked == value {
		return
	}
	c.checked = value
	c.dirty = true
}

func (c *Checkbox) Update(dt int) {
	c.b.Update(dt)
}

func (c *Checkbox) Draw(surface *ebiten.Image) {
	c.t.Draw(surface)
	c.b.Draw(surface)
}

func (c *Checkbox) Resize(rect []int) {
	c.Rect(rect)
	w0, h0 := c.rect.Size()
	x, y := c.rect.X, c.rect.Y
	c.b.Resize([]int{x, y, h0, h0})
	x += h0
	c.t.Resize([]int{x, y, w0 - h0, h0})
}
