package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
)

func GetCelsiusFromFahrenheit(f float64) float64 {
	return (5.0 / 9) * (f - 32)
}

func GetFahrenheitFromCelsius(c float64) float64 {
	return c*9/5 + 32
}

type SceneTemp struct {
	eui.SceneBase
	layout *eui.BoxLayout
}

func NewSceneTemp() *SceneTemp {
	var a, c *eui.InputBox
	s := &SceneTemp{}
	s.layout = eui.NewHLayout()
	a = eui.NewDigitInputBox("01234", 5, func(ib *eui.InputBox) {
		if digit, err := ib.GetDigit(); err == nil {
			n := GetFahrenheitFromCelsius(digit)
			c.SetText(fmt.Sprintf("%.2f", n))
			fmt.Println(digit, n, a.GetText())
		}
	})
	s.layout.Add(a)
	b := eui.NewText("Celsius =")
	s.layout.Add(b)
	c = eui.NewDigitInputBox("43210", 5, func(ib *eui.InputBox) {
		if digit, err := ib.GetDigit(); err == nil {
			n := GetCelsiusFromFahrenheit(digit)
			a.SetText(fmt.Sprintf("%.2f", n))
			fmt.Println(digit, n, c.GetText())
		}
	})
	s.layout.Add(c)
	d := eui.NewText("Fahrenheit")
	s.layout.Add(d)
	s.Resize()
	return s
}

func (s *SceneTemp) Update(dt int) {
	for _, v := range s.layout.GetContainer() {
		v.Update(dt)
	}
}

func (s *SceneTemp) Draw(surface *ebiten.Image) {
	for _, v := range s.layout.GetContainer() {
		v.Draw(surface)
	}
}

func (s *SceneTemp) Resize() {
	w, h := eui.GetUi().Size()
	s.layout.Resize([]int{0, 0, w, h})
}

func NewGame() *eui.Ui {
	u := eui.GetUi()
	u.SetTitle("Convert temperature")
	k := 2
	w, h := 320*k, 200*k
	u.SetSize(w, h)
	return u
}

func main() {
	eui.Init(NewGame())
	eui.Run(NewSceneTemp())
	eui.Quit(func() {})
}
