package eui

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

type SnackBar struct {
	*Container
	msg     *Label
	quitBtn *Button
	timer   *Timer
}

func NewSnackBar(message string) *SnackBar {
	s := &SnackBar{Container: NewContainer(NewLayoutHorizontalPercent([]int{10, 90}, 3))}
	s.quitBtn = NewButton("X", func(b *Button) {
		GetUi().HideModal()
		log.Println("SnackBar:Close")
	})
	s.msg = NewLabel(message)
	s.msg.SetBg(colornames.Blue)
	s.msg.SetFg(colornames.Yellow)
	s.SetViewType(ViewModal)
	s.Add(s.quitBtn)
	s.Add(s.msg)
	return s
}

func (b *SnackBar) Hit(pt Point[int]) Drawabler {
	if !pt.In(b.rect) || b.IsHidden() {
		return nil
	}
	if pt.In(b.quitBtn.Rect()) {
		log.Println("SnackBar:Hit:", b.Rect(), b.msg.Text())
		return b.quitBtn
	}
	return nil
}

func (s *SnackBar) SetBg(value color.Color) *SnackBar { s.msg.SetBg(value); return s }
func (s *SnackBar) SetFg(value color.Color) *SnackBar { s.msg.SetFg(value); return s }
func (s *SnackBar) SetText(value string) *SnackBar    { s.msg.SetText(value); return s }
func (s *SnackBar) ShowTime(durration time.Duration) *SnackBar {
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
		s.Container.Draw(surface)
	}
}
func (s *SnackBar) Close() { s.timer = nil; s.ImageReset() }
