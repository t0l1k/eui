package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/utils"
	"golang.org/x/image/colornames"
)

func main() {
	title := "RGB Calculator"
	eui.Init(eui.GetUi().SetTitle(title).SetSize(400, 300))
	eui.Run(func() *eui.Scene {
		var (
			r, g, b, a *eui.Slider
		)
		s := eui.NewScene(eui.NewLayoutVerticalPercent([]int{5, 90}, 5))

		calcColor := eui.NewLabel("Color")
		rLbl := eui.NewLabel("R:")
		gLbl := eui.NewLabel("G:")
		bLbl := eui.NewLabel("B:")
		aLbl := eui.NewLabel("A:")

		colorFrom := func(r, g, b, a float64) color.Color {
			rValue := utils.ValueFromPercent(r*100, 255)
			gValue := utils.ValueFromPercent(g*100, 255)
			bValue := utils.ValueFromPercent(b*100, 255)
			aValue := utils.ValueFromPercent(a*100, 255)
			col := color.RGBA{uint8(rValue), uint8(gValue), uint8(bValue), uint8(aValue)}
			return col
		}

		setCalcColor := func() {
			col := colorFrom(r.Value(), g.Value(), b.Value(), a.Value())
			calcColor.SetBg(col)
			for k, v := range colornames.Map {
				if v == col {
					calcColor.SetText(k)
					return
				}
			}
			calcColor.SetText("Color")
		}

		setSliderLbls := func(float64) {
			rLbl.SetText(fmt.Sprintf("R:%.0f", utils.ValueFromPercent(r.Value()*100, 255)))
			gLbl.SetText(fmt.Sprintf("G:%.0f", utils.ValueFromPercent(g.Value()*100, 255)))
			bLbl.SetText(fmt.Sprintf("B:%.0f", utils.ValueFromPercent(b.Value()*100, 255)))
			aLbl.SetText(fmt.Sprintf("A:%.0f", utils.ValueFromPercent(a.Value()*100, 255)))
			setCalcColor()
		}

		r = eui.NewSlider(0, 255, 0.5, eui.Horizontal, func(data float64) {
			setSliderLbls(data)
		})
		g = eui.NewSlider(0, 255, 0.5, eui.Horizontal, func(data float64) {
			setSliderLbls(data)
		})
		b = eui.NewSlider(0, 255, 0.5, eui.Horizontal, func(data float64) {
			setSliderLbls(data)
		})
		a = eui.NewProgress(0, 255, 1, eui.Horizontal, func(data float64) {})

		setSliderLbls(0)

		btnNext := eui.NewButton("Next random color", func(btn *eui.Button) {
			rV := rand.Intn(255)
			gV := rand.Intn(255)
			bV := rand.Intn(255)
			r.SetValue(utils.PercentOf(float64(rV), 255) * 0.01)
			g.SetValue(utils.PercentOf(float64(gV), 255) * 0.01)
			b.SetValue(utils.PercentOf(float64(bV), 255) * 0.01)
			setCalcColor()
		})

		calcCont := eui.NewContainer(eui.NewVBoxLayout(3))
		rgbCont := eui.NewContainer(eui.NewVBoxLayout(3))
		rCont := eui.NewContainer(eui.NewLayoutHorizontalPercent([]int{20, 80}, 3)).Add(rLbl).Add(r)
		gCont := eui.NewContainer(eui.NewLayoutHorizontalPercent([]int{20, 80}, 3)).Add(gLbl).Add(g)
		bCont := eui.NewContainer(eui.NewLayoutHorizontalPercent([]int{20, 80}, 3)).Add(bLbl).Add(b)
		aCont := eui.NewContainer(eui.NewLayoutHorizontalPercent([]int{20, 80}, 3)).Add(aLbl).Add(a)

		s.Add(eui.NewTopBar(title, nil).SetShowStoppwatch(true).SetUseStopwatch())

		calcCont.Add(calcColor)

		rgbCont.Add(rCont)
		rgbCont.Add(gCont)
		rgbCont.Add(bCont)
		rgbCont.Add(aCont)

		calcCont.Add(rgbCont)
		calcCont.Add(btnNext)

		s.Add(calcCont)
		return s
	}())
}
