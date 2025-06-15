package main

import (
	"github.com/t0l1k/eui"
)

type SceneTestPlot struct {
	eui.SceneBase
	plot *eui.Plot
}

func NewSceneTestPlot() *SceneTestPlot {
	s := &SceneTestPlot{}
	xArr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	yArr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	values := []int{1, 2, 3, 4, 5, 5, 6, 7, 7, 8}
	s.plot = eui.NewPlot(xArr, yArr, values, "Game Score", "Game", "Level")
	s.Add(s.plot)
	s.Resize()
	return s
}

func (s *SceneTestPlot) Resize() {
	w0, h0 := eui.GetUi().Size()
	rect := eui.NewRect([]int{0, 0, w0, h0})
	margin := int(float64(rect.GetLowestSize()) * 0.1)
	x, y := margin, margin
	w, h := w0-margin*2, h0-margin*2
	s.plot.Resize([]int{x, y, w, h})
}

func NewGame() *eui.Ui {
	u := eui.GetUi()
	u.SetTitle("Test Plot")
	k := 30
	w, h := 9*k, 6*k
	u.SetSize(w, h)
	return u
}

func main() {
	eui.Init(NewGame())
	eui.Run(NewSceneTestPlot())
	eui.Quit(func() {})
}
