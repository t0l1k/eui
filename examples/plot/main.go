package main

import (
	"log"

	"github.com/t0l1k/eui"
)

func main() {
	eui.Init(eui.GetUi().SetTitle("Test Plot").SetSize(400, 400))
	eui.Run(func() *eui.Scene {
		s := eui.NewScene(eui.NewStackLayout(50))
		xArr := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		yArr := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		values := []float64{1, 2, 3, 4, 5, 5, 6, 7, 7, 8}
		plot := eui.NewPlot(xArr, yArr, values, "Game Score", "Game", "Level").
			AddValues([]float64{0, 1, 2, 3, 4, 4, 5, 6, 6, 7}).
			AddValues(func() (result []float64) {
				sum := 0.0
				for i, v := range values {
					sum += v
					result = append(result, sum/float64(i+1))
				}
				log.Println("res:", values, result)
				return result
			}())
		s.Add(plot)
		return s
	}())
	eui.Quit(func() {})
}
