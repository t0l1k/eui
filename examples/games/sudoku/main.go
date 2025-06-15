package main

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/sudoku/app"
)

func main() {
	eui.Init(app.NewGameSudoku())
	eui.Run(app.NewSceneSudoku())
	eui.Quit(func() {})
}
