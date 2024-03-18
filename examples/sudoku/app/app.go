package app

import "github.com/t0l1k/eui"

func NewGameSudoku() *eui.Ui {
	u := eui.GetUi()
	u.SetTitle("Собери поле")
	k := 60
	w, h := 9*k, 6*k
	u.SetSize(w, h)
	u.GetTheme().Set(eui.ViewBg, eui.Black)
	return u
}
