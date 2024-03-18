package app

import (
	"fmt"

	"github.com/t0l1k/eui"
)

type DialogSelect struct {
	eui.DrawableBase
	title    *eui.Text
	btnClose *eui.Button
	cSize    *eui.ComboBox
	btnsDiff []*DiffButton
	size     *eui.SubjectBase
}

func NewDialogSelect(f func(b *eui.Button)) *DialogSelect {
	d := &DialogSelect{}
	d.title = eui.NewText("Выбрать размер поля и сложность")
	d.Add(d.title)
	d.btnClose = eui.NewButton("X", func(b *eui.Button) { d.Visible(false) })
	d.Add(d.btnClose)
	data := []interface{}{2, 3, 4}
	idx := 0
	d.size = eui.NewSubject()
	d.size.SetValue(data[idx])
	d.cSize = eui.NewComboBox("Размер поля", data, idx, func(cb *eui.ComboBox) {
		d.size.SetValue(cb.Value())
		str := fmt.Sprintf("Выбран размер поля %vx%v", d.size, d.size)
		d.cSize.SetText(str)
	})
	d.Add(d.cSize)
	for i := 0; i < 3; i++ {
		btn := NewDiffButton(d.size.Value().(int), Difficult(i), "3:00", f)
		d.size.Attach(btn)
		d.btnsDiff = append(d.btnsDiff, btn)
		d.Add(btn)
	}
	return d
}

func (d *DialogSelect) Visible(value bool) {
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
		case *DiffButton:
			vT.Visible(value)
			if value {
				vT.Enable()
			} else {
				vT.Disable()
			}
		}
	}
}

func (d *DialogSelect) Resize(rect []int) {
	d.Rect(eui.NewRect(rect))
	w0, h0 := d.GetRect().Size()
	h1 := float64(h0) / 5
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
