package timer

import (
	"fmt"
	"log"
	"time"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/7Guis/data"
	"github.com/t0l1k/eui/utils"
)

func NewTimerDemo(fn func(*eui.Button)) *eui.Container {
	var (
		duration, elapsed time.Duration
		running           bool
	)
	duration = 30 * time.Second
	label := eui.NewLabel("")
	progress := eui.NewProgress(0, 1, 0, eui.Horizontal, func(data float64) {})

	reset := func() {
		elapsed = 0
		progress.SetValue(0)
		label.SetText(fmt.Sprintf("%.1f", duration.Seconds()))
		log.Println("Обнулить")
	}

	slider := eui.NewSlider(0, 1, 0.30, eui.Horizontal, func(data float64) {
		duration = time.Duration(data*100) * time.Second
		if !running {
			reset()
		}
		log.Println("Slider", data)
	})
	button := eui.NewButton("Start", func(b *eui.Button) {
		running = !running
		if running {
			b.SetText("Stop")
			log.Println("Начали")
		} else {
			b.SetText("Start")
			log.Println("Остановили")
		}
	})

	eui.GetUi().TickListener().Connect(func(data eui.Event) {
		if !running {
			return
		}
		td := data.Value.(eui.TickData)
		dt := td.Duration()
		elapsed += dt
		if elapsed >= duration {
			elapsed = 0
			running = false
			button.SetText("Start")
			log.Println("Отсчитались")
		}
		progress.SetValue(utils.PercentOf(elapsed.Seconds(), duration.Seconds()) / 100)
		label.SetText(fmt.Sprintf("%.1f", (duration - elapsed).Seconds()))
	})

	cont := eui.NewContainer(eui.NewVBoxLayout(1))
	cont.Add(slider)
	cont.Add(progress)
	cont.Add(label)
	cont.Add(button)

	scene := eui.NewContainer(eui.NewLayoutVerticalPercent([]int{5, 95}, 1))
	scene.Add(eui.NewTopBar(data.Timer, fn).SetButtonText(data.QuitDemo))
	scene.Add(cont)
	return scene
}
