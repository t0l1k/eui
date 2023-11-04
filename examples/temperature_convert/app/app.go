package app

import "github.com/t0l1k/eui"

func NewGame() *eui.Ui {
	u := eui.GetUi()
	u.SetTitle("Convert temperature")
	k := 2
	w, h := 320*k, 200*k
	u.SetSize(w, h)
	return u
}
