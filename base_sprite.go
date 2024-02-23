package eui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type SpriteBase struct {
	LayoutBase
	Visible, disabled bool
	image             *ebiten.Image
	bg, fg            color.Color
}

func (s *SpriteBase) GetBg() color.Color   { return s.bg }
func (s *SpriteBase) Bg(value color.Color) { s.bg = value; s.Dirty = true }
func (s *SpriteBase) GetFg() color.Color   { return s.fg }
func (s *SpriteBase) Fg(value color.Color) { s.fg = value; s.Dirty = true }

func (s *SpriteBase) IsDisabled() bool { return !s.disabled }
func (s *SpriteBase) Enable()          { s.disabled = true }
func (s *SpriteBase) Disable()         { s.disabled = false }

func (s *SpriteBase) Image() *ebiten.Image         { return s.image }
func (s *SpriteBase) SetImage(image *ebiten.Image) { s.image = image }
func (s *SpriteBase) ImageReset()                  { s.image = nil; s.Dirty = true }
func (s *SpriteBase) Layout() {
	if s.Image() == nil {
		w0, h0 := s.GetRect().Size()
		s.image = ebiten.NewImage(w0, h0)
	} else {
		s.image.Clear()
	}
	s.Dirty = false
}

func (s *SpriteBase) Update(dt int)              {}
func (d *SpriteBase) Draw(surface *ebiten.Image) {}