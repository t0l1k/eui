package eui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

type Cursor struct {
	*Drawable
	show bool
}

func NewCursor(bg, fg color.Color) *Cursor {
	c := &Cursor{Drawable: NewDrawable()}
	c.Bg(bg)
	c.Fg(fg)
	return c
}

func (s *Cursor) IsVisible() bool { return s.show }
func (s *Cursor) Visible(value bool) {
	if s.show == value {
		return
	}
	s.show = value
	s.MarkDirty()
}

func (c *Cursor) Layout() {
	c.Drawable.Layout()
	c.Image().Fill(color.Transparent)
	if c.show {
		c.Image().Fill(colornames.Red)
	}
	c.ClearDirty()
}

func (s *Cursor) Update(dt int) {
	if s.IsDisabled() {
		return
	}
}
func (s *Cursor) Draw(surface *ebiten.Image) {
	if s.IsDisabled() {
		return
	}
	if s.IsDirty() {
		s.Layout()
	}
	op := &ebiten.DrawImageOptions{}
	x, y := s.Rect().Pos()
	op.GeoM.Translate(float64(x), float64(y))
	surface.DrawImage(s.Image(), op)
}

func (c *Cursor) Resize(rect Rect[int]) {
	c.SetRect(rect)
	c.ImageReset()
}
