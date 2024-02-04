package eui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type SceneBase struct {
	ContainerBase
}

func (s *SceneBase) Update(dt int) {
	for _, v := range s.GetContainer() {
		v.Update(dt)
	}
}

func (s *SceneBase) Draw(surface *ebiten.Image) {
	for _, v := range s.GetContainer() {
		v.Draw(surface)
	}
}

func (s *SceneBase) Entered() {}

func (s *SceneBase) Resize() {
	w0, h0 := GetUi().Size()
	for _, v := range s.GetContainer() {
		v.Resize([]int{0, 0, w0, h0})
	}
}

func (s *SceneBase) Quit() {
	for _, v := range s.GetContainer() {
		v.Close()
	}
}
