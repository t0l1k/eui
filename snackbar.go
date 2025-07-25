package eui

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

type SnackBar struct {
	*Drawable
	msg   *Text
	timer *Timer
}

func NewSnackBar(message string) *SnackBar {
	s := &SnackBar{Drawable: NewDrawable()}
	s.msg = NewText(message)
	s.msg.Bg(colornames.Blue)
	s.msg.Fg(colornames.Yellow)
	s.Resize(NewRect([]int{0, 0, 0, 0}))
	return s
}

func (s *SnackBar) Bg(value color.Color) *SnackBar {
	s.msg.Bg(value)
	return s
}

func (s *SnackBar) Fg(value color.Color) *SnackBar {
	s.msg.Fg(value)
	return s
}

func (s *SnackBar) SetText(value string) *SnackBar {
	s.msg.SetText(value)
	return s
}

func (s *SnackBar) Show(durration time.Duration) *SnackBar {
	s.SetHidden(false)
	s.timer = NewTimer(durration, func() {
		log.Println("Show snackbar done", s.msg.GetText())
		GetUi().HideModal()
	}).On()
	GetUi().ShowModal(s)
	log.Println("Begin show snackbar", s)
	return s
}

func (s *SnackBar) Draw(surface *ebiten.Image) {
	if s.timer != nil {
		s.msg.Draw(surface)
	}
}

func (s *SnackBar) Resize(r Rect[int]) {
	w0, h0 := GetUi().Size()
	rect := NewRect([]int{0, 0, w0, h0})
	sz := int(float64(rect.GetLowestSize()) * 0.1)
	x, y := rect.X+((w0-sz*8)/2), rect.Y+rect.Bottom()-sz*2
	w, h := w0-sz*8, sz
	s.msg.Resize(NewRect([]int{x, y, w, h}))
}

func (s *SnackBar) Close() {
	s.timer = nil
	s.ImageReset()
}
