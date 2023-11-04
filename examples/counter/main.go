package main

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/counter/app"
	"github.com/t0l1k/eui/examples/counter/app/scene_counter"
)

func main() {
	eui.Init(app.NewGame())
	eui.Run(scene_counter.NewSceneCounter())
	eui.Quit()
}
