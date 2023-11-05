package main

import (
	"strconv"

	"github.com/t0l1k/eui"
)

type SceneTestListView struct {
	eui.SceneBase
	lst *eui.ListView
}

func NewSceneTestListView() *SceneTestListView {
	s := &SceneTestListView{}
	var list []string
	for i := 0; i < 22; i++ {
		list = append(list, "Item "+strconv.Itoa(i))
	}
	s.lst = eui.NewListView(list, 30)
	s.Add(s.lst)
	s.Resize()
	return s
}

func (s *SceneTestListView) Resize() {
	s.lst.Resize([]int{25, 25, 350, 350})
}

func NewGame() *eui.Ui {
	u := eui.GetUi()
	u.SetTitle("Test ListView")
	k := 2
	w, h := 200*k, 200*k
	u.SetSize(w, h)
	return u
}

func main() {
	eui.Init(NewGame())
	eui.Run(NewSceneTestListView())
	eui.Quit()
}
