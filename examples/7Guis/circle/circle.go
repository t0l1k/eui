package circle

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/t0l1k/eui"
	"golang.org/x/image/colornames"
)

type CircleModel struct{ x, y, r int }

type Circle struct {
	*eui.Drawable
	CircleModel
	hover   bool
	changed *eui.Signal[*Circle]
	state   eui.ViewState
}

func NewCircle(fn eui.SlotFunc[*Circle]) *Circle {
	c := &Circle{Drawable: eui.NewDrawable()}
	c.SetBg(color.Transparent)
	c.SetFg(colornames.Yellow)
	c.changed = eui.NewSignal(func(a, b *Circle) bool { return false })
	c.changed.Connect(fn)
	return c
}
func (c *Circle) Reset(x, y, r int) {
	c.x = x
	c.y = y
	c.r = r
	c.SetRect(eui.NewRect([]int{x - r, y - r, r * 2, r * 2}))
	c.MarkDirty()
}
func (b *Circle) WantBlur() bool { return true }
func (b *Circle) Hit(pt eui.Point[int]) eui.Drawabler {
	if !pt.In(b.Rect()) || b.IsHidden() || b.IsDisabled() {
		if b.hover {
			b.hover = false
			b.MarkDirty()
			return b
		}
		return nil
	}
	if !b.hover {
		b.hover = true
		b.MarkDirty()
		return b
	}
	return b
}
func (c *Circle) MouseUp(md eui.MouseData) { eui.GetUi().ShowModal(NewEditCircleDialog(c)) }
func (c *Circle) Layout() {
	c.Drawable.Layout()
	x := float32(c.x - c.Rect().X)
	y := float32(c.y - c.Rect().Y)
	r := float32(c.r)
	if c.hover {
		vector.FillCircle(c.Image(), x, y, r, colornames.Aqua, true)
	}
	vector.StrokeCircle(c.Image(), x, y, r, 1, c.Fg(), true)
	c.ClearDirty()
}
func (c *Circle) Draw(surface *ebiten.Image) {
	if c.IsHidden() {
		return
	}
	if c.IsDirty() {
		c.Layout()
	}
	c.Drawable.Draw(surface)
}
func (c *Circle) Model() CircleModel { return CircleModel{c.x, c.y, c.r} }
func (c *Circle) String() string     { return fmt.Sprintf("Circle[%v,%v,%v,%v]", c.x, c.y, c.r, c.state) }
