package main

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/minesweeper/app"
	"github.com/t0l1k/eui/examples/minesweeper/app/scenes/scene_main"
)

func main() {
	eui.Init(app.NewGame())
	eui.Run(scene_main.NewSceneSelectGame())
	eui.Quit()
}
