package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/t0l1k/eui"
)

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
		gridCont := eui.NewContainer(eui.NewGridLayout(2, 2, 5))
		gridCont.Add(eui.NewLabel("1" + title).SetFontSize(12))
		gridCont.Add(eui.NewLabel("1" + title + " 2" + title).SetFontSize(18))
		gridCont.Add(eui.NewLabel("1" + title + " 2" + title + " 3" + title).SetFontSize(24))
		gridCont.Add(eui.NewLabel("1" + title + " 2" + title + " 3" + title + " 4" + title).SetFontSize(30))
		hAlign := []string{"Left", "Center", "Right"}
		vAlign := []string{"Up", "Center", "Down"}
		i := 0
		btn := eui.NewButton("Next align", func(b *eui.Button) {
			x, y := i%3, i/3
			for _, lbl := range gridCont.Children() {
				lbl.(*eui.Label).SetAlign(text.Align(x), text.Align(y))
			}
			a.Children()[2].(*eui.Label).SetText(hAlign[x] + " " + vAlign[y])
			log.Println("align", i, x, y, hAlign[x], vAlign[y])
			i++
			if i > 8 {
				i = 0
			}
		})

		a.Add(eui.NewTopBar(title, nil))
		a.Add(gridCont)
		a.Add(eui.NewLabel(hAlign[1] + " " + vAlign[1]))
		a.Add(btn)
		a.Add(eui.NewGridBackground(50))
		return a
	}())
	eui.Quit(func() {})
}
