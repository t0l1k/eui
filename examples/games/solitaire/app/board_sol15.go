package app

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/solitaire/game"
)

type BoardSol15 struct {
	eui.DrawableBase
	game           *game.Lauout15
	layout         *eui.GridLayoutRightDown
	fn             func(*eui.Button)
	deck           *game.DeckCards52
	sw             *eui.Stopwatch
	historyOfMoves [][]*game.Card
	moveIdx        int
}

func NewBoardSol15(fn func(*eui.Button)) *BoardSol15 {
	b := &BoardSol15{}
	b.fn = fn
	b.deck = game.NewDeckCards52().Shuffle()
	b.game = game.NewLayout15(b.deck)
	b.layout = eui.NewGridLayoutRightDown(14, 8)
	b.sw = eui.NewStopwatch()
	b.Visible(true)
	return b
}

func (b *BoardSol15) Setup(full bool) {
	if full {
		b.deck = game.NewDeckCards52().Shuffle()
		b.sw.Reset()
	}
	b.game.Reset(b.deck)
	b.layout.ResetContainerBase()
	for i := 0; i < 15; i++ {
		for p := 0; p < 4; p++ {
			cell := b.game.Row(i)[p]
			cardIcon := NewCardIcon(cell, b.fn)
			cell.Attach(cardIcon)
			b.layout.Add(cardIcon)
		}
		if (i+1)%3 == 0 && (i > 0) {
			for i := 0; i < 14; i++ {
				lbl := eui.NewText("*")
				b.layout.Add(lbl)
			}
		} else {
			lbl := eui.NewText("*")
			b.layout.Add(lbl)
		}
	}
	b.historyOfMoves = nil
	b.moveIdx = 0
	b.backupGame()
	b.sw.Start()
	b.Resize(b.GetRect().GetArr()) // обязательно после обнуления контейнеров
}

func (b *BoardSol15) MakeMove(move int) {
	if b.game.MakeMove(move) {
		if b.game.IsSolved() {
			b.sw.Stop()
			sb := eui.NewSnackBar("Пасьянс собран за " + b.sw.String() + ". Победа!!!").Show(5000)
			b.Add(sb)
		}
		b.moveIdx++
		b.backupGame()
	}
}

func (b *BoardSol15) backupGame() {
	deck := b.game.GetDeck()
	b.historyOfMoves = b.historyOfMoves[:b.moveIdx]
	b.historyOfMoves = append(b.historyOfMoves, deck)
	// b.moveIdx = len(b.historyOfMoves) - 1
	fmt.Println("deck:", deck, b.moveIdx, len(b.historyOfMoves))
}

func (b *BoardSol15) Update(dt int) {
	if !b.IsVisible() {
		return
	}
	for _, v := range b.layout.GetContainer() {
		v.Update(dt)
	}
	b.DrawableBase.Update(dt)
}

func (b *BoardSol15) Draw(surface *ebiten.Image) {
	if !b.IsVisible() {
		return
	}
	for _, v := range b.layout.GetContainer() {
		v.Draw(surface)
	}
	b.DrawableBase.Draw(surface)
}

func (b *BoardSol15) Resize(rect []int) {
	b.Rect(eui.NewRect(rect))
	margin := float64(b.GetRect().GetLowestSize()) * 0.005
	x, y, w, h := b.GetRect().GetRect()
	b.layout.Resize([]int{x, y, w, h})
	b.layout.SetCellMargin(margin)
}
