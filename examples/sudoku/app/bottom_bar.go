package app

import (
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
)

const (
	aUndo = "Отменить"
	aDel  = "Удалить"
	aNote = "Заметка"
)

var actsStr = []string{aUndo, aDel, aNote}

type BottomBar struct {
	eui.DrawableBase
	layoutActs, layoutNums               *eui.BoxLayout
	actBtns                              []eui.Drawabler
	dim                                  int
	show                                 bool
	fn                                   func(*eui.Button)
	actUndo, actDel, actNotes, actNumber bool
}

func NewBottomBar(fn func(*eui.Button)) *BottomBar {
	b := &BottomBar{}
	b.fn = fn
	b.layoutActs = eui.NewHLayout()
	b.layoutNums = eui.NewHLayout()
	b.Bg(eui.Red)
	b.Fg(eui.Silver)
	return b
}

func (b *BottomBar) Setup(dim int) {
	b.dim = dim
	b.layoutActs.ResetContainerBase()
	b.layoutNums.ResetContainerBase()
	for _, v := range actsStr {
		btn := eui.NewButton(v, b.fn)
		b.layoutActs.Add(btn)
		b.actBtns = append(b.actBtns, btn)
	}
	size := b.dim * b.dim
	for i := 0; i < size; i++ {
		btn := NewBtn(b.fn)
		btn.SetValue(i + 1)
		btn.SetCount(0)
		b.layoutNums.Add(btn)
		b.actBtns = append(b.actBtns, btn)
	}
	b.Resize(b.GetRect().GetArr()) // обязательно после обнуления контейнеров
	b.setBtnClrs()
}

func (b *BottomBar) IsVisible() bool      { return b.show }
func (b *BottomBar) Visible(value bool)   { b.show = value }
func (b *BottomBar) IsActUndo() bool      { return b.actUndo }
func (b *BottomBar) IsActDel() bool       { return b.actDel }
func (b *BottomBar) IsActNotes() bool     { return b.actNotes }
func (b *BottomBar) ShowNotes(value bool) { b.actNotes = value }

func (b *BottomBar) SetAct(btn *eui.Button) (result bool) {
	b.setBtnClrs()
	switch btn.GetText() {
	case aUndo:
		b.actUndo = true
		btn.Bg(eui.Yellow)
		result = false
	case aDel:
		b.actDel = true
		btn.Bg(eui.Yellow)
		result = false
	case aNote:
		b.actNotes = !b.actNotes
		if b.actNotes {
			btn.Bg(eui.Yellow)
		}
		result = false
	default:
		for _, v := range b.layoutNums.GetContainer() {
			if v.(*BottomBarNr).GetText() == btn.GetText() {
				v.(*BottomBarNr).Bg(eui.Yellow)
			}
		}
		b.actNumber = true
		result = true
	}
	return result
}

func (b *BottomBar) setBtnClrs() {
	b.actUndo = false
	b.actDel = false
	b.actNumber = false
	for _, v := range b.actBtns {
		switch vv := v.(type) {
		case *eui.Button:
			if vv.GetText() == aNote && b.actNotes {
				vv.Bg(eui.Yellow)
			} else if vv.GetBg() == eui.Yellow {
				vv.Bg(eui.Silver)
			}
		case *BottomBarNr:
			if vv.GetBg() == eui.Yellow {
				vv.Bg(eui.Silver)
			}
		}
	}
}

func (b *BottomBar) UpdateNrs(counts map[int]int) {
	size := b.dim * b.dim
	for k, v := range counts {
		for _, btn := range b.layoutNums.GetContainer() {
			if btn.(*BottomBarNr).GetValue() == strconv.Itoa(k) {
				btn.(*BottomBarNr).SetCount(size - v)
			}
		}
	}
}

func (b *BottomBar) Update(dt int) {
	if !b.IsVisible() {
		return
	}
	for _, v := range b.layoutActs.GetContainer() {
		v.Update(dt)
	}
	for _, v := range b.layoutNums.GetContainer() {
		v.Update(dt)
	}
}

func (b *BottomBar) Draw(surface *ebiten.Image) {
	if !b.IsVisible() {
		return
	}
	for _, v := range b.layoutActs.GetContainer() {
		v.Draw(surface)
	}
	for _, v := range b.layoutNums.GetContainer() {
		v.Draw(surface)
	}
}

func (b *BottomBar) Resize(rect []int) {
	b.Rect(eui.NewRect(rect))
	b.SpriteBase.Resize(rect)
	w0, h0 := b.GetRect().Size()
	x, y := b.GetRect().Pos()
	h1 := h0 / 3
	b.layoutActs.Resize([]int{x, y, w0, h1})
	y += h1
	b.layoutNums.Resize([]int{x, y, w0, h1 * 2})
	b.ImageReset()
}
