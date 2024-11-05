package sols

import "github.com/t0l1k/eui/examples/games/solitaire/sols/deck"

type CardGame interface {
	Reset(*deck.DeckCards52)
	GetDeck() []*deck.Card
	SetDeck([]*deck.Card)
	Index(*deck.Card) (Column, int)
	MakeMove(Column) bool
	AvailableMoves() int
	IsSolved() bool
}
