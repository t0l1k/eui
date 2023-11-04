package main

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/analog_clock/app"
	"github.com/t0l1k/eui/examples/analog_clock/app/scene_clock"
)

func main() {
	eui.Init(app.NewGame())
	eui.Run(scene_clock.NewSceneAnalogClock())
	eui.Quit()
}
