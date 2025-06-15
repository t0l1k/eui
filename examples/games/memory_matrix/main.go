package main

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/memory_matrix/app"
)

func main() {
	eui.Init(app.NewApp())
	eui.Run(app.NewSceneMain())
	eui.Quit(func() {})
}
