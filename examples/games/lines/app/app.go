package app

import (
	"github.com/t0l1k/eui"
	"golang.org/x/image/colornames"
)

const (
	GameTitle = "Очищай поле вовремя aka Color Lines"
	BNew      = "Новая Игра"

	FieldSizeSmall  = 7
	FieldSizeNormal = 9
	FieldSizeBig    = 11
)

func NewGameLines() *eui.Ui {
	u := eui.GetUi()
	u.SetTitle(GameTitle)
	k := 60
	w, h := 9*k, 6*k
	u.SetSize(w, h)
	u.Theme().Set(eui.ViewBg, colornames.Black)
	return u
}
