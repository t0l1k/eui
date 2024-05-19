package app

import (
	"time"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/sudoku/game"
)

type DiffButton struct {
	eui.DrawableBase
	title, lblScore *eui.Text
	btn             *eui.Button
	dim             *game.Dim
	diff            game.Difficult
	score           *game.Score
	f               func(b *eui.Button)
}

func NewDiffButton(dim *game.Dim, diff game.Difficult, score *game.Score, f func(b *eui.Button)) *DiffButton {
	d := &DiffButton{}
	d.dim = dim
	d.score = score
	d.diff = diff
	d.f = f
	d.title = eui.NewText(diff.String())
	d.Add(d.title)
	d.lblScore = eui.NewText(score.String())
	d.Add(d.lblScore)
	d.btn = eui.NewButton("Запустить", f)
	d.Add(d.btn)
	return d
}

func (d *DiffButton) GetData() (*game.Dim, string, game.Difficult) {
	return d.dim, d.score.String(), d.diff
}

func (d *DiffButton) SetScore(value time.Duration) {
	d.score.Last(value)
	d.lblScore.SetText(d.score.String())
}

func (d *DiffButton) UpdateData(value interface{}) {
	switch v := value.(type) {
	case int:
		d.dim.W = v
		d.dim.H = v
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
