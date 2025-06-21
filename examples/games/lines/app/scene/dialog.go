package scene

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/lines/app"
)

type Dialog struct {
	*eui.Container
	btnHide, btnNew  *eui.Button
	title, message   *eui.Text
	comboSelGameDiff *eui.ComboBox
	diff             int
	dialogFunc       func(b *eui.Button)
}

func NewDialog(title string, f func(btn *eui.Button)) *Dialog {
	d := &Dialog{Container: eui.NewContainer(eui.NewAbsoluteLayout())}
	d.dialogFunc = f
	d.title = eui.NewText(title)
	d.Add(d.title)
	d.btnHide = eui.NewButton("X", func(b *eui.Button) {
		d.Visible(false)
	})
	d.Add(d.btnHide)
	d.btnNew = eui.NewButton(app.BNew, f)
	d.Add(d.btnNew)
	d.message = eui.NewText("")
	d.Add(d.message)
	data := []interface{}{app.FieldSizeSmall, app.FieldSizeNormal, app.FieldSizeBig}
	d.diff = data[0].(int)
	d.comboSelGameDiff = eui.NewComboBox("Выбор размер поля", data, 0, func(cb *eui.ComboBox) {
		d.diff = cb.Value().(int)
	})
	d.Add(d.comboSelGameDiff)
	return d
}

func (d *Dialog) Visible(value bool) {
	d.Traverse(func(child eui.Drawabler) {
		child.Visible(value)
		if value {
			child.Enable()
		} else {
			child.Disable()
		}
	}, false)
}

func (d *Dialog) SetTitle(title string) {
	d.title.SetText(title)
}

func (d *Dialog) Resize(rect eui.Rect) {
	d.SetRect(rect)
	x, y := d.Rect().Pos()
	w, h := d.Rect().W, d.Rect().H/4
	d.title.Resize(eui.NewRect([]int{x, y, w - h, h}))
	d.btnHide.Resize(eui.NewRect([]int{x + w - h, y, h, h}))
	y += h
	d.message.Resize(eui.NewRect([]int{x, y, w, h}))
	y += h
	d.comboSelGameDiff.Resize(eui.NewRect([]int{x, y, w, h}))
	y += h
	d.btnNew.Resize(eui.NewRect([]int{x, y, w, h}))
}
