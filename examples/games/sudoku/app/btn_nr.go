package app

import (
	"image/color"
	"strconv"

	"github.com/t0l1k/eui"
	"golang.org/x/image/colornames"
)

type BottomBarNr struct {
	*eui.Container
	valueBtn *eui.Button
	countLbl *eui.Text
}

func NewBtn(fn func(btn *eui.Button)) *BottomBarNr {
	b := &BottomBarNr{Container: eui.NewContainer(eui.NewAbsoluteLayout())}
	b.valueBtn = eui.NewButton("", fn)
	b.Add(b.valueBtn)
	b.countLbl = eui.NewText("")
	b.countLbl.Fg(colornames.Black)
	b.Add(b.countLbl)
	return b
}

func (d *BottomBarNr) SetHidden(value bool) {
	d.Drawable.SetHidden(value)
	d.Traverse(func(c eui.Drawabler) { c.SetHidden(value) }, false)
}

func (b *BottomBarNr) GetBg() color.Color   { return b.valueBtn.GetBg() }
func (b *BottomBarNr) Bg(value color.Color) { b.valueBtn.Bg(value) }
func (b *BottomBarNr) GetText() string      { return b.valueBtn.GetText() }
func (b *BottomBarNr) GetValue() string     { return b.valueBtn.GetText() }
func (b *BottomBarNr) SetValue(value int)   { b.valueBtn.SetText(strconv.Itoa(value)) }
func (b *BottomBarNr) GetCount() string     { return b.countLbl.GetText() }
func (b *BottomBarNr) SetCount(value int) {
	nr := strconv.Itoa(value)
	b.countLbl.SetText(nr)
	if value == 0 {
		b.countLbl.Bg(colornames.Silver)
	} else {
		b.countLbl.Bg(colornames.Yellow)
	}
}

func (b *BottomBarNr) SetRect(rect eui.Rect[int]) {
	b.Container.SetRect(rect)
	b.valueBtn.SetRect(rect)
	x, y := b.Rect().Pos()
	w, h := b.Rect().Size()
	x += w - w/3
	b.countLbl.SetRect(eui.NewRect([]int{x, y, w / 3, h / 3}))
}
