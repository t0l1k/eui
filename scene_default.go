package eui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type SceneDefault struct {
	BoxLayout
}

func (s *SceneDefault) Update(dt int) {
	for _, v := range s.Container {
		v.Update(dt)
		vv, ok := v.(*BoxLayout)
		if ok {
			for _, v := range vv.Container {
				v.Update(dt)
			}
		}
	}
}

func (s *SceneDefault) Draw(surface *ebiten.Image) {
	for _, v := range s.Container {
		v.Draw(surface)
		vv, ok := v.(*BoxLayout)
		if ok {
			for _, v := range vv.Container {
				v.Draw(surface)
			}
		}
	}
}

func (s *SceneDefault) Entered() {}

func (s *SceneDefault) Resize() {
	w0, h0 := GetUi().Size()
	s.BoxLayout.Resize([]int{0, 0, w0, h0})
}

func (s *SceneDefault) Quit() {
	for _, v := range s.Container {
		v.Close()
		vv, ok := v.(*BoxLayout)
		if ok {
			for _, v := range vv.Container {
				v.Close()
			}
		}
	}
}
