package app

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/solitaire/sols"
	"github.com/t0l1k/eui/examples/games/solitaire/sols/deck"
)

var title = "Собери пасьянс"

func NewGame() *eui.Ui {
	u := eui.GetUi()
	u.SetTitle(title)
	k := 90
	w, h := 9*k, 6*k
	u.SetSize(w, h)
	return u
}

type Sols interface {
	eui.Drawabler
	Setup(bool)
	Game() sols.CardGame
	MakeMove(sols.Column)
	AvailableMoves() (int, string)
	Stopwatch() *eui.Stopwatch
	GetHistory() [][]*deck.Card
	GetMoveNr() int
	SetMoveNr(int)
}
