package main

import (
	"fmt"

	"github.com/t0l1k/eui"
)

func GetCelsiusFromFahrenheit(f float64) float64 { return (5.0 / 9) * (f - 32) }
func GetFahrenheitFromCelsius(c float64) float64 { return c*9/5 + 32 }

type SceneTemp struct{ *eui.Scene }

func NewSceneTemp() *SceneTemp {
	var a, c *eui.TextInputLine
	s := &SceneTemp{Scene: eui.NewScene(eui.NewHBoxLayout(2))}
	a = eui.NewTextInputLine(func(ib *eui.TextInputLine) {
		if digit, err := ib.Digit(); err == nil {
			n := GetFahrenheitFromCelsius(digit)
			c.SetText(fmt.Sprintf("%.2f", n))
			fmt.Println(digit, n, a.Text())
		}
	})
	s.Add(a)
	s.Add(eui.NewText("Celsius ="))
	c = eui.NewTextInputLine(func(ib *eui.TextInputLine) {
		if digit, err := ib.Digit(); err == nil {
			n := GetCelsiusFromFahrenheit(digit)
			a.SetText(fmt.Sprintf("%.2f", n))
			fmt.Println(digit, n, c.Text())
		}
	})
	s.Add(c)
	s.Add(eui.NewText("Fahrenheit"))
	return s
}

func main() {
	eui.Init(eui.GetUi().SetTitle("Convert temperature").SetSize(600, 120))
	eui.Run(NewSceneTemp())
	eui.Quit(func() {})
}
