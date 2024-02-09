package eui

type Checkbox struct {
	View
	btn     *Button
	txt     *Text
	checked bool
}

func NewCheckbox(text string, f func(c *Checkbox)) *Checkbox {
	c := &Checkbox{}
	c.SetupCheckbox(text, f)
	return c
}

func (c *Checkbox) SetupCheckbox(text string, f func(c *Checkbox)) {
	c.SetupView()
	theme := GetUi().theme
	c.Bg(theme.Get(CheckboxBg))
	c.Fg(theme.Get(CheckboxFg))
	c.btn = NewButton("", func(b *Button) {
		c.checked = !c.checked
		f(c)
		c.SetChecked(c.checked)
	})
	c.Add(c.btn)
	c.txt = NewText(text)
	c.txt.Bg(c.bg)
	c.txt.Fg(c.fg)
	c.Add(c.txt)
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
	c.dirty = true
}

func (c *Checkbox) Resize(rect []int) {
	c.Rect(rect)
	w0, h0 := c.rect.Size()
	x, y := c.rect.X, c.rect.Y
	c.btn.Resize([]int{x, y, h0, h0})
	x += h0
	c.txt.Resize([]int{x, y, w0 - h0, h0})
}
