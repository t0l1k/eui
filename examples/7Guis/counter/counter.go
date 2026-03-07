package counter

import (
	"strconv"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/7Guis/data"
)

func NewCounterDemo(fn func(*eui.Button)) *eui.Container {
	label := eui.NewLabel("")
	counter := eui.NewSignal(func(a, b int) bool { return a == b })
	counter.ConnectAndFire(func(data int) {
		label.SetText(strconv.Itoa(data))
	}, 0)

	counterScene := eui.NewContainer(eui.NewLayoutVerticalPercent([]int{5, 70, 25}, 1))
	counterScene.Add(eui.NewTopBar(data.Counter, fn).SetButtonText(data.QuitDemo))
	counterScene.Add(label)
	counterScene.Add(eui.NewButton("Inc", func(b *eui.Button) {
		counter.Emit(counter.Value() + 1)
	}))
	return counterScene
}
