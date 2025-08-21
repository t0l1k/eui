package app

import (
	"fmt"
	"strconv"
	"time"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/solitaire/sols"
	"github.com/t0l1k/eui/examples/games/solitaire/sols/deck"
	"github.com/t0l1k/eui/examples/games/solitaire/sols/sol15"
)

type BoardSol15 struct {
	*eui.Container
	game           *sol15.Lauout15
	layout         *eui.Container
	fn             func(*eui.Button)
	deck           *deck.DeckCards52
	sw             *eui.Stopwatch
	historyOfMoves [][]*deck.Card
	moveIdx        int
}

func NewBoardSol15(fn func(*eui.Button)) *BoardSol15 {
	b := &BoardSol15{Container: eui.NewContainer(eui.NewAbsoluteLayout())}
	b.fn = fn
	b.deck = deck.NewDeckCards52().Shuffle()
	b.game = sol15.NewLayout15(b.deck)
	b.layout = eui.NewContainer(eui.NewGridLayout(14, 8, 1))
	b.Add(b.layout)
	b.sw = eui.NewStopwatch()
	return b
}

func (b *BoardSol15) Setup(newDeck bool) {
	if newDeck {
		b.deck = deck.NewDeckCards52().Shuffle()
		b.sw.Reset()
	}
	b.game.Reset(b.deck)
	b.layout.ResetContainer()
	for i := 0; i < 15; i++ {
		for p := 0; p < 4; p++ {
			cell := b.game.Column(sols.Column(i))[p]
			cardIcon := NewCardIcon(cell, b.fn)
			cell.Connect(cardIcon.UpdateData)
			b.layout.Add(cardIcon)
		}
		if (i+1)%3 == 0 && (i > 0 && i < 14) {
			for i := 0; i < 14; i++ {
				lbl := eui.NewLabel(" ")
				b.layout.Add(lbl)
			}
		} else if i < 14 {
			lbl := eui.NewLabel(" ")
			b.layout.Add(lbl)
		}
	}
	b.historyOfMoves = nil
	b.moveIdx = 0
	b.backupGame()
	b.sw.Start()
	b.SetRect(b.Rect()) // обязательно после обнуления контейнеров
	if b.Rect().IsEmpty() {
		return
	}
	b.Layout()
}

func (b *BoardSol15) MakeMove(move sols.Column) {
	if b.game.MakeMove(move) {
		if b.game.IsSolved() {
			b.sw.Stop()
			eui.NewSnackBar("Пасьянс собран за " + eui.FormatSmartDuration(b.sw.Duration(), false) + ". Победа!!!").ShowTime(5 * time.Second)
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

func (b *BoardSol15) SetRect(rect eui.Rect[int]) {
	b.Container.SetRect(rect)
	b.layout.SetRect(rect)
}
