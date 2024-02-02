package main

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/lines/app"
	"github.com/t0l1k/eui/examples/lines/app/scene"
)

func main() {
	eui.Init(app.NewGameLines())
	eui.Run(scene.NewSceneGame())
	eui.Quit()
}
