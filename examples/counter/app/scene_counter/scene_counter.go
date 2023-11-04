package scene_counter

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/counter/app/model"
)

type SceneCounter struct {
	eui.SceneDefault
}

func NewSceneCounter() *SceneCounter {
	sc := &SceneCounter{}

	counter := eui.NewIntVar(model.Value())
	lblCounter := eui.NewText(counter.String(), eui.Yellow, eui.Black)
	counter.Attach(lblCounter)
	sc.Add(lblCounter)

	btnInc := eui.NewButton("+", eui.Green, eui.Black, func(b *eui.Button) {
		model.Inc()
		counter.Set(model.Value())
	})

	btnDec := eui.NewButton("-", eui.Red, eui.Black, func(b *eui.Button) {
		model.Dec()
		counter.Set(model.Value())
	})

	horLayout := eui.NewHLayout()
	horLayout.Add(btnInc)
	horLayout.Add(btnDec)
	sc.Add(horLayout)
	sc.Resize()
	return sc
}
