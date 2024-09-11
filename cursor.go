package eui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui/colors"
)

type Cursor struct {
	DrawableBase
	show bool
}

func NewCursor(bg, fg color.Color) *Cursor {
	c := &Cursor{}
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
	s.Dirty = true
}

func (c *Cursor) Layout() {
	c.SpriteBase.Layout()
	c.Image().Fill(colors.SetInvisible(c.GetBg()))
	if c.show {
		c.Image().Fill(colors.Red)
	}
	c.Dirty = false
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
	if s.Dirty {
		s.Layout()
	}
	op := &ebiten.DrawImageOptions{}
	x, y := s.GetRect().Pos()
	op.GeoM.Translate(float64(x), float64(y))
	surface.DrawImage(s.Image(), op)
}

func (c *Cursor) Resize(rect []int) {
	c.Rect(NewRect(rect))
	c.ImageReset()
}
