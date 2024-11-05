package app

import (
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/solitaire/sols"
	"github.com/t0l1k/eui/examples/games/solitaire/sols/deck"
	"github.com/t0l1k/eui/examples/games/solitaire/sols/fc"
)

type BoardFreecell struct {
	eui.DrawableBase
	layoutFC, layoutHome *eui.BoxLayout
	layoutCols           []*eui.BoxLayout
	game                 *fc.Freecell
	fn                   func(*eui.Button)
	deck                 *deck.DeckCards52
	sw                   *eui.Stopwatch
	historyOfMoves       [][]*deck.Card
	moveIdx              int
}

func NewBoardFreecell(fn func(*eui.Button)) *BoardFreecell {
	b := &BoardFreecell{}
	b.fn = fn
	b.deck = deck.NewDeckCards52()
	b.game = fc.NewFreecell(b.deck)
	b.layoutFC = eui.NewHLayout()
	b.layoutHome = eui.NewHLayout()
	for i := 0; i < 8; i++ {
		b.layoutCols = append(b.layoutCols, eui.NewVLayout())
	}
	b.sw = eui.NewStopwatch()
	return b
}

func (b *BoardFreecell) Setup(resetDeck bool) {
	if resetDeck {
		b.deck = deck.NewDeckCards52().Shuffle()
		b.sw.Reset()
	}
	b.game.Reset(b.deck)
	b.layoutFC.ResetContainerBase()
	b.layoutHome.ResetContainerBase()
	for i := 0; i < 8; i++ {
		b.layoutCols[i].ResetContainerBase()
	}

	var (
		idx int
	)
	for _, colName := range fc.ColFree {
		cell := b.game.Column(colName)[0]
		if cell.IsEmpty() {
			cardIcon := NewCardIcon(cell, b.fn)
			cell.Attach(cardIcon)
			b.layoutFC.Add(cardIcon)
			idx++
		}
	}
	for _, colName := range fc.ColHouses {
		cell := b.game.Column(colName)[0]
		cardIcon := NewCardIcon(cell, b.fn)
		cell.Attach(cardIcon)
		b.layoutHome.Add(cardIcon)
		idx++
	}

	for i, colName := range fc.Cols {
		cells := b.game.Column(colName)
		for _, cell := range cells {
			cardIcon := NewCardIcon(cell, b.fn)
			cell.Attach(cardIcon)
			b.layoutCols[i].Add(cardIcon)
			idx++
		}
	}
	log.Println("Setup", resetDeck, b.game)

	b.historyOfMoves = nil
	b.moveIdx = 0
	b.backupGame()
	b.sw.Start()
	b.Resize(b.GetRect().GetArr()) // обязательно после обнуления контейнеров
}

func (b *BoardFreecell) MakeMove(move sols.Column) {
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

func (b *BoardFreecell) Game() sols.CardGame { return b.game }

func (b *BoardFreecell) AvailableMoves() (int, string) {
	moves := b.game.AvailableMoves()
	str := "Ход:" + strconv.Itoa(len(b.historyOfMoves)) + " доступно:" + strconv.Itoa(moves) + " ходов"
	return moves, str
}

func (b *BoardFreecell) Stopwatch() *eui.Stopwatch  { return b.sw }
func (b *BoardFreecell) GetHistory() [][]*deck.Card { return b.historyOfMoves }
func (b *BoardFreecell) GetMoveNr() int             { return b.moveIdx }
func (b *BoardFreecell) SetMoveNr(value int)        { b.moveIdx = value }

func (b *BoardFreecell) backupGame() {
	deck := b.game.GetDeck()
	b.historyOfMoves = b.historyOfMoves[:b.moveIdx]
	b.historyOfMoves = append(b.historyOfMoves, deck)
}

func (b *BoardFreecell) Update(dt int) {
	if !b.IsVisible() {
		return
	}
	for _, v := range b.layoutFC.GetContainer() {
		v.Update(dt)
	}
	for _, v := range b.layoutHome.GetContainer() {
		v.Update(dt)
	}
	for _, layout := range b.layoutCols {
		for _, v := range layout.GetContainer() {
			v.Update(dt)
		}
	}
	b.DrawableBase.Update(dt)
}

func (b *BoardFreecell) Draw(surface *ebiten.Image) {
	if !b.IsVisible() {
		return
	}
	for _, v := range b.layoutFC.GetContainer() {
		v.Draw(surface)
	}
	for _, v := range b.layoutHome.GetContainer() {
		v.Draw(surface)
	}
	for _, layout := range b.layoutCols {
		for _, v := range layout.GetContainer() {
			v.Draw(surface)
		}
	}
	b.DrawableBase.Draw(surface)
}

func (b *BoardFreecell) Resize(rect []int) {
	b.Rect(eui.NewRect(rect))
	x0, y0, w0, h0 := b.GetRect().GetRect()
	cellSize := b.GetRect().GetLowestSize() / 8
	x := x0 + (w0-cellSize*8)/2
	y := y0 + (h0-cellSize*8)/2
	w, h := cellSize, cellSize
	b.layoutFC.Resize([]int{x, y, w * 4, h})
	b.layoutHome.Resize([]int{x + cellSize*4, y, w * 4, h})
	y += cellSize
	for i, layout := range b.layoutCols {
		layout.Resize([]int{x + cellSize*i, y, w, h * 8})
	}
}
