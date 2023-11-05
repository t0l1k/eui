package eui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type SceneBase struct {
	BoxLayout
}

func (s *SceneBase) Update(dt int) {
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

func (s *SceneBase) Draw(surface *ebiten.Image) {
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

func (s *SceneBase) Entered() {}

func (s *SceneBase) Resize() {
	w0, h0 := GetUi().Size()
	s.BoxLayout.Resize([]int{0, 0, w0, h0})
}

func (s *SceneBase) Quit() {
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
