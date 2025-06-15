package eui

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

type SnackBar struct {
	DrawableBase
	msg   *Text
	timer *Timer
}

func NewSnackBar(message string) *SnackBar {
	s := &SnackBar{}
	s.msg = NewText(message)
	s.msg.Bg(colornames.Blue)
	s.msg.Fg(colornames.Yellow)
	s.Add(s.msg)
	s.Resize([]int{})
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

func (s *SnackBar) Show(durration int) *SnackBar {
	s.Visible(true)
	s.timer = NewTimer(durration).On()
	log.Println("Begin show snackbar", s)
	return s
}

func (s *SnackBar) Update(dt int) {
	if s.timer == nil {
		return
	}
	if !s.timer.IsDone() {
		s.timer.Update(dt)
	} else {
		s.Visible(false)
		log.Println("Show snackbar done", s.msg.GetText())
		s.Close()
	}
}

func (s *SnackBar) Draw(surface *ebiten.Image) {
	if s.timer != nil {
		s.msg.Draw(surface)
	}
}

func (s *SnackBar) Resize(r []int) {
	w0, h0 := GetUi().Size()
	rect := NewRect([]int{0, 0, w0, h0})
	sz := int(float64(rect.GetLowestSize()) * 0.1)
	x, y := rect.X+((w0-sz*8)/2), rect.Y+rect.Bottom()-sz*2
	w, h := w0-sz*8, sz
	s.msg.Resize([]int{x, y, w, h})
}

func (s *SnackBar) Close() {
	s.timer = nil
	s.msg = nil
	s.ImageReset()
	s.ResetContainerBase()
	s = nil
}
