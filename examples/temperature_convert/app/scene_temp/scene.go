package scene_temp

import (
	"fmt"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/temperature_convert/app/model"
)

type SceneTemp struct {
	eui.SceneBase
}

func NewSceneTemp() *SceneTemp {
	var a, c *eui.InputBox
	s := &SceneTemp{}
	s.BoxLayout.SetHorizontal()
	a = eui.NewInputBox("0", 10, func(ib *eui.InputBox) {
		n := model.GetFahrenheitFromCelsius(ib.GetDigit())
		c.SetDigit(fmt.Sprintf("%.2f", n))
		fmt.Println(c.GetDigit(), n)
	})
	a.Name("Input Celsius")
	s.Add(a)
	b := eui.NewText("Celsius =")
	b.Name("Text Celsius")
	s.Add(b)
	c = eui.NewInputBox("0", 10, func(ib *eui.InputBox) {
		n := model.GetCelsiusFromFahrenheit(ib.GetDigit())
		a.SetDigit(fmt.Sprintf("%.2f", n))
		fmt.Println(a.GetDigit(), n)
	})
	c.Name("Input Fahrenheit")
	s.Add(c)
	d := eui.NewText("Fahrenheit")
	d.Name("Text Fahrenheit")
	s.Add(d)
	s.Resize()
	return s
}
