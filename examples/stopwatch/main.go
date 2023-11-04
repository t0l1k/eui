package main

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/stopwatch/app"
	"github.com/t0l1k/eui/examples/stopwatch/app/scenes"
)

func main() {
	eui.Init(app.NewGame())
	eui.Run(scenes.NewSceneStopwatch())
	eui.Quit()
}
