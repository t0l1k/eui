package app

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/sudoku/game"
)

type DialogSelect struct {
	eui.DrawableBase
	title     *eui.Text
	btnClose  *eui.Button
	cSize     *eui.ComboBox
	btnsDiff  []*DiffButton
	modes     *eui.SubjectBase
	show      bool
	gamesData *game.GamesData
}

func NewDialogSelect(gamesData *game.GamesData, f func(b *eui.Button)) *DialogSelect {
	d := &DialogSelect{}
	d.gamesData = gamesData
	d.title = eui.NewText("Выбрать размер поля и сложность")
	d.Add(d.title)
	d.btnClose = eui.NewButton("X", func(b *eui.Button) { d.Visible(false); eui.GetUi().Pop() })
	d.Add(d.btnClose)
	var data []interface{}
	for dim := range *gamesData {
		data = append(data, dim)
	}
	idx := 0
	d.modes = eui.NewSubject()
	d.modes.Attach(d)
	d.modes.SetValue(data[idx])
	d.cSize = eui.NewComboBox(fmt.Sprintf("Размер поля %v", data[idx].(game.Dim)), data, idx, func(cb *eui.ComboBox) {
		d.modes.SetValue(cb.Value())
		str := fmt.Sprintf("Размер поля %v", d.modes)
		d.cSize.SetText(str)
	})
	d.Add(d.cSize)
	for i := 0; i < 4; i++ {
		dim := d.modes.Value().(game.Dim)
		btn := NewDiffButton(dim, game.NewDiff(game.Difficult(i)), f)
		d.modes.Attach(btn)
		d.btnsDiff = append(d.btnsDiff, btn)
		d.Add(btn)
	}
	return d
}

func (d *DialogSelect) UpdateData(value interface{}) {
	switch v := value.(type) {
	case game.Dim:
		for dim, diffs := range *d.gamesData {
			for diff := range diffs {
				if v.Eq(dim) {
					for _, v := range d.btnsDiff {
						if v.diff.Eq(diff) {
							res := d.gamesData.GetLastBest(dim, diff)
							v.SetScore(res)
						}
					}
				}
			}
		}
	}
}

func (d *DialogSelect) Update(dt int) {
	if !d.IsVisible() {
		return
	}
	d.DrawableBase.Update(dt)
}

func (d *DialogSelect) Draw(surface *ebiten.Image) {
	if !d.IsVisible() {
		return
	}
	d.DrawableBase.Draw(surface)
}

func (d *DialogSelect) IsVisible() bool    { return d.show }
func (d *DialogSelect) Visible(value bool) { d.show = value }

func (d *DialogSelect) Resize(rect []int) {
	d.Rect(eui.NewRect(rect))
	w0, h0 := d.GetRect().Size()
	h1 := float64(h0) / 6
	hTop := h1 * 0.8
	wTop := float64(w0) * 0.8
	x, y := d.GetRect().Pos()
	d.title.Resize([]int{x, y, int(wTop), int(hTop)})
	d.btnClose.Resize([]int{x + w0 - int(hTop), y, int(hTop), int(hTop)})
	y += int(h1)
	d.cSize.Resize([]int{x, y, w0, int(h1)})
	y += int(h1)
	for _, v := range d.btnsDiff {
		v.Resize([]int{x, y, w0, int(h1)})
		y += int(h1)
	}
	d.ImageReset()
}
