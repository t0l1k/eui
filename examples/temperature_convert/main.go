package main

import (
	"fmt"

	"github.com/t0l1k/eui"
)

func GetCelsiusFromFahrenheit(f float64) float64 { return (5.0 / 9) * (f - 32) }
func GetFahrenheitFromCelsius(c float64) float64 { return c*9/5 + 32 }

type SceneTemp struct{ *eui.Scene }

func NewSceneTemp() *SceneTemp {
	var a, c *eui.InputBox
	s := &SceneTemp{Scene: eui.NewScene(eui.NewHBoxLayout(2))}
	a = eui.NewDigitInputBox("01234", 5, func(ib *eui.InputBox) {
		if digit, err := ib.GetDigit(); err == nil {
			n := GetFahrenheitFromCelsius(digit)
			c.SetText(fmt.Sprintf("%.2f", n))
			fmt.Println(digit, n, a.GetText())
		}
	})
	s.Add(a)
	s.Add(eui.NewText("Celsius ="))
	c = eui.NewDigitInputBox("43210", 5, func(ib *eui.InputBox) {
		if digit, err := ib.GetDigit(); err == nil {
			n := GetCelsiusFromFahrenheit(digit)
			a.SetText(fmt.Sprintf("%.2f", n))
			fmt.Println(digit, n, c.GetText())
		}
	})
	s.Add(c)
	s.Add(eui.NewText("Fahrenheit"))
	return s
}

func main() {
	eui.Init(eui.GetUi().SetTitle("Convert temperature").SetSize(320, 80))
	eui.Run(NewSceneTemp())
	eui.Quit(func() {})
}
