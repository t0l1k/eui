package app

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/sudoku/game"
)

type DialogSelect struct {
	*eui.Container
	title      *eui.Label
	btnClose   *eui.Button
	lblHistory *eui.Label
	history    *eui.ListView
	cSize      *eui.SpinBox[game.Dim]
	btnsDiff   []*DiffButton
	modes      *eui.Signal[game.Dim]
	gamesData  *game.GamesData
	margin     int
}

func NewDialogSelect(gamesData *game.GamesData, fn func(b *eui.Button)) *DialogSelect {
	d := &DialogSelect{Container: eui.NewContainer(eui.NewLayoutVerticalPercent([]int{10, 90}, 1))}
	d.gamesData = gamesData
	d.title = eui.NewLabel("Выбрать размер поля и сложность")
	d.btnClose = eui.NewButton("X", func(b *eui.Button) { d.Hide(); eui.GetUi().Pop() })
	data := gamesData.SortedDims()
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
		"Размер поля",
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

	contTitle := eui.NewContainer(eui.NewLayoutHorizontalPercent([]int{5, 95}, 1))
	contTitle.Add(d.btnClose)
	contTitle.Add(d.title)
	d.Add(contTitle)

	contDialog := eui.NewContainer(eui.NewLayoutHorizontalPercent([]int{70, 30}, 1))

	contBtns := eui.NewContainer(eui.NewVBoxLayout(5))
	contBtns.Add(d.cSize)
	for i := 0; i < 5; i++ {
		dim := d.modes.Value()
		btn := NewDiffButton(dim, game.NewDiff(game.Difficult(i)), fn)
		d.modes.Connect(func(data game.Dim) { btn.dim = data })
		d.btnsDiff = append(d.btnsDiff, btn)
		contBtns.Add(btn)
	}

	contDialog.Add(contBtns)

	contHist := eui.NewContainer(eui.NewLayoutVerticalPercent([]int{10, 90}, 1))
	d.lblHistory = eui.NewLabel("История игр")
	d.history = eui.NewListView()
	contHist.Add(d.lblHistory)
	contHist.Add(d.history)

	contDialog.Add(contHist)
	d.Add(contDialog)
	return d
}
