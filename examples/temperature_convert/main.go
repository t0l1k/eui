package main

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/temperature_convert/app"
	"github.com/t0l1k/eui/examples/temperature_convert/app/scene_temp"
)

func main() {
	eui.Init(app.NewGame())
	eui.Run(scene_temp.NewSceneTemp())
	eui.Quit()
}
