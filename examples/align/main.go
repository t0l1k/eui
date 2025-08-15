package main

import "github.com/t0l1k/eui"

const title = "Text Align Demo"

func main() {
	eui.Init(func() *eui.Ui {
		u := eui.GetUi()
		u.SetTitle(title)
		k := 2
		w, h := 320*k, 200*k
		u.SetSize(w, h)
		return u
	}())
	eui.Run(func() *eui.Scene {
		a := eui.NewScene(eui.NewLayoutVerticalPercent([]int{5, 70, 15, 10}, 5))
		a.Add(eui.NewTopBar(title, nil))
		gridCont := eui.NewContainer(eui.NewGridLayout(2, 2, 5))
		gridCont.Add(eui.NewLabel("1" + title).SetFontSize(12))
		gridCont.Add(eui.NewLabel("1" + title + " 2" + title).SetFontSize(18))
		gridCont.Add(eui.NewLabel("1" + title + " 2" + title + " 3" + title).SetFontSize(24))
		gridCont.Add(eui.NewLabel("1" + title + " 2" + title + " 3" + title + " 4" + title).SetFontSize(30))
		a.Add(gridCont)
		a.Add(eui.NewLabel(eui.LabelAlignCenter.String()))
		i := 0
		a.Add(eui.NewButton("Next align", func(b *eui.Button) {
			for _, lbl := range gridCont.Children() {
				lbl.(*eui.Label).SetAlign(eui.LabelAlign(i))
			}
			a.Children()[2].(*eui.Label).SetText(eui.LabelAlign(i).String())
			i++
			if i > 5 {
				i = 0
			}
		}))
		return a
	}())
	eui.Quit(func() {})
}
