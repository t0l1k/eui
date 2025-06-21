package eui

import "github.com/hajimehoshi/ebiten/v2"

type Checkbox struct {
	*Drawable
	btn     *Button
	txt     *Text
	checked bool
}

func NewCheckbox(text string, f func(c *Checkbox)) *Checkbox {
	c := &Checkbox{Drawable: NewDrawable()}
	c.SetupCheckbox(text, f)
	c.Visible(true)
	return c
}

func (c *Checkbox) SetupCheckbox(text string, f func(c *Checkbox)) {
	theme := GetUi().theme
	c.Bg(theme.Get(CheckboxBg))
	c.Fg(theme.Get(CheckboxFg))
	c.btn = NewButton("", func(b *Button) {
		c.checked = !c.checked
		f(c)
		c.SetChecked(c.checked)
	})
	c.txt = NewText(text)
	c.txt.Bg(c.bg)
	c.txt.Fg(c.fg)
}

func (c *Checkbox) GetText() string      { return c.txt.GetText() }
func (c *Checkbox) SetText(value string) { c.txt.SetText(value) }
func (c *Checkbox) IsChecked() bool      { return c.checked }
func (c *Checkbox) SetChecked(value bool) {
	c.checked = value
	if c.checked {
		c.btn.SetText("*")
	} else {
		c.btn.SetText(" ")
	}
	c.MarkDirty()
}

func (c *Checkbox) Layout() {
	c.Drawable.Layout()
	c.ClearDirty()
}

func (c *Checkbox) Update(dt int) {
	if c.disabled {
		return
	}
	c.btn.Update(dt)
}

func (b *Checkbox) Draw(surface *ebiten.Image) {
	if !b.IsVisible() {
		return
	}
	if b.IsDirty() {
		b.Layout()
		b.btn.Layout()
		b.txt.Layout()
	}
	b.txt.Draw(surface)
	b.btn.Draw(surface)
}

func (c *Checkbox) Resize(rect Rect) {
	c.SetRect(rect)
	w0, h0 := c.rect.Size()
	x, y := c.rect.X, c.rect.Y
	c.btn.Resize(NewRect([]int{x, y, h0, h0}))
	x += h0
	c.txt.Resize(NewRect([]int{x, y, w0 - h0, h0}))
}
