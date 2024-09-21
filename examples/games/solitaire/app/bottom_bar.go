package app

import (
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
)

const (
	actNextSol = " Следующий пасьянс"
	actNew     = "Новая игра"
	actReset   = "Играть сначала"
	// actAvailibleMove = "Доступные ходы"
	actBackwardMove = "<"
	actForwardMove  = ">"
)

var actStrs = []string{actNextSol, actNew, actReset, actBackwardMove, actForwardMove}

type BottomBar struct {
	eui.DrawableBase
	layout          *eui.BoxLayout
	fn              func(*eui.Button)
	board           *BoardSol15
	varSw, varMoves *eui.SubjectBase
}

func NewBottomBar(fn func(*eui.Button)) *BottomBar {
	b := &BottomBar{}
	b.layout = eui.NewHLayout()
	b.fn = fn
	b.varSw = eui.NewSubject()
	b.varMoves = eui.NewSubject()
	b.Visible(true)
	return b
}

func (b *BottomBar) Setup(board *BoardSol15) {
	b.board = board
	b.layout.ResetContainerBase()
	for _, str := range actStrs {
		b.layout.Add(eui.NewButton(str, b.fn))
	}
	movesText := eui.NewText("Ходов:0")
	b.varMoves.Attach(movesText)
	b.layout.Add(movesText)
	swText := eui.NewText("")
	b.varSw.Attach(swText)
	b.layout.Add(swText)
}

func (s *BottomBar) updateMoveCount() {
	str := "Ход:" + strconv.Itoa(len(s.board.historyOfMoves)) + " доступно:" + strconv.Itoa(s.board.game.AvailableMoves()) + " ходов"
	s.varMoves.SetValue(str)
}

func (b *BottomBar) Update(dt int) {
	if !b.IsVisible() {
		return
	}
	for _, v := range b.layout.GetContainer() {
		v.Update(dt)
	}
	b.varSw.SetValue(b.board.sw.StringShort())
}

func (b *BottomBar) Draw(surface *ebiten.Image) {
	if !b.IsVisible() {
		return
	}
	for _, v := range b.layout.GetContainer() {
		v.Draw(surface)
	}
}

func (b *BottomBar) Resize(rect []int) {
	b.SpriteBase.Rect(eui.NewRect(rect))
	b.Rect(eui.NewRect(rect))
	b.layout.Resize(rect)
}
