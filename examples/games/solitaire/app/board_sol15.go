package app

import (
	"fmt"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/solitaire/sols"
	"github.com/t0l1k/eui/examples/games/solitaire/sols/deck"
	"github.com/t0l1k/eui/examples/games/solitaire/sols/sol15"
)

type BoardSol15 struct {
	eui.DrawableBase
	game           *sol15.Lauout15
	layout         *eui.GridLayoutRightDown
	fn             func(*eui.Button)
	deck           *deck.DeckCards52
	sw             *eui.Stopwatch
	historyOfMoves [][]*deck.Card
	moveIdx        int
}

func NewBoardSol15(fn func(*eui.Button)) *BoardSol15 {
	b := &BoardSol15{}
	b.fn = fn
	b.deck = deck.NewDeckCards52().Shuffle()
	b.game = sol15.NewLayout15(b.deck)
	b.layout = eui.NewGridLayoutRightDown(14, 8)
	b.sw = eui.NewStopwatch()
	return b
}

func (b *BoardSol15) Setup(newDeck bool) {
	if newDeck {
		b.deck = deck.NewDeckCards52().Shuffle()
		b.sw.Reset()
	}
	b.game.Reset(b.deck)
	b.layout.ResetContainerBase()
	for i := 0; i < 15; i++ {
		for p := 0; p < 4; p++ {
			cell := b.game.Column(sols.Column(i))[p]
			cardIcon := NewCardIcon(cell, b.fn)
			cell.Connect(cardIcon.UpdateData)
			b.layout.Add(cardIcon)
		}
		if (i+1)%3 == 0 && (i > 0 && i < 14) {
			for i := 0; i < 14; i++ {
				lbl := eui.NewText(" ")
				b.layout.Add(lbl)
			}
		} else if i < 14 {
			lbl := eui.NewText(" ")
			b.layout.Add(lbl)
		}
	}
	b.historyOfMoves = nil
	b.moveIdx = 0
	b.backupGame()
	b.sw.Start()
	b.Resize(b.GetRect().GetArr()) // обязательно после обнуления контейнеров
}

func (b *BoardSol15) MakeMove(move sols.Column) {
	if b.game.MakeMove(move) {
		if b.game.IsSolved() {
			b.sw.Stop()
			eui.NewSnackBar("Пасьянс собран за " + b.sw.String() + ". Победа!!!").Show(5000)
		}
		b.moveIdx++
		b.backupGame()
	}
}

func (b *BoardSol15) Game() sols.CardGame { return b.game }

func (b *BoardSol15) AvailableMoves() (int, string) {
	moves := b.game.AvailableMoves()
	str := "Ход:" + strconv.Itoa(len(b.historyOfMoves)) + " доступно:" + strconv.Itoa(moves) + " ходов"
	return moves, str
}

func (b *BoardSol15) Stopwatch() *eui.Stopwatch  { return b.sw }
func (b *BoardSol15) GetHistory() [][]*deck.Card { return b.historyOfMoves }
func (b *BoardSol15) GetMoveNr() int             { return b.moveIdx }
func (b *BoardSol15) SetMoveNr(value int)        { b.moveIdx = value }

func (b *BoardSol15) backupGame() {
	deck := b.game.GetDeck()
	b.historyOfMoves = b.historyOfMoves[:b.moveIdx]
	b.historyOfMoves = append(b.historyOfMoves, deck)
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
