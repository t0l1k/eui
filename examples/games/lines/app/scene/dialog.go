package scene

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/lines/app"
)

type Dialog struct {
	eui.DrawableBase
	btnHide, btnNew  *eui.Button
	title, message   *eui.Text
	comboSelGameDiff *eui.ComboBox
	diff             int
	dialogFunc       func(b *eui.Button)
}

func NewDialog(title string, f func(btn *eui.Button)) *Dialog {
	d := &Dialog{}
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
	for _, v := range d.GetContainer() {
		switch vT := v.(type) {
		case *eui.Text:
			vT.Visible(value)
			if value {
				vT.Enable()
			} else {
				vT.Disable()
			}
		case *eui.Button:
			vT.Visible(value)
			if value {
				vT.Enable()
			} else {
				vT.Disable()
			}
		case *eui.ComboBox:
			vT.Visible(value)
			if value {
				vT.Enable()
			} else {
				vT.Disable()
			}
		}
	}
}

func (d *Dialog) SetTitle(title string) {
	d.title.SetText(title)
}

func (d *Dialog) Resize(rect []int) {
	d.Rect(eui.NewRect(rect))
	d.SpriteBase.Resize(rect)
	x, y := d.GetRect().Pos()
	w, h := d.GetRect().W, d.GetRect().H/4
	d.title.Resize([]int{x, y, w - h, h})
	d.btnHide.Resize([]int{x + w - h, y, h, h})
	y += h
	d.message.Resize([]int{x, y, w, h})
	y += h
	d.comboSelGameDiff.Resize([]int{x, y, w, h})
	y += h
	d.btnNew.Resize([]int{x, y, w, h})
}
