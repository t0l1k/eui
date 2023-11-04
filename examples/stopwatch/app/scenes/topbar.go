package scenes

import (
	"fmt"

	"github.com/t0l1k/eui"
)

type TopBar struct {
	eui.View
	btnQuit         *eui.Button
	lblTitle, lblTm *eui.Text
	tm              *eui.StringVar
	Stopwatch       *eui.Stopwatch
}

func NewTopBar(title string) *TopBar {
	t := &TopBar{}
	t.Stopwatch = eui.NewStopwatch()
	fg := eui.Yellow
	bg := eui.Blue
	bg2 := eui.Gray
	t.SetupView(bg2)
	sq := "<"
	if eui.GetUi().IsMainScene() {
		sq = "x"
	}
	t.btnQuit = eui.NewButton(sq, fg, bg, func(b *eui.Button) {
		eui.GetUi().Pop()
	})
	t.Add(t.btnQuit)
	t.lblTitle = eui.NewText(title, bg, fg)
	t.Add(t.lblTitle)
	t.tm = eui.NewStringVar(t.Stopwatch.String())
	t.lblTm = eui.NewText(t.tm.Get(), bg, fg)
	t.tm.Attach(t.lblTm)
	t.Add(t.lblTm)
	t.Stopwatch.Start()
	return t
}

func (t *TopBar) Update(dt int) {
	t.View.Update(dt)
	str := fmt.Sprintf("%v", t.Stopwatch.String2())
	t.tm.Set(str)
}

func (t *TopBar) Resize(arr []int) {
	t.View.Resize(arr)
	x, y, w, h := 0, 0, t.GetRect().H, t.GetRect().H
	t.btnQuit.Resize([]int{x, y, w, h})
	x += h
	w = h * 3
	t.lblTitle.Resize([]int{x, y, w, h})
	x = t.GetRect().W - w
	t.lblTm.Resize([]int{x, y, w, h})
	t.Dirty(true)
}
