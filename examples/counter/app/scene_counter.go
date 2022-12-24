package app

import (
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	ui "github.com/t0l1k/eui"
)

type SceneCounter struct {
	topBar                      *TopBar
	counter                     Counter
	lblCounter                  *ui.Label
	btnInc, btnDec, btnNewScene *ui.Button
	ui.ContainerDefault
}

func NewSceneCounter() *SceneCounter {
	sc := &SceneCounter{}
	rect := []int{0, 0, 1, 1}
	sc.btnInc = ui.NewButton("+", rect, ui.Green, ui.Black, func(b *ui.Button) {
		sc.counter.Inc()
		sc.lblCounter.SetText(strconv.Itoa(int(sc.counter.Get())))
	})
	sc.Add(sc.btnInc)
	sc.btnDec = ui.NewButton("-", rect, ui.Red, ui.Black, func(b *ui.Button) {
		sc.counter.Dec()
		sc.lblCounter.SetText(strconv.Itoa(int(sc.counter.Get())))
	})
	sc.Add(sc.btnDec)
	sc.lblCounter = ui.NewLabel(sc.counter.String(), rect, ui.Yellow, ui.Black)
	sc.Add(sc.lblCounter)
	sc.topBar = NewTopBar()
	sc.Add(sc.topBar)
	sc.btnNewScene = ui.NewButton("New Scene", rect, ui.Green, ui.Black, func(b *ui.Button) {
		if b.IsMouseDownRight() {
			ui.Push(NewSceneGame())
		}
	})
	sc.Add(sc.btnNewScene)
	return sc
}

func (sc *SceneCounter) Entered() {
	sc.Resize()
}
func (sc *SceneCounter) Update(dt int) {
	for _, v := range sc.Container {
		v.Update(dt)
	}
}

func (sc *SceneCounter) Draw(surface *ebiten.Image) {
	for _, v := range sc.Container {
		v.Draw(surface)
	}
}

func (sc *SceneCounter) Resize() {
	w, h := ebiten.WindowSize()
	rect := ui.NewRect([]int{0, 0, w, h})
	w, h = int(float64(rect.H)*0.3), int(float64(rect.H)*0.1)
	x, y := rect.W/2-w/2, int(float64(h)*1.1)
	sc.lblCounter.Resize([]int{x, y, w, h})
	w2 := int(float64(rect.W) * 0.2)
	w3 := int(float64(w2) * 0.1)
	x, y = rect.W/2-(w2+w3), y+int(float64(h)*2)
	sc.btnInc.Resize([]int{x, y, w2, h})
	x = rect.W/2 + (w3)
	sc.btnDec.Resize([]int{x, y, w2, h})
	sc.topBar.Resize()
	x, y = rect.W/2-w/2, rect.H-h
	sc.btnNewScene.Resize([]int{x, y, w, h})
}

func (sc *SceneCounter) Close() {
	for _, v := range sc.Container {
		v.Close()
	}
}
