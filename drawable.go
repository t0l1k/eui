package eui

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Drawable struct {
	state          ViewState
	rect           Rect[int]
	dirty, pressed bool
	image          *ebiten.Image
	bg, fg         color.Color
	shadowSize     int
	shadowColor    color.Color
}

func NewDrawable() *Drawable {
	return &Drawable{state: StateNormal, dirty: true, bg: color.Transparent, fg: color.Black, shadowSize: 0, shadowColor: color.RGBA{0, 0, 0, 128}}
}

func (s *Drawable) State() ViewState         { return s.state }
func (s *Drawable) SetState(value ViewState) { s.state = value; s.MarkDirty() }

func (s *Drawable) GetBg() color.Color         { return s.bg }
func (s *Drawable) Bg(value color.Color)       { s.bg = value; s.MarkDirty() }
func (s *Drawable) GetFg() color.Color         { return s.fg }
func (s *Drawable) Fg(value color.Color)       { s.fg = value; s.MarkDirty() }
func (s *Drawable) Shadow() (int, color.Color) { return s.shadowSize, s.shadowColor }
func (s *Drawable) SetShadow(size int, color color.Color) {
	s.shadowSize = size
	s.shadowColor = color
	s.MarkDirty()
}

func (s *Drawable) IsHidden() bool { return s.state.IsHidden() }
func (s *Drawable) Show()          { s.SetState(StateNormal); s.MarkDirty() }
func (s *Drawable) Hide()          { s.SetState(StateHidden); s.ImageReset() }

func (s *Drawable) IsDisabled() bool { return s.state.IsDisabled() }
func (s *Drawable) Enable()          { s.SetState(StateNormal); s.MarkDirty() }
func (s *Drawable) Disable()         { s.SetState(StateDisabled); s.MarkDirty() }

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
	w, h := s.Rect().Size()
	if s.Image() == nil {
		s.image = ebiten.NewImage(w, h)
		log.Println("Drawable:Layout:Image:nil", s.Rect())
	} else {
		s.image.Clear()
		// log.Println("Drawable:Layout:Image:clear", s.Rect())
	}
	s.Image().Fill(s.bg)
	x, y := 0, 0
	if s.shadowSize > 0 {
		off := s.shadowSize / 2
		if s.pressed {
			x += off
			y += off
		}
		vector.DrawFilledRect(s.Image(), float32(x), float32(y), float32(w), float32(h), s.shadowColor, true)
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
func (s *Drawable) SetRect(rect Rect[int]) { s.rect = rect; s.ImageReset() }

func (s *Drawable) Close() {}
