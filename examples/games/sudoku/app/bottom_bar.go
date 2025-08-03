package app

import (
	"strconv"
	"strings"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/sudoku/game"
	"golang.org/x/image/colornames"
)

const (
	aAccept = "Применить игру"
	aUndo   = "Отменить ход"
	aDel    = "Удалить"
	aNote   = "Заметка"
)

var actsStr = []string{aUndo, aDel, aNote}

type BottomBar struct {
	*eui.Container
	layoutActs, layoutNums                          *eui.Container
	actBtns                                         []eui.Drawabler
	fn                                              func(*eui.Button)
	actAccept, actUndo, actDel, actNotes, actNumber bool
	varSw, varDiff                                  *eui.Signal[string]
	board                                           *Board
}

func NewBottomBar(fn func(*eui.Button)) *BottomBar {
	b := &BottomBar{Container: eui.NewContainer(eui.NewAbsoluteLayout())}
	b.fn = fn
	b.layoutActs = eui.NewContainer(eui.NewVBoxLayout(1))
	b.Add(b.layoutActs)
	b.layoutNums = eui.NewContainer(eui.NewGridLayout(2, 2, 1))
	b.Add(b.layoutNums)
	return b
}

func (b *BottomBar) Setup(board *Board) {
	b.board = board
	b.layoutActs.ResetContainer()
	b.layoutNums.ResetContainer()
	b.layoutNums.SetLayout(eui.NewGridLayout(float64(board.dim.W), float64(board.dim.H), 1))
	b.layoutActs.Add(eui.NewText(Title))
	b.layoutActs.Add(eui.NewText(board.dim.String()))
	diffText := eui.NewText("")
	b.layoutActs.Add(diffText)
	b.varDiff = eui.NewSignal(func(a, b string) bool { return a == b })
	b.varDiff.ConnectAndFire(func(data string) { diffText.SetText(data) }, b.board.GetDiffStr())
	swStr := eui.NewText("")
	b.varSw = eui.NewSignal(func(a, b string) bool { return a == b })
	b.varSw.Connect(func(data string) { swStr.SetText(data) })
	b.varSw.Emit(b.board.sw.StringShort())
	b.layoutActs.Add(swStr)
	b.actBtns = nil
	if b.board.diff.String() == game.Manual.String() {
		btn := eui.NewButton(aAccept, b.fn)
		b.layoutActs.Add(btn)
		b.actBtns = append(b.actBtns, btn)
	}
	for _, v := range actsStr {
		btn := eui.NewButton(v, b.fn)
		b.layoutActs.Add(btn)
		b.actBtns = append(b.actBtns, btn)
	}
	for i := 0; i < b.board.dim.Size(); i++ {
		btn := NewBtn(b.fn)
		btn.SetValue(i + 1)
		btn.SetCount(0)
		b.layoutNums.Add(btn)
		b.actBtns = append(b.actBtns, btn)
	}
	b.SetRect(b.Rect()) // обязательно после обнуления контейнеров
	b.setBtnClrs()
}

func (d *BottomBar) SetHidden(value bool) {
	d.Drawable.SetHidden(value)
	d.Traverse(func(c eui.Drawabler) { c.SetHidden(value); c.MarkDirty() }, false)
}
func (b *BottomBar) IsActAccept() bool    { return b.actAccept }
func (b *BottomBar) IsActUndo() bool      { return b.actUndo }
func (b *BottomBar) IsActDel() bool       { return b.actDel }
func (b *BottomBar) IsActNotes() bool     { return b.actNotes }
func (b *BottomBar) ShowNotes(value bool) { b.actNotes = value }

func (b *BottomBar) SetAct(btn *eui.Button) (result bool) {
	b.setBtnClrs()
	btnStr := (btn.GetText())
	if strings.HasPrefix(btnStr, aUndo) {
		btnStr = aUndo
	}
	switch btnStr {
	case aAccept:
		b.actAccept = true
		btn.Bg(colornames.Yellow)
		result = false
	case aUndo:
		b.actUndo = true
		btn.Bg(colornames.Yellow)
		result = false
	case aDel:
		b.actDel = true
		btn.Bg(colornames.Yellow)
		result = false
	case aNote:
		b.actNotes = !b.actNotes
		if b.actNotes {
			btn.Bg(colornames.Yellow)
		}
		result = false
	default:
		for _, v := range b.layoutNums.Childrens() {
			if v.(*BottomBarNr).GetText() == btn.GetText() {
				v.(*BottomBarNr).Bg(colornames.Yellow)
			}
		}
		b.actNumber = true
		result = true
	}
	return result
}

func (b *BottomBar) setBtnClrs() {
	b.actAccept = false
	b.actUndo = false
	b.actDel = false
	b.actNumber = false
	for _, v := range b.actBtns {
		switch vv := v.(type) {
		case *eui.Button:
			if vv.GetText() == aNote && b.actNotes {
				vv.Bg(colornames.Yellow)
			} else if vv.GetBg() == colornames.Yellow {
				vv.Bg(colornames.Silver)
			}
		case *BottomBarNr:
			if vv.GetBg() == colornames.Yellow {
				vv.Bg(colornames.Silver)
			}
		}
	}
}

func (b *BottomBar) UpdateNrs(counts map[int]int) {
	b.varDiff.Emit(b.board.GetDiffStr())
	size := b.board.dim.Size()
	for k, v := range counts {
		for _, btn := range b.layoutNums.Childrens() {
			if btn.(*BottomBarNr).GetValue() == strconv.Itoa(k) {
				btn.(*BottomBarNr).SetCount(size - v)
			}
		}
	}
}

func (b *BottomBar) UpdateUndoBtn(count int) {
	for _, btn := range b.layoutActs.Childrens() {
		switch btn := btn.(type) {
		case *eui.Button:
			if strings.HasPrefix(btn.GetText(), aUndo) {
				btn.SetText(aUndo + ":" + strconv.Itoa(count))
			}
		}
	}
}

func (b *BottomBar) Update(dt int) {
	b.Container.Update(dt)
	if b.board == nil {
		return
	}
	if b.board.sw.IsRun() {
		b.varSw.Emit(b.board.sw.StringShort())
	}
}

func (b *BottomBar) SetRect(rect eui.Rect[int]) {
	b.Container.SetRect(rect)
	w0, h0 := b.Rect().Size()
	x, y := b.Rect().Pos()
	h1 := h0 / 2
	b.layoutActs.SetRect(eui.NewRect([]int{x, y, w0, h1}))
	y += h1
	b.layoutNums.SetRect(eui.NewRect([]int{x, y, w0, h1}))
	b.ImageReset()
}
