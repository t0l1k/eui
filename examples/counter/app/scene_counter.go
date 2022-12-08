package app

import (
	"fmt"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	ui "github.com/t0l1k/eui"
)

type SceneCounter struct {
	counter                 Counter
	lblCounter              *ui.Label
	btnQuit, btnInc, btnDec *ui.Button
	ui.SceneDefault
}

func NewSceneCounter() *SceneCounter {
	sc := &SceneCounter{}
	rect := []int{0, 0, 1, 1}
	sc.btnQuit = ui.NewButton("<", rect, ui.GreenYellow, ui.Black, func(b *ui.Button) { ui.Pop() })
	sc.Add(sc.btnQuit)
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
	return sc
}

func (sc *SceneCounter) Add(d ui.Drawable) {
	sc.Container = append(sc.Container, d)
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
	sc.Rect = ui.NewRect([]int{0, 0, w, h})
	x, y, w, h := 0, 0, int(float64(sc.Rect.H)*0.05), int(float64(sc.Rect.H)*0.05)
	sc.btnQuit.Resize([]int{x, y, w, h})
	w, h = int(float64(sc.Rect.H)*0.3), int(float64(sc.Rect.H)*0.1)
	x, y = sc.Rect.W/2-w/2, y+int(float64(h)*1.1)
	sc.lblCounter.Resize([]int{x, y, w, h})
	w2 := int(float64(sc.Rect.W) * 0.2)
	w3 := int(float64(w2) * 0.1)
	fmt.Println(w, w2, w3)
	x, y = sc.Rect.W/2-(w2+w3), y+int(float64(h)*2)
	sc.btnInc.Resize([]int{x, y, w2, h})
	x = sc.Rect.W/2 + (w3)
	sc.btnDec.Resize([]int{x, y, w2, h})
}
func (sc *SceneCounter) Quit() {
	for _, v := range sc.Container {
		v.Close()
	}
}
