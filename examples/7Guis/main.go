package main

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/7Guis/mainscene"
)

func main() {
	eui.Init(eui.GetUi().SetTitle("7Guis").SetSize(800, 600))
	eui.Run(mainscene.NewMainScene())
}
