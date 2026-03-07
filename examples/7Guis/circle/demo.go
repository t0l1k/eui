package circle

import (
	"log"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/7Guis/circle/undo"
	"github.com/t0l1k/eui/examples/7Guis/data"
)

func NewCircleDemo(fn func(*eui.Button)) *eui.Container {
	var (
		btnUndo, btnRedo *eui.Button
	)
	history := undo.NewUndoer(0)
	canvas := NewCanvas()

	// Сохраняем снимок состояния всех кругов (только данные, не объекты)
	saveCanvasState := func() {
		circles := make([]CircleModel, 0)
		for _, child := range canvas.Children() {
			if c, ok := child.(*Circle); ok {
				circles = append(circles, c.Model())
			}
		}
		history.Save(circles)
	}

	// Восстанавливаем состояние из снимка
	restoreCanvasState := func() {
		state := history.State()
		if state != nil {
			if models, ok := state.([]CircleModel); ok {
				// Обновляем существующие круги или создаем новые
				children := canvas.Children()

				// Если кол-во кругов совпадает - обновляем существующие
				if len(children) == len(models) {
					for i, model := range models {
						if c, ok := children[i].(*Circle); ok {
							c.Reset(model.x, model.y, model.r)
						}
					}
				} else {
					// Если не совпадает - очищаем и пересоздаем
					canvas.ResetContainer()
					for _, model := range models {
						circle := NewCircle(func(data *Circle) {
							saveCanvasState()
						})
						circle.Reset(model.x, model.y, model.r)
						canvas.Add(circle)
					}
				}
				canvas.MarkDirty()
				log.Println("Restored:", len(models), "circles")
			}
		}
	}

	btnUndo = eui.NewButton("Undo", func(b *eui.Button) {
		history.Undo()
		restoreCanvasState()
		log.Println("Undo executed")
	})

	btnRedo = eui.NewButton("Redo", func(b *eui.Button) {
		history.Redo()
		restoreCanvasState()
		log.Println("Redo executed")
	})

	canvas.clicked.Connect(func(data eui.Point[int]) {
		circle := NewCircle(func(data *Circle) {
			// При изменении параметров круга - сохраняем состояние
			saveCanvasState()
		})
		circle.Reset(data.X, data.Y, 30)
		canvas.Add(circle)
		saveCanvasState()
		log.Println("Circle added at", data)
	})

	// Сохраняем начальное состояние (пустой холст)
	saveCanvasState()

	btns := eui.NewContainer(eui.NewHBoxLayout(5)).Add(btnUndo).Add(btnRedo)
	cont := eui.NewContainer(eui.NewLayoutVerticalPercent([]int{5, 95}, 1)).Add(btns).Add(canvas)
	scene := eui.NewContainer(eui.NewLayoutVerticalPercent([]int{5, 95}, 1))
	scene.Add(eui.NewTopBar(data.Circle, fn).SetButtonText(data.QuitDemo)).Add(cont)
	return scene
}
