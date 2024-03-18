package app

import (
	"github.com/t0l1k/eui"
)

type Difficult int

const (
	Easy Difficult = iota
	Normal
	Hard
)

func (d Difficult) Size() int { return int(Hard) }
func (d Difficult) String() (res string) {
	switch d {
	case Easy:
		res = "Легко"
	case Normal:
		res = "Нормально"
	case Hard:
		res = "Сложно"
	}
	return res
}

type DiffButton struct {
	eui.DrawableBase
	title, bestScore *eui.Text
	btn              *eui.Button
	size             int
	diff             Difficult
	score            string
	f                func(b *eui.Button)
}

func NewDiffButton(size int, diff Difficult, score string, f func(b *eui.Button)) *DiffButton {
	d := &DiffButton{}
	d.size = size
	d.score = score
	d.diff = diff
	d.f = f
	d.title = eui.NewText(diff.String())
	d.Add(d.title)
	d.bestScore = eui.NewText(score)
	d.Add(d.bestScore)
	d.btn = eui.NewButton("Запустить", f)
	d.Add(d.btn)
	return d
}

func (d *DiffButton) GetData() (int, string, Difficult) {
	return d.size, d.score, d.diff
}

func (d *DiffButton) UpdateData(value interface{}) {
	switch v := value.(type) {
	case int:
		d.size = v
	}
}

func (d *DiffButton) Visible(value bool) {
	for _, v := range d.GetContainer() {
		switch vT := v.(type) {
		case *eui.Text:
			vT.Visible = value
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
	d.bestScore.Resize([]int{x, y, w2, int(h1)})
	d.btn.Resize([]int{x + w2, y, w2 * 2, int(h1)})
	d.ImageReset()
}
