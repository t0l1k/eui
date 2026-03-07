package cells

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/7Guis/data"
)

func NewCellsDemo(fn func(*eui.Button)) *eui.Container {
	scene := eui.NewContainer(eui.NewLayoutVerticalPercent([]int{5, 95}, 1))
	scene.Add(eui.NewTopBar(data.Cells, fn).SetButtonText(data.QuitDemo))
	cont := eui.NewContainer(eui.NewVBoxLayout(1))
	cont.Add(eui.NewLabel("В разработке..."))
	scene.Add(cont)
	return scene
}
