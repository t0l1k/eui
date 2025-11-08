package eui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/t0l1k/eui/utils"
	"golang.org/x/image/colornames"
)

type Slider struct {
	*Drawable
	min, max               float64
	value                  *Signal[float64]
	orientation            Orientation
	OnChange               SlotFunc[float64]
	thumbOn                bool
	thumbSize              float64
	valueColor, thumbColor color.Color
	trackRect              Rect[int]
}

// Умею показать прогресс от 0 до 1, горизонтально/вертикально с подпиской при обновлении значения
func NewProgress(min, max, initial float64, o Orientation, fn SlotFunc[float64]) *Slider {
	p := &Slider{Drawable: NewDrawable(), min: min, max: max, value: NewSignal(func(a, b float64) bool { return a == b }), orientation: o, thumbOn: false}
	p.value.Emit(initial)
	p.value.Connect(fn)
	p.SetBg(colornames.Navy)
	p.SetFg(colornames.Orange)
	p.valueColor = colornames.Fuchsia
	p.thumbColor = colornames.Teal
	return p
}

func NewSlider(min, max, initial float64, o Orientation, fn SlotFunc[float64]) *Slider {
	s := NewProgress(min, max, initial, o, fn)
	s.thumbOn = true
	return s
}
func (s *Slider) Hit(pt Point[int]) Drawabler {
	if !pt.In(s.rect) || s.IsHidden() || s.IsDisabled() {
		return nil
	}
	return s
}
func (s *Slider) WantBlur() bool { return true }
func (b *Slider) MouseDown(md MouseData) {
	if !b.thumbOn {
		return
	}
	b.pressed = true
	b.updateValue(md.pos)
}
func (b *Slider) MouseUp(md MouseData) {
	b.pressed = false
}
func (b *Slider) MouseDrag(md MouseData) {
	if !b.thumbOn {
		return
	}
	if b.pressed {
		b.updateValue(md.pos)
	}
}
func (s *Slider) Value() float64 { return s.value.Value() }
func (s *Slider) SetValue(value float64) *Slider {
	if value > 1.0 {
		s.value.Emit(1.0)
		s.MarkDirty()
		return s
	}
	s.value.Emit(utils.Clamp(value, s.min, s.max))
	s.MarkDirty()
	return s
}
func (s *Slider) updateValue(pt Point[int]) {
	var value float64
	x, y := s.trackRect.Pos()
	w, h := s.trackRect.Size()
	if s.orientation == Horizontal {
		relX := float64(pt.X - x)
		value = relX / float64(w)
	} else {
		relY := float64(pt.Y - y)
		value = 1.0 - (relY / float64(h))
	}
	s.SetValue(utils.Clamp(value, s.min, s.max))
}

func (p *Slider) Layout() {
	p.Drawable.Layout()
	p.Image().Fill(p.bg)
	w0, h0 := p.rect.Width(), p.rect.Height()
	w, h := p.trackRect.Width(), p.trackRect.Height()
	vector.StrokeRect(p.image, 0, 0, float32(w0), float32(h0), 1, p.Fg(), true)

	if p.orientation == Horizontal {
		valueWidth := float64(w) * p.Value()
		vector.DrawFilledRect(p.Image(), float32(p.thumbSize/2), float32(h0/3), float32(w), float32(h0/3), p.fg, true)
		vector.DrawFilledRect(p.Image(), float32(p.thumbSize/2), float32(h0/3), float32(valueWidth), float32(h0/3), p.valueColor, true)
		if p.thumbOn {
			thumbX := valueWidth
			vector.DrawFilledRect(p.Image(), float32(thumbX), 0, float32(p.thumbSize), float32(h0), p.thumbColor, true)
		}
	} else {
		valueHeight := float64(h) * p.Value()
		vector.DrawFilledRect(p.Image(), float32(w0/3), float32(p.thumbSize/2), float32(w0/3), float32(h), p.fg, true)
		vector.DrawFilledRect(p.Image(), float32(w0/3), float32(h0)-float32(valueHeight)-float32(p.thumbSize)/2, float32(w0/3), float32(valueHeight), p.valueColor, true)

		if p.thumbOn {
			thumbY := float64(p.Rect().Height()) - valueHeight - p.thumbSize
			vector.DrawFilledRect(p.Image(), 0, float32(thumbY), float32(h0), float32(p.thumbSize), p.thumbColor, true)
		}
	}
	p.ClearDirty()
}

func (s *Slider) Draw(surface *ebiten.Image) {
	if s.IsHidden() {
		return
	}
	if s.IsDirty() {
		s.Layout()
	}
	s.Drawable.Draw(surface)
}

func (s *Slider) SetRect(rect Rect[int]) {
	s.Drawable.SetRect(rect)
	w, h := rect.Size()
	x, y := rect.Pos()
	switch s.orientation {
	case Horizontal:
		s.thumbSize = float64(h) / 3
		x += int(s.thumbSize / 2)
		w -= int(s.thumbSize)
	case Vertical:
		s.thumbSize = float64(w) / 3
		y += int(s.thumbSize / 2)
		h -= int(s.thumbSize)
	}
	s.trackRect = NewRect([]int{x, y, w, h})
}
