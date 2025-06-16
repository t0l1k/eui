package app

import (
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
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
	eui.DrawableBase
	layoutActs                                      *eui.BoxLayout
	layoutNums                                      *eui.GridLayoutRightDown
	actBtns                                         []eui.Drawabler
	show                                            bool
	fn                                              func(*eui.Button)
	actAccept, actUndo, actDel, actNotes, actNumber bool
	varSw, varDiff                                  *eui.Signal[string]
	board                                           *Board
}

func NewBottomBar(fn func(*eui.Button)) *BottomBar {
	b := &BottomBar{}
	b.fn = fn
	b.layoutActs = eui.NewVLayout()
	b.layoutNums = eui.NewGridLayoutRightDown(2, 2)
	return b
}

func (b *BottomBar) Setup(board *Board) {
	b.board = board
	b.layoutActs.ResetContainerBase()
	b.layoutNums.ResetContainerBase()
	b.layoutNums.SetDim(float64(board.dim.W), float64(board.dim.H))
	b.layoutActs.Add(eui.NewText(title))
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
	b.Resize(b.GetRect().GetArr()) // обязательно после обнуления контейнеров
	b.setBtnClrs()
}

func (b *BottomBar) IsVisible() bool      { return b.show }
func (b *BottomBar) Visible(value bool)   { b.show = value }
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
		for _, v := range b.layoutNums.GetContainer() {
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
		for _, btn := range b.layoutNums.GetContainer() {
			if btn.(*BottomBarNr).GetValue() == strconv.Itoa(k) {
				btn.(*BottomBarNr).SetCount(size - v)
			}
		}
	}
}

func (b *BottomBar) UpdateUndoBtn(count int) {
	for _, btn := range b.layoutActs.GetContainer() {
		switch btn := btn.(type) {
		case *eui.Button:
			if strings.HasPrefix(btn.GetText(), aUndo) {
				btn.SetText(aUndo + ":" + strconv.Itoa(count))
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
	b.varSw.Emit(b.board.sw.StringShort())
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
	h1 := h0 / 2
	b.layoutActs.Resize([]int{x, y, w0, h1})
	y += h1
	b.layoutNums.Resize([]int{x, y, w0, h1})
	b.ImageReset()
}
