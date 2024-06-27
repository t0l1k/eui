package app

import "github.com/t0l1k/eui"

var title = "Собери поле sudoku"

func NewGameSudoku() *eui.Ui {
	u := eui.GetUi()
	u.SetTitle(title)
	k := 90
	w, h := 9*k, 6*k
	u.SetSize(w, h)
	u.GetTheme().Set(eui.ViewBg, eui.Black)
	return u
}
