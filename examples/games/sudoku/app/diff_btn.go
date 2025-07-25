package app

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/sudoku/game"
)

type DiffButton struct {
	*eui.Container
	title, lblScore *eui.Text
	btn             *eui.Button
	dim             game.Dim
	diff            game.Difficult
	f               func(b *eui.Button)
}

func NewDiffButton(dim game.Dim, diff game.Difficult, f func(b *eui.Button)) *DiffButton {
	d := &DiffButton{Container: eui.NewContainer(eui.NewAbsoluteLayout())}
	d.dim = dim
	d.diff = diff
	d.f = f
	d.title = eui.NewText(diff.String())
	d.Add(d.title)
	d.lblScore = eui.NewText("")
	d.Add(d.lblScore)
	d.btn = eui.NewButton("Запустить", f)
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

func (d *DiffButton) SetHidden(value bool) {
	d.Drawable.SetHidden(value)
	d.Traverse(func(c eui.Drawabler) { c.SetHidden(value) }, false)
}

func (d *DiffButton) Resize(rect eui.Rect[int]) {
	d.SetRect(rect)
	w0, h0 := d.Rect().Size()
	h1 := float64(h0) / 2
	hTop := h1 * 0.8
	wTop := float64(w0) * 0.8
	x, y := d.Rect().Pos()
	d.title.Resize(eui.NewRect([]int{x, y, int(wTop), int(hTop)}))
	y += int(h1)
	w2 := w0 / 3
	d.lblScore.Resize(eui.NewRect([]int{x, y, w2, int(h1)}))
	d.btn.Resize(eui.NewRect([]int{x + w2, y, w2 * 2, int(h1)}))
	d.ImageReset()
}
