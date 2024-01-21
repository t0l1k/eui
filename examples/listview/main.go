package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/t0l1k/eui"
)

type SceneTestListView struct {
	eui.SceneBase
	lstText, lstButtons, lstCheckBoxs *eui.ListView
	btnRemoveSelected                 *eui.Button
}

func NewSceneTestListView() *SceneTestListView {
	s := &SceneTestListView{}
	var list []string
	for i := 0; i < 24; i++ {
		list = append(list, "Item "+strconv.Itoa(i))
	}
	theme := eui.GetUi().GetTheme()
	bg := theme.Get(eui.ListViewItemBg)
	fg := theme.Get(eui.ListViewItemFg)
	s.lstText = eui.NewListView()
	s.lstText.SetupListViewText(list, 30, 2, bg, fg)
	s.Add(s.lstText)

	s.lstButtons = eui.NewListView()
	s.lstButtons.SetupListViewButtons(list, 30, 1, bg, fg, func(b *eui.Button) {
		log.Println("pressed:", b.GetText())
	})
	s.Add(s.lstButtons)

	s.lstCheckBoxs = eui.NewListView()
	s.lstCheckBoxs.SetupListViewCheckBoxs(list, 30, 1, bg, fg, func(b *eui.Checkbox) {
		log.Println("pressed:", b.GetText())
	})
	s.Add(s.lstCheckBoxs)

	s.btnRemoveSelected = eui.NewButton("Remove Selected", func(b *eui.Button) {
		list = nil
		for _, v := range s.lstCheckBoxs.GetCheckBoxes() {
			if v.IsChecked() {
				fmt.Println("selected:", v.GetText())
				continue
			}
			list = append(list, v.GetText())
		}
		s.lstCheckBoxs.Reset()
		s.lstCheckBoxs.SetupListViewCheckBoxs(list, 30, 1, bg, fg, func(b *eui.Checkbox) {
			log.Println("pressed:", b.GetText())
		})

	})
	s.Add(s.btnRemoveSelected)

	s.Resize()
	return s
}

func (s *SceneTestListView) Resize() {
	s.lstButtons.Resize([]int{25, 25, 250, 250})
	s.lstText.Resize([]int{300, 25, 250, 250})
	s.lstCheckBoxs.Resize([]int{600, 25, 250, 250})
	s.btnRemoveSelected.Resize([]int{600, 300, 250, 50})
}

func NewGame() *eui.Ui {
	u := eui.GetUi()
	u.SetTitle("Test ListView")
	k := 2
	w, h := 500*k, 200*k
	u.SetSize(w, h)
	return u
}

func main() {
	eui.Init(NewGame())
	eui.Run(NewSceneTestListView())
	eui.Quit()
}
