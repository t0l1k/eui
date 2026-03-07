package temper

import (
	"fmt"
	"strconv"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/7Guis/data"
)

func NewTemperatureDemo(fn func(*eui.Button)) *eui.Container {
	var (
		cField, fField *eui.InputLine
	)
	cField = eui.NewInputLine(func(til *eui.InputLine) {}).SetPlaceholder("Цельсия")
	cField.TextChanged().Connect(func(data string) {
		v, err := strconv.ParseFloat(data, 64)
		// v, err := til.Digit()
		if err == nil {
			result := v*9/5 + 32.0
			fField.SetText(fmt.Sprintf("%.1f", result))
		}
	})

	fField = eui.NewInputLine(func(til *eui.InputLine) {}).SetPlaceholder("Фаренгейта")
	fField.TextChanged().Connect(func(data string) {
		v, err := strconv.ParseFloat(data, 64)
		if err == nil {
			result := (v - 32) * 5 / 9.0
			cField.SetText(fmt.Sprintf("%.1f", result))
		}
	})

	contC := eui.NewContainer(eui.NewVBoxLayout(1))
	contC.Add(cField)
	contC.Add(eui.NewLabel("градус Цельсия"))

	contF := eui.NewContainer(eui.NewVBoxLayout(1))
	contF.Add(fField)
	contF.Add(eui.NewLabel("градус Фаренгейта"))

	cont := eui.NewContainer(eui.NewLayoutHorizontalPercent([]int{45, 10, 45}, 1))
	cont.Add(contC)
	cont.Add(eui.NewLabel("="))
	cont.Add(contF)

	tempScene := eui.NewContainer(eui.NewLayoutVerticalPercent([]int{5, 5, 90}, 1))
	tempScene.Add(eui.NewTopBar(data.TempConv, fn).SetButtonText(data.QuitDemo))
	tempScene.Add(eui.NewLabel("Конвертер температур градус Цельсия/Фаренгейта"))
	tempScene.Add(cont)
	return tempScene
}
