package app

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/sudoku/game"
)

type DiffButton struct {
	*eui.Container
	title, lblScore *eui.Label
	btn             *eui.Button
	dim             game.Dim
	diff            game.Difficult
	fn              func(b *eui.Button)
}

func NewDiffButton(dim game.Dim, diff game.Difficult, fn func(b *eui.Button)) *DiffButton {
	d := &DiffButton{Container: eui.NewContainer(eui.NewAbsoluteLayout())}
	d.dim = dim
	d.diff = diff
	d.fn = fn
	d.title = eui.NewLabel(diff.String())
	d.Add(d.title)
	d.lblScore = eui.NewLabel("")
	d.Add(d.lblScore)
	d.btn = eui.NewButton("Запустить", fn)
	d.Add(d.btn)
	return d
}

func (d *DiffButton) GetData() (game.Dim, game.Difficult) {
	return d.dim, d.diff
}

func (d *DiffButton) SetScore(value string) {
	d.lblScore.SetText(value)
}

func (d *DiffButton) UpdateData(value interface{}) {
	switch v := value.(type) {
	case game.Dim:
		d.dim = v
	}
}

func (d *DiffButton) SetRect(rect eui.Rect[int]) {
	d.Container.SetRect(rect)
	w0, h0 := d.Rect().Size()
	h1 := float64(h0) / 2
	hTop := h1 * 0.8
	wTop := float64(w0) * 0.8
	x, y := d.Rect().Pos()
	d.title.SetRect(eui.NewRect([]int{x, y, int(wTop), int(hTop)}))
	y += int(h1)
	w2 := w0 / 3
	d.lblScore.SetRect(eui.NewRect([]int{x, y, w2, int(h1)}))
	d.btn.SetRect(eui.NewRect([]int{x + w2, y, w2 * 2, int(h1)}))
	d.ImageReset()
}
