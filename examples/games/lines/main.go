package main

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/lines/app"
	"github.com/t0l1k/eui/examples/games/lines/app/scene"
)

func main() {
	eui.Init(app.NewGameLines())
	eui.Run(scene.NewSceneGame())
	eui.Quit(func() {})
}
