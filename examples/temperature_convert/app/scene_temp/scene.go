package scene_temp

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/temperature_convert/app/model"
)

type SceneTemp struct {
	eui.SceneBase
	layout *eui.BoxLayout
}

func NewSceneTemp() *SceneTemp {
	var a, c *eui.InputBox
	s := &SceneTemp{}
	s.layout = eui.NewHLayout()
	a = eui.NewInputBox("0", 10, func(ib *eui.InputBox) {
		n := model.GetFahrenheitFromCelsius(ib.GetDigit())
		c.SetDigit(fmt.Sprintf("%.2f", n))
		fmt.Println(c.GetDigit(), n)
	})
	s.layout.Add(a)
	b := eui.NewText("Celsius =")
	s.layout.Add(b)
	c = eui.NewInputBox("0", 10, func(ib *eui.InputBox) {
		n := model.GetCelsiusFromFahrenheit(ib.GetDigit())
		a.SetDigit(fmt.Sprintf("%.2f", n))
		fmt.Println(a.GetDigit(), n)
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
