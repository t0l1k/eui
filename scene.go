package eui

import "log"

type Scene struct{ *Container }

func NewScene(lay Layouter) *Scene { return &Scene{Container: NewContainer(lay)} }
func (s *Scene) Entered()          {}
func (s *Scene) SetRect(rect Rect[int]) {
	s.Container.SetRect(rect)
	s.MarkDirty()
	log.Println("SceneBase:Resize", rect)
}
func (s *Scene) Quit() {
	for _, v := range s.Childrens() {
		v.Close()
	}
}
