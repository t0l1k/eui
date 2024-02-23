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
	for i := 0; i < 54; i++ {
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
	w0, h0 := eui.GetUi().Size()
	rect := eui.NewRect([]int{0, 0, w0, h0})
	margin := int(float64(rect.GetLowestSize()) * 0.1)
	x, y := margin, margin
	w, h := margin*3, h0-margin*2
	s.lstButtons.Resize([]int{x, y, w, h})
	x += margin*3 + margin
	s.lstText.Resize([]int{x, y, w, h})
	x += margin*3 + margin
	h -= margin * 3
	s.lstCheckBoxs.Resize([]int{x, y, w, h})
	y += h + margin
	h = margin
	s.btnRemoveSelected.Resize([]int{x, y, w, h})
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
