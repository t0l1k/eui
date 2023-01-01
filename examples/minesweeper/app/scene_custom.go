package app

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	ui "github.com/t0l1k/eui"
)

type SceneCustom struct {
	ui.ContainerDefault
	topBar                        *TopBar
	optRows, optColumns, optMines *ui.Combobox
	btnStart                      *ui.Button
}

func NewSceneCustom() *SceneCustom {
	s := &SceneCustom{}
	s.topBar = NewTopBar(ui.GetUi().GetTitle() + "Select Custom Game")
	s.Add(s.topBar)
	rect := []int{0, 0, 1, 1}
	fg := ui.Red
	bg := ui.Blue
	var numbers []interface{}
	for i := 5; i < 100; i += 5 {
		numbers = append(numbers, i)
	}
	var rows, columns, mines int
	s.optRows = ui.NewCombobox("Rows", rect, bg, fg, numbers, 0, func(c *ui.Combobox) {
		rows = s.optRows.Value().(int)
	})
	s.Add(s.optRows)
	s.optColumns = ui.NewCombobox("Columns", rect, bg, fg, numbers, 0, func(c *ui.Combobox) {
		columns = s.optColumns.Value().(int)
	})
	s.Add(s.optColumns)
	s.optMines = ui.NewCombobox("Mines", rect, bg, fg, numbers, 0, func(c *ui.Combobox) {
		mines = s.optMines.Value().(int)
	})
	s.Add(s.optMines)
	s.btnStart = ui.NewButton("Start", rect, bg, fg, func(b *ui.Button) {
		sc := NewSceneGame(rows, columns, mines)
		ui.Push(sc)
		sc.topBar.lblTitle.SetText(fmt.Sprintf("%v Custom(%v,%v,%v)", ui.GetUi().GetTitle(), rows, columns, mines))

	})
	s.Add(s.btnStart)
	return s
}

func (s *SceneCustom) Entered() {
	s.Resize()
}

func (s *SceneCustom) Update(dt int) {
	for _, c := range s.Container {
		c.Update(dt)
	}
}

func (s *SceneCustom) Draw(surface *ebiten.Image) {
	for _, c := range s.Container {
		c.Draw(surface)
	}
}

func (s *SceneCustom) Resize() {
	s.topBar.Resize()
	w, h := ebiten.WindowSize()
	hTop := int(float64(h) * 0.05)
	rect := ui.NewRect([]int{0, hTop, w, h - hTop})
	w1, h1 := int(float64(w)*0.6), rect.H/6
	x, y := rect.CenterX()-w1/2, hTop
	y += h1
	s.optRows.Resize([]int{x, y, w1, h1 - 2})
	y += h1
	s.optColumns.Resize([]int{x, y, w1, h1 - 2})
	y += h1
	s.optMines.Resize([]int{x, y, w1, h1 - 2})
	y += h1
	s.btnStart.Resize([]int{x, y, w1, h1 - 2})
}

func (s *SceneCustom) Close() {
	for _, c := range s.Container {
		c.Close()
	}
}
