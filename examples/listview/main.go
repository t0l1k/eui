package main

import (
	"log"
	"strconv"

	"github.com/t0l1k/eui"
)

type SceneTestListView struct {
	eui.SceneBase
	lstText, lstButtons *eui.ListView
}

func NewSceneTestListView() *SceneTestListView {
	s := &SceneTestListView{}
	var list []string
	for i := 0; i < 62; i++ {
		list = append(list, "Item "+strconv.Itoa(i))
	}
	theme := eui.GetUi().GetTheme()
	bg := theme.Get(eui.ListViewItemBg)
	fg := theme.Get(eui.ListViewItemFg)
	s.lstText = eui.NewListView()
	s.lstText.SetupListViewText(list, 30, 2, bg, fg)
	s.Add(s.lstText)

	s.lstButtons = eui.NewListView()
	s.lstButtons.SetupListViewButtons(list, 30, 2, bg, fg, func(b *eui.Button) {
		log.Println("pressed:", b.GetText())
	})
	s.Add(s.lstButtons)

	s.Resize()
	return s
}

func (s *SceneTestListView) Resize() {
	s.lstText.Resize([]int{25, 25, 350, 350})
	s.lstButtons.Resize([]int{400, 25, 350, 350})
}

func NewGame() *eui.Ui {
	u := eui.GetUi()
	u.SetTitle("Test ListView")
	k := 2
	w, h := 400*k, 200*k
	u.SetSize(w, h)
	return u
}

func main() {
	eui.Init(NewGame())
	eui.Run(NewSceneTestListView())
	eui.Quit()
}
