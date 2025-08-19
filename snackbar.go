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
	msg   *Label
	timer *Timer
}

func NewSnackBar(message string) *SnackBar {
	s := &SnackBar{Drawable: NewDrawable()}
	s.msg = NewLabel(message)
	s.msg.SetBg(colornames.Blue)
	s.msg.SetFg(colornames.Yellow)
	s.SetRect(NewRect([]int{0, 0, 0, 0}))
	return s
}

func (s *SnackBar) Bg(value color.Color) *SnackBar {
	s.msg.SetBg(value)
	return s
}

func (s *SnackBar) Fg(value color.Color) *SnackBar {
	s.msg.SetFg(value)
	return s
}

func (s *SnackBar) SetText(value string) *SnackBar {
	s.msg.SetText(value)
	return s
}

func (s *SnackBar) ShowTime(durration time.Duration) *SnackBar {
	s.Show()
	s.timer = NewTimer(durration, func() {
		log.Println("Show snackbar done", s.msg.Text())
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

func (s *SnackBar) SetRect(r Rect[int]) {
	w0, h0 := GetUi().Size()
	rect := NewRect([]int{0, 0, w0, h0})
	sz := int(float64(rect.GetLowestSize()) * 0.1)
	x, y := rect.X+((w0-sz*8)/2), rect.Y+rect.Bottom()-sz*2
	w, h := w0-sz*8, sz
	s.msg.SetRect(NewRect([]int{x, y, w, h}))
}

func (s *SnackBar) Close() {
	s.timer = nil
	s.ImageReset()
}
