package scene_counter

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/counter/app/model"
)

type SceneCounter struct {
	eui.SceneBase
}

func NewSceneCounter() *SceneCounter {
	sc := &SceneCounter{}

	counter := eui.NewIntVar(model.Value())
	lblCounter := eui.NewText(counter.String())
	counter.Attach(lblCounter)
	sc.Add(lblCounter)

	btnInc := eui.NewButton("+", func(b *eui.Button) {
		model.Inc()
		counter.SetValue(model.Value())
	})

	btnDec := eui.NewButton("-", func(b *eui.Button) {
		model.Dec()
		counter.SetValue(model.Value())
	})

	horLayout := eui.NewHLayout()
	horLayout.Add(btnInc)
	horLayout.Add(btnDec)
	sc.Add(horLayout)
	sc.Resize()
	return sc
}
