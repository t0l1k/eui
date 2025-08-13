package app

import (
	"fmt"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/sudoku/game"
)

type DialogSelect struct {
	*eui.Container
	title      *eui.Text
	btnClose   *eui.Button
	lblHistory *eui.Text
	history    *eui.ListView
	cSize      *eui.SpinBox[game.Dim]
	btnsDiff   []*DiffButton
	modes      *eui.Signal[game.Dim]
	gamesData  *game.GamesData
	margin     int
}

func NewDialogSelect(gamesData *game.GamesData, fn func(b *eui.Button)) *DialogSelect {
	d := &DialogSelect{Container: eui.NewContainer(eui.NewAbsoluteLayout())}
	d.gamesData = gamesData
	d.title = eui.NewText("Выбрать размер поля и сложность")
	d.Add(d.title)
	d.btnClose = eui.NewButton("X", func(b *eui.Button) { d.Hide(); eui.GetUi().Pop() })
	d.Add(d.btnClose)
	var data []game.Dim
	for _, dim := range gamesData.SortedDims() {
		data = append(data, dim)
	}
	idx := 0
	d.modes = eui.NewSignal[game.Dim]()
	d.modes.Connect(func(value game.Dim) {
		for dim, diffs := range *d.gamesData {
			for diff := range diffs {
				if value.Eq(dim) {
					for _, v := range d.btnsDiff {
						if v.diff.Eq(diff) {
							res := d.gamesData.GetLastBest(dim, diff)
							v.SetScore(res)
						}
					}
				}
			}
		}
	})
	d.modes.Emit(data[idx])
	d.cSize = eui.NewSpinBox(
		fmt.Sprintf("Размер поля %v", data[idx]),
		data,
		0,
		func(d game.Dim) string { return d.String() },
		func(a, b game.Dim) int { return a.W*a.H - b.W*b.H },
		func(a, b game.Dim) bool { return a.Eq(b) },
		false,
		0,
	)
	d.cSize.SelectedValue.Connect(func(data game.Dim) {
		d.modes.Emit(data)
	})

	d.Add(d.cSize)
	for i := 0; i < 5; i++ {
		dim := d.modes.Value()
		btn := NewDiffButton(dim, game.NewDiff(game.Difficult(i)), fn)
		d.modes.Connect(func(data game.Dim) { btn.dim = data })
		d.btnsDiff = append(d.btnsDiff, btn)
		d.Add(btn)
	}
	d.lblHistory = eui.NewText("История игр")
	d.Add(d.lblHistory)
	d.history = eui.NewListView()
	d.Add(d.history)
	return d
}

func (d *DialogSelect) SetRect(rect eui.Rect[int]) {
	d.Container.SetRect(rect)
	w0, h0 := d.Rect().Size()
	h1 := float64(h0) / 7
	d.margin = int(h1 / 2)
	hTop := h1 * 0.8
	wTop := float64(w0) * 0.7
	x0, y0 := d.Rect().Pos()
	d.title.SetRect(eui.NewRect([]int{x0, y0, w0 - int(hTop), int(hTop)}))
	x := x0 + (w0 - int(hTop))
	y := y0
	d.btnClose.SetRect(eui.NewRect([]int{x, y, int(hTop), int(hTop)}))
	x = x0
	y += int(h1)
	d.cSize.SetRect(eui.NewRect([]int{x, y, int(wTop), int(h1)}))
	x = x0 + w0 - (w0 - int(wTop))
	d.lblHistory.SetRect(eui.NewRect([]int{x, y, w0 - int(wTop), int(hTop)}))
	y += int(h1)
	d.history.SetRect(eui.NewRect([]int{x, y, (w0 - int(wTop)), h0 - int(h1*2)}))
	x = x0
	for _, v := range d.btnsDiff {
		v.SetRect(eui.NewRect([]int{x, y, int(wTop), int(h1)}))
		y += int(h1)
	}
	d.ImageReset()
}
