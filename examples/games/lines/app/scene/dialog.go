package scene

import (
	"fmt"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/lines/app"
)

type Dialog struct {
	*eui.Container
	btnHide, btnNew  *eui.Button
	title, message   *eui.Text
	comboSelGameDiff *eui.SpinBox[int]
	diff             int
	dialogFunc       func(b *eui.Button)
}

func NewDialog(title string, f func(btn *eui.Button)) *Dialog {
	d := &Dialog{Container: eui.NewContainer(eui.NewAbsoluteLayout())}
	d.dialogFunc = f
	d.title = eui.NewText(title)
	d.Add(d.title)
	d.btnHide = eui.NewButton("X", func(b *eui.Button) {
		d.Hide()
	})
	d.Add(d.btnHide)
	d.btnNew = eui.NewButton(app.BNew, f)
	d.Add(d.btnNew)
	d.message = eui.NewText("")
	d.Add(d.message)
	data := []int{app.FieldSizeSmall, app.FieldSizeNormal, app.FieldSizeBig}
	d.diff = data[1]
	d.comboSelGameDiff = eui.NewSpinBox(
		"Выбор размер поля",
		data,
		0,
		func(i int) string { return fmt.Sprintf("%v", i) },
		func(a, b int) int { return a - b },
		func(a, b int) bool { return a == b },
		true,
		2,
	)
	d.comboSelGameDiff.SelectedValue.Connect(func(value int) {
		d.diff = value
	})
	d.Add(d.comboSelGameDiff)
	return d
}

func (d *Dialog) SetTitle(title string) {
	d.title.SetText(title)
}

func (d *Dialog) SetRect(rect eui.Rect[int]) {
	d.Container.SetRect(rect)
	x, y := d.Rect().Pos()
	w, h := d.Rect().W, d.Rect().H/4
	d.title.SetRect(eui.NewRect([]int{x, y, w - h, h}))
	d.btnHide.SetRect(eui.NewRect([]int{x + w - h, y, h, h}))
	y += h
	d.message.SetRect(eui.NewRect([]int{x, y, w, h}))
	y += h
	d.comboSelGameDiff.SetRect(eui.NewRect([]int{x, y, w, h}))
	y += h
	d.btnNew.SetRect(eui.NewRect([]int{x, y, w, h}))
}
