package app

import (
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
	board           Sols
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

func (b *BottomBar) Setup(board Sols) {
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

func (s *BottomBar) UpdateMoveCount() bool {
	moves, str := s.board.AvailableMoves()
	s.varMoves.SetValue(str)
	return moves == 0
}

func (b *BottomBar) Update(dt int) {
	if !b.IsVisible() {
		return
	}
	for _, v := range b.layout.GetContainer() {
		v.Update(dt)
	}
	b.varSw.SetValue(b.board.Stopwatch().StringShort())
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
