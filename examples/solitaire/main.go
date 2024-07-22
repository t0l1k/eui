package main

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/solitaire/app"
)

func main() {
	eui.Init(app.NewGame())
	eui.Run(app.NewSceneMain())
	eui.Quit()
}
