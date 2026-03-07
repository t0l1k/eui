package circle

import (
	"log"

	"github.com/t0l1k/eui"
	"golang.org/x/image/colornames"
)

func NewEditCircleDialog(c *Circle) eui.Drawabler {
	w, h := 300, 150
	uw, uh := eui.GetUi().Size()
	x, y := (uw-w)/2, (uh-h)/2

	cont := eui.NewContainer(eui.NewVBoxLayout(10))
	cont.SetRect(eui.NewRect([]int{x, y, w, h}))
	cont.SetBg(colornames.Gray)

	lbl := eui.NewLabel("Diameter")

	initialD := float64(c.r) * 0.01
	var tempD float64 = initialD

	update := func() {
		c.r = int(tempD * 100)
		c.SetRect(eui.NewRect([]int{c.x - c.r, c.y - c.r, c.r * 2, c.r * 2}))
	}

	slider := eui.NewSlider(0.05, 1, initialD, eui.Horizontal, func(v float64) {
		tempD = v
		update()
	})

	btns := eui.NewContainer(eui.NewHBoxLayout(10))
	btns.Add(eui.NewButton("Cancel", func(b *eui.Button) {
		tempD = initialD
		update()
		eui.GetUi().HideModal()
	}))
	btns.Add(eui.NewButton("Apply", func(b *eui.Button) {
		update()
		c.changed.Emit(c)
		log.Println("Update radius", c)
		eui.GetUi().HideModal()
	}))
	return cont.Add(lbl).Add(slider).Add(btns)
}
