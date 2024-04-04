package app

import (
	"image/color"

	"github.com/t0l1k/eui"
)

type BottomBarNr struct {
	eui.DrawableBase
	valueBtn *eui.Button
	countLbl *eui.Text
}

func NewBtn(fn func(btn *eui.Button)) *BottomBarNr {
	b := &BottomBarNr{}
	b.valueBtn = eui.NewButton("", fn)
	b.Add(b.valueBtn)
	b.countLbl = eui.NewText("")
	b.Add(b.countLbl)
	return b
}

func (b *BottomBarNr) GetBg() color.Color    { return b.valueBtn.GetBg() }
func (b *BottomBarNr) Bg(value color.Color)  { b.valueBtn.Bg(value) }
func (b *BottomBarNr) GetText() string       { return b.valueBtn.GetText() }
func (b *BottomBarNr) GetValue() string      { return b.valueBtn.GetText() }
func (b *BottomBarNr) SetValue(value string) { b.valueBtn.SetText(value) }
func (b *BottomBarNr) GetCount() string      { return b.countLbl.GetText() }
func (b *BottomBarNr) SetCount(value string) { b.countLbl.SetText(value) }

func (b *BottomBarNr) Resize(rect []int) {
	b.Rect(eui.NewRect(rect))
	b.valueBtn.Resize(rect)
	x, y := b.GetRect().Pos()
	w, h := b.GetRect().Size()
	x += w - w/3
	b.countLbl.Resize([]int{x, y, w / 3, h / 3})
}
