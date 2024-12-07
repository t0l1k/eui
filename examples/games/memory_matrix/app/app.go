package app

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/colors"
)

const (
	title = "Вспомни Матрицу"
)

func NewApp() *eui.Ui {
	u := eui.GetUi()
	u.SetTitle(title)
	k := 80
	w, h := 9*k, 6*k
	u.SetSize(w, h)
	u.GetTheme().Set(eui.ViewBg, colors.Navy)
	return u
}
