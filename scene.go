package eui

import "log"

type Scene struct{ *Container }

func NewScene(lay Layouter) *Scene { return &Scene{Container: NewContainer(lay)} }
func (s *Scene) Entered()          {}
func (s *Scene) Resize() {
	w0, h0 := GetUi().Size()
	rect := NewRect([]int{0, 0, w0, h0})
	s.SetRect(rect)
	s.MarkDirty()
	log.Println("SceneBase:Resize", rect)
}
func (s *Scene) Quit() {
	for _, v := range s.Childrens() {
		v.Close()
	}
}
