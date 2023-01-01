package main

import (
	ui "github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/minesweeper/app"
)

func main() {
	ui.Init(app.NewGame())
	ui.Run(app.NewSceneMain())
	ui.Quit()
}
