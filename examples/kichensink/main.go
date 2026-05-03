package main

import (
	"fmt"
	"image"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/kichensink/res"
)

func NewTestScene() *eui.Scene {
	sc := eui.NewScene(eui.NewAbsoluteLayout())

	btn1 := eui.NewButton("Alert", func(b *eui.Button) {
		eui.NewSnackBar("Hello, world!").ShowTime(3 * time.Second)
	})
	btn1.SetRect(eui.NewRect([]int{10, 450 - 40, 80, 30}))
	sc.Add(btn1)

	slider := eui.NewSlider(0, 100, 0.5, eui.Horizontal, func(data float64) {})
	slider.SetRect(eui.NewRect([]int{10, 150, 100, 20}))
	sc.Add(slider)

	progress := eui.NewProgress(0, 1, 0.25, eui.Horizontal, func(data float64) {})
	progress.SetRect(eui.NewRect([]int{10, 180, 100, 20}))
	sc.Add(progress)

	var runningTask bool
	btn2 := eui.NewButton("Start", func(b *eui.Button) {
		progress.SetValue(0)
		runningTask = true
		go func() {
			for runningTask && progress.Value() < 1.0 {
				time.Sleep(100 * time.Millisecond)
				progress.SetValue(progress.Value() + 0.1)
			}
			if runningTask {
				eui.NewSnackBar("Simulated task finished!").ShowTime(3 * time.Second)
				runningTask = false
			}
		}()
	})
	btn2.SetRect(eui.NewRect([]int{10, 210, 80, 30}))
	sc.Add(btn2)

	items := make([]string, 100)
	for i := range items {
		items[i] = fmt.Sprintf("Item %d", i)
	}
	list := eui.NewListView()
	list.SetupListViewButtons(items, 30, 1, eui.GetUi().Theme().Get(eui.ButtonBg), eui.GetUi().Theme().Get(eui.ButtonFg), func(b *eui.Button) {
		eui.NewSnackBar(fmt.Sprintf("Selected: %s", b.Text())).ShowTime(3 * time.Second)
	})
	list.SetRect(eui.NewRect([]int{450 / 2, 10, 150, 170}))
	sc.Add(list)

	cb1 := eui.NewCheckbox("I eat food", func(c *eui.Checkbox) {})
	cb1.SetRect(eui.NewRect([]int{list.Rect().X, list.Rect().Y + list.Rect().Height() + 50, 200, 20}))
	cb1.SetChecked(true)
	sc.Add(cb1)

	cb2 := eui.NewCheckbox("I drink water", func(c *eui.Checkbox) {})
	cb2.SetRect(eui.NewRect([]int{cb1.Rect().X, cb1.Rect().Y + cb1.Rect().Height() + 10, 200, 20}))
	cb2.SetChecked(true)
	sc.Add(cb2)

	cb3 := eui.NewCheckbox("I exercise regularly", func(c *eui.Checkbox) {})
	cb3.SetRect(eui.NewRect([]int{cb1.Rect().X, cb2.Rect().Y + cb2.Rect().Height() + 10, 200, 20}))
	sc.Add(cb3)

	tf1 := eui.NewInputLine(func(il *eui.InputLine) {
		eui.NewSnackBar(fmt.Sprintf("You entered: %s", il.Text())).ShowTime(3 * time.Second)
	})
	tf1.SetPlaceholder("Type text and press Enter")
	tf1.SetRect(eui.NewRect([]int{10, 10, 150, 30}))
	sc.Add(tf1)

	tf2 := eui.NewInputLine(func(il *eui.InputLine) {})
	tf2.SetPlaceholder("Another input")
	tf2.SetRect(eui.NewRect([]int{10, 50, 150, 30}))
	sc.Add(tf2)

	logoIcon := eui.NewIcon(eui.GetUi().RM().LoadImage(res.LogoPNG))
	logoIcon.SetRect(eui.NewRect([]int{450 - 90, 10, 80, 40}))
	sc.Add(logoIcon)

	spinner := eui.GetUi().RM().LoadImage(res.SpinnerPNG)
	loadSpinner := func() []*ebiten.Image {
		var icons []*ebiten.Image // x:0 y:24 w: 26 h: 26
		var nImg *ebiten.Image
		w, h := 26, 26
		for i := 0; i < 5; i++ {
			nImg = ebiten.NewImage(w, h)
			x, y := i*w, 24
			x += i
			nImg.DrawImage(spinner.SubImage(image.Rect(x, y, x+w, y+h)).(*ebiten.Image), &ebiten.DrawImageOptions{})
			icons = append(icons, nImg)
		}
		return icons
	}
	spinnerIcon := eui.NewAnimation(loadSpinner(), 100*time.Millisecond)
	spinnerIcon.SetPos(eui.NewPoint[float64](450-74, 600-74))
	// spinnerIcon.SetRect(eui.NewRect([]int{450 - 74, 600 - 74, 64, 64}))
	sc.Add(spinnerIcon)

	return sc
}

func main() {
	title := "Витрина элементов"
	eui.Init(eui.GetUi().SetTitle(title).SetSize(800, 600))
	eui.Run(NewTestScene())
}
