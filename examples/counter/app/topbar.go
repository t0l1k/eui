package app

import (
	"github.com/hajimehoshi/ebiten/v2"
	ui "github.com/t0l1k/eui"
)

type TopBar struct {
	btnQuit  *ui.Button
	lblTitle *ui.Label
	ui.ContainerDefault
}

func NewTopBar() *TopBar {
	tb := &TopBar{}
	rect := []int{0, 0, 1, 1}
	bg := ui.Yellow
	fg := ui.Navy
	tb.lblTitle = ui.NewLabel(ui.GetUi().GetTitle(), rect, bg, fg)
	tb.Add(tb.lblTitle)
	tb.btnQuit = ui.NewButton("<", rect, ui.GreenYellow, ui.Black, func(b *ui.Button) { ui.Pop() })
	tb.Add(tb.btnQuit)
	return tb
}

func (tb *TopBar) Update(dt int) {
	for _, v := range tb.Container {
		v.Update(dt)
	}
}

func (tb *TopBar) Draw(surface *ebiten.Image) {
	for _, v := range tb.Container {
		v.Draw(surface)
	}
}

func (tb *TopBar) Resize() {
	w, h := ebiten.WindowSize()
	rect := ui.NewRect([]int{0, 0, w, h})
	x, y, w, h := 0, 0, int(float64(rect.H)*0.05), int(float64(rect.H)*0.05)
	tb.btnQuit.Resize([]int{x, y, w, h})
	x, w = h, h*5
	tb.lblTitle.Resize([]int{x, y, w, h})
}

func (tb *TopBar) Close() {
	for _, v := range tb.Container {
		v.Close()
	}
}
