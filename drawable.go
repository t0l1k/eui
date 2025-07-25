package eui

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Drawable struct {
	rect                    Rect[int]
	dirty, hidden, disabled bool
	image                   *ebiten.Image
	bg, fg                  color.Color
}

func NewDrawable() *Drawable             { return &Drawable{} }
func (s *Drawable) GetBg() color.Color   { return s.bg }
func (s *Drawable) Bg(value color.Color) { s.bg = value; s.MarkDirty() }
func (s *Drawable) GetFg() color.Color   { return s.fg }
func (s *Drawable) Fg(value color.Color) { s.fg = value; s.MarkDirty() }

func (s *Drawable) IsHidden() bool { return s.hidden }
func (s *Drawable) SetHidden(value bool) {
	s.hidden = value
	if !value {
		s.Enable()
	} else {
		s.Disable()
	}
	s.MarkDirty()
}
func (s *Drawable) IsDisabled() bool { return s.disabled }
func (s *Drawable) Enable()          { s.disabled = false }
func (s *Drawable) Disable()         { s.disabled = true }

func (s *Drawable) IsDirty() bool { return s.dirty }
func (s *Drawable) MarkDirty()    { s.dirty = true }
func (s *Drawable) ClearDirty()   { s.dirty = false }

func (c *Drawable) Traverse(action func(d Drawabler), reverse bool) { action(c) }

func (s *Drawable) Image() *ebiten.Image         { return s.image }
func (s *Drawable) SetImage(image *ebiten.Image) { s.image = image; s.MarkDirty() }
func (s *Drawable) ImageReset()                  { s.image = nil; s.MarkDirty() }
func (s *Drawable) Layout() {
	if s.Rect().IsEmpty() {
		panic("Drawable:Layout:Image:Rect.empty")
	}
	if s.Image() == nil {
		w0, h0 := s.Rect().Size()
		s.image = ebiten.NewImage(w0, h0)
		log.Println("Drawable:Layout:Image:nil", s.Rect())
	} else {
		s.image.Clear()
		// log.Println("Drawable:Layout:Image:clear", s.Rect())
	}
	s.ClearDirty()
}

func (s *Drawable) Update(dt int) {
	if s.IsDisabled() {
		return
	}
}
func (s *Drawable) Draw(surface *ebiten.Image) {
	if s.IsHidden() {
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

func (s *Drawable) Rect() Rect[int]        { return s.rect }
func (s *Drawable) SetRect(rect Rect[int]) { s.rect = rect; s.MarkDirty() }
func (s *Drawable) Resize(rect Rect[int])  { s.rect = rect; s.MarkDirty() }

func (s *Drawable) Close() {}
