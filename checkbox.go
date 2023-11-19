package eui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Checkbox struct {
	View
	b       *Button
	t       *Text
	checked bool
}

func NewCheckbox(text string, f func(c *Checkbox)) *Checkbox {
	c := &Checkbox{}
	c.SetupCheckbox(text, f)
	return c
}

func (c *Checkbox) SetupCheckbox(text string, f func(c *Checkbox)) {
	theme := GetUi().theme
	c.Bg(theme.Get(CheckboxBg))
	c.Fg(theme.Get(CheckboxFg))
	c.b = NewButton("", func(b *Button) {
		c.checked = !c.checked
		f(c)
		c.SetChecked(c.checked)
	})
	c.t = NewText(text)
	c.t.Bg(c.bg)
	c.t.Fg(c.fg)
}

func (c *Checkbox) IsChecked() bool { return c.checked }
func (c *Checkbox) SetChecked(value bool) {
	c.checked = value
	if c.checked {
		c.b.SetText("*")
	} else {
		c.b.SetText(" ")
	}
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
