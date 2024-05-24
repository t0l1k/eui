package app

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/sudoku/game"
)

type DiffButton struct {
	eui.DrawableBase
	title, lblScore *eui.Text
	btn             *eui.Button
	dim             game.Dim
	diff            game.Difficult
	f               func(b *eui.Button)
}

func NewDiffButton(dim game.Dim, diff game.Difficult, f func(b *eui.Button)) *DiffButton {
	d := &DiffButton{}
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

func (d *DiffButton) Visible(value bool) {
	for _, v := range d.GetContainer() {
		switch vT := v.(type) {
		case *eui.Text:
			vT.Visible(value)
		case *eui.Button:
			vT.Visible(value)
		}
	}
}

func (d *DiffButton) Resize(rect []int) {
	d.Rect(eui.NewRect(rect))
	w0, h0 := d.GetRect().Size()
	h1 := float64(h0) / 2
	hTop := h1 * 0.8
	wTop := float64(w0) * 0.8
	x, y := d.GetRect().Pos()
	d.title.Resize([]int{x, y, int(wTop), int(hTop)})
	y += int(h1)
	w2 := w0 / 3
	d.lblScore.Resize([]int{x, y, w2, int(h1)})
	d.btn.Resize([]int{x + w2, y, w2 * 2, int(h1)})
	d.ImageReset()
}
