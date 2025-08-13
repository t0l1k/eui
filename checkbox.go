package eui

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Checkbox struct {
	*Drawable
	txt     string
	checked bool
	onClick *Signal[*Checkbox]
}

func NewCheckbox(txt string, fn func(c *Checkbox)) *Checkbox {
	c := &Checkbox{Drawable: NewDrawable(), onClick: NewSignal[*Checkbox](), txt: txt}
	c.onClick.Connect(fn)
	theme := GetUi().theme
	c.Bg(theme.Get(CheckboxBg))
	c.Fg(theme.Get(CheckboxFg))
	return c
}

func (c *Checkbox) Text() string          { return c.txt }
func (c *Checkbox) SetText(value string)  { c.txt = value; c.MarkDirty() }
func (c *Checkbox) IsChecked() bool       { return c.checked }
func (c *Checkbox) SetChecked(value bool) { c.checked = value; c.MarkDirty() }

func (c *Checkbox) Hit(pt Point[int]) Drawabler {
	if !pt.In(c.rect) {
		return nil
	}
	log.Println("Checkbox:Hit:", c.Text(), c.checked, c.Rect(), pt)
	return c
}
func (c *Checkbox) MouseUp(md MouseData) {
	if c.onClick != nil {
		c.checked = !c.checked
		c.onClick.Emit(c)
		c.MarkDirty()
	}
}
func (c *Checkbox) WantBlur() bool { return true }

func (c *Checkbox) Layout() {
	c.Drawable.Layout()
	w, h := c.rect.Size()
	str := ""
	if c.checked {
		str = "*"
	}
	lblCheck := NewText(str)
	margin := int(float64(c.Rect().GetLowestSize()) * 0.03)
	r := NewRect([]int{int(margin), int(margin), h - int(margin*2), h - int(margin*2)})
	lblCheck.SetRect(r)
	lblCheck.Bg(c.GetFg())
	lblCheck.Fg(c.GetBg())
	lblCheck.Layout()
	lblCheck.Draw(c.Image())

	lblTxt := NewText(c.txt)
	lblTxt.SetRect(NewRect([]int{h + int(margin), int(margin), (w - h) - int(margin*2), h - int(margin*2)}))
	lblTxt.Bg(c.GetBg())
	lblTxt.Fg(c.GetFg())
	lblTxt.Layout()
	lblTxt.Draw(c.Image())

	vector.StrokeRect(c.Image(), 0, 0, float32(w), float32(h), float32(margin), c.State().Color(), true)

	c.ClearDirty()
}

func (c *Checkbox) Draw(surface *ebiten.Image) {
	if c.IsHidden() {
		return
	}
	if c.IsDirty() {
		c.Layout()
	}
	c.Drawable.Draw(surface)
}
