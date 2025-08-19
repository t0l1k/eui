package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/t0l1k/eui"
)

type SceneTestListView struct {
	*eui.Scene
	lstText, lstButtons, lstCheckBoxs *eui.ListView
	btnRemoveSelected                 *eui.Button
}

func NewSceneTestListView() *SceneTestListView {
	s := &SceneTestListView{Scene: eui.NewScene(eui.NewAbsoluteLayout())}
	var list []string
	for i := 0; i < 54; i++ {
		list = append(list, "Item "+strconv.Itoa(i))
	}
	theme := eui.GetUi().Theme()
	bg := theme.Get(eui.ListViewItemBg)
	fg := theme.Get(eui.ListViewItemFg)
	s.lstText = eui.NewListView()
	s.lstText.SetupListViewText(list, 30, 2, bg, fg)
	s.Add(s.lstText)

	s.lstButtons = eui.NewListView()
	s.lstButtons.SetupListViewButtons(list, 30, 1, bg, fg, func(b *eui.Button) {
		log.Println("pressed:", b.Text())
	})
	s.Add(s.lstButtons)

	s.lstCheckBoxs = eui.NewListView()
	s.lstCheckBoxs.SetupListViewCheckBoxs(list, 30, 1, bg, fg, func(b *eui.Checkbox) {
		log.Println("pressed:", b.Text())
	})
	s.Add(s.lstCheckBoxs)

	s.btnRemoveSelected = eui.NewButton("Remove Selected", func(b *eui.Button) {
		list = nil
		for _, v := range s.lstCheckBoxs.GetCheckBoxes() {
			if v.IsChecked() {
				fmt.Println("selected:", v.Text())
				continue
			}
			list = append(list, v.Text())
		}
		s.lstCheckBoxs.Reset()
		s.lstCheckBoxs.SetupListViewCheckBoxs(list, 30, 1, bg, fg, func(b *eui.Checkbox) {
			log.Println("pressed:", b.Text())
		})

	})
	s.Add(s.btnRemoveSelected)
	return s
}

func (s *SceneTestListView) SetRect(rect eui.Rect[int]) {
	_, h0 := rect.Size()
	margin := int(float64(rect.GetLowestSize()) * 0.1)
	x, y := margin, margin
	w, h := margin*3, h0-margin*2
	s.lstButtons.SetRect(eui.NewRect([]int{x, y, w, h}))
	x += margin*3 + margin
	s.lstText.SetRect(eui.NewRect([]int{x, y, w, h}))
	x += margin*3 + margin
	h -= margin * 3
	s.lstCheckBoxs.SetRect(eui.NewRect([]int{x, y, w, h}))
	y += h + margin
	h = margin
	s.btnRemoveSelected.SetRect(eui.NewRect([]int{x, y, w, h}))
	log.Println("SceneTestListView:Resize:", rect, s.lstButtons.Rect(), s.lstCheckBoxs.Rect(), s.lstText.Rect(), s.btnRemoveSelected.Rect())
}

func NewGame() *eui.Ui {
	k := 2
	w, h := 500*k, 200*k
	return eui.GetUi().SetTitle("Test ListView").SetSize(w, h)
}

func main() {
	eui.Init(NewGame())
	eui.Run(NewSceneTestListView())
	eui.Quit(func() {})
}
