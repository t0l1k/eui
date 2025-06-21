package app

import (
	"github.com/t0l1k/eui"
	"golang.org/x/image/colornames"
)

var Title = "Собери поле sudoku"

func NewGameSudoku() *eui.Ui {
	k := 90
	w, h := 9*k, 6*k
	u := eui.GetUi().SetTitle(Title).SetSize(w, h)
	u.GetTheme().Set(eui.ViewBg, colornames.Black)
	return u
}
