package app

import "github.com/t0l1k/eui"

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
	MakeMove(int)
}
