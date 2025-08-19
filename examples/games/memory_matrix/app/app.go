package app

import (
	"github.com/t0l1k/eui"
	"golang.org/x/image/colornames"
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
	u.Theme().Set(eui.ViewBg, colornames.Navy)
	return u
}
