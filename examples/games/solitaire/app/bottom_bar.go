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
	*eui.Container
	layout          *eui.Container
	fn              func(*eui.Button)
	board           Sols
	varSw, varMoves *eui.Signal[string]
}

func NewBottomBar(fn func(*eui.Button)) *BottomBar {
	b := &BottomBar{Container: eui.NewContainer(eui.NewAbsoluteLayout())}
	b.layout = eui.NewContainer(eui.NewHBoxLayout(1))
	b.Add(b.layout)
	b.fn = fn
	b.varSw = eui.NewSignal(func(a, b string) bool { return a == b })
	b.varMoves = eui.NewSignal(func(a, b string) bool { return a == b })
	// b.Hide()
	return b
}

func (b *BottomBar) Setup(board Sols) {
	b.board = board
	b.layout.ResetContainer()
	for _, str := range actStrs {
		b.layout.Add(eui.NewButton(str, b.fn))
	}
	movesText := eui.NewLabel("Ходов:0")
	b.varMoves.Connect(func(data string) { movesText.SetText(data) })
	b.layout.Add(movesText)
	swText := eui.NewLabel("")
	b.varSw.Connect(func(data string) { swText.SetText(data) })
	b.layout.Add(swText)
}

func (s *BottomBar) UpdateMoveCount() bool {
	moves, str := s.board.AvailableMoves()
	s.varMoves.Emit(str)
	return moves == 0
}

func (b *BottomBar) Update(dt int) {
	if b.IsHidden() {
		return
	}
	for _, v := range b.layout.Children() {
		v.Update(dt)
	}
	b.varSw.Emit(b.board.Stopwatch().StringShort())
}

func (b *BottomBar) Draw(surface *ebiten.Image) {
	if b.IsHidden() {
		return
	}
	for _, v := range b.layout.Children() {
		v.Draw(surface)
	}
}

func (b *BottomBar) SetRect(rect eui.Rect[int]) {
	b.Container.SetRect(rect)
	b.layout.SetRect(rect)
}
