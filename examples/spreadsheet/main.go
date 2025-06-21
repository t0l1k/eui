package main

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/spreadsheet/app"
)

type SceneSpreadSheet struct {
	*eui.Scene
	ssView *app.SpreadsheetView
}

func NewSceneSpreadSheet() *SceneSpreadSheet {
	sc := &SceneSpreadSheet{Scene: eui.NewScene(eui.NewVBoxLayout(1))}
	sc.ssView = app.NewSpreadSheetView(5, 25)
	sc.Add(sc.ssView)
	sc.Resize()
	return sc
}

func NewGame() *eui.Ui {
	u := eui.GetUi()
	u.SetTitle("Spreadsheet example")
	k := 3
	w, h := 320*k, 200*k
	u.SetSize(w, h)
	return u
}

func main() {
	eui.Init(NewGame())
	eui.Run(NewSceneSpreadSheet())
	eui.Quit(func() {})
}
