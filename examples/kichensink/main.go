package main

import (
	"log"
	"time"

	"github.com/t0l1k/eui"
)

func main() {
	title := "Витрина элементов"
	eui.Init(eui.GetUi().SetTitle(title).SetSize(800, 600))
	eui.Run(func() *eui.Scene {
		s := eui.NewScene(eui.NewLayoutVerticalPercent([]int{5, 95}, 10))

		conHor := eui.NewContainer(eui.NewHBoxLayout(3))
		conVSliders := eui.NewContainer(eui.NewVBoxLayout(3))
		sl01 := eui.NewSlider(0, 255, 0.25, eui.Horizontal, func(data float64) {
			log.Println("sl01 new value:", data)
		})
		sl02 := eui.NewSlider(0, 255, 0.25, eui.Vertical, func(data float64) {
			log.Println("sl02 new value:", data)
		})
		progressH := eui.NewProgress(0, 100, 0.5, eui.Horizontal, func(data float64) {})
		progressV := eui.NewProgress(0, 100, 0.5, eui.Vertical, func(data float64) {})

		var runningTask bool
		runProgress := func() {
			runningTask = true
			for runningTask && progressH.Value() <= 1.0 {
				time.Sleep(10 * time.Millisecond)
				newValue := progressH.Value() + 0.01
				progressH.SetValue(newValue)
				progressV.SetValue(newValue)
				if newValue > 1.0 {
					runningTask = false
					eui.NewSnackBar("Task finished!").ShowTime(3 * time.Second)
				}
				log.Println("runningTask", newValue, progressH.Value(), progressV.Value())
			}
		}

		btnProgress := eui.NewButton("Start", func(b *eui.Button) {
			progressH.SetValue(0)
			go runProgress()
		})

		s.Add(eui.NewTopBar(title, nil).SetShowStoppwatch(true).SetUseStopwatch())

		conVSliders.Add(sl01)
		conVSliders.Add(progressH)
		conVSliders.Add(btnProgress)
		conHor.Add(conVSliders)
		conHor.Add(sl02)
		conHor.Add(progressV)

		s.Add(conHor)
		return s
	}())
}
