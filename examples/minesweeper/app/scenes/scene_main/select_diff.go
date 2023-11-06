package scene_main

import (
	"fmt"
	"strconv"

	"github.com/t0l1k/eui/examples/minesweeper/app/scenes/scene_game"

	"github.com/t0l1k/eui"
)

type SelectDiff struct {
	eui.SceneBase
	frame                          *eui.BoxLayout
	topBar                         *eui.TopBar
	comboRow, comboCol, comboMines *eui.ComboBox
	btnExec                        *eui.Button
}

func NewSelectDiff(title string) *SelectDiff {
	s := &SelectDiff{}

	s.topBar = eui.NewTopBar(title)
	s.Add(s.topBar)

	s.frame = eui.NewVLayout()

	lblTitle := eui.NewText("Настрой сложность")
	s.frame.Add(lblTitle)

	column, row, mines := 5, 5, 5
	var data []interface{}
	for i := 5; i < 100; i += 5 {
		data = append(data, i)
	}
	s.comboCol = eui.NewComboBox("Сколько рядов", data, 0, func(combo *eui.ComboBox) {
		row = combo.Value().(int)
	})
	s.frame.Add(s.comboCol)
	s.comboRow = eui.NewComboBox("Сколько столбиков", data, 0, func(combo *eui.ComboBox) {
		column = combo.Value().(int)
	})
	s.frame.Add(s.comboRow)
	s.comboMines = eui.NewComboBox("Сколько мин", data, 0, func(combo *eui.ComboBox) {
		mines = combo.Value().(int)
		fmt.Println("set mines", mines)
	})
	s.frame.Add(s.comboMines)
	s.btnExec = eui.NewButton("Запустить игру", func(b *eui.Button) {
		str := "Игра настроена на " + strconv.Itoa(column) + " рядов" + strconv.Itoa(row) + " рядов" + strconv.Itoa(mines) + " мин"
		game := scene_game.NewSceneGame(str, row, column, mines)
		eui.GetUi().Push(game)
	})
	s.frame.Add(s.btnExec)
	s.Add(s.frame)
	s.Resize()
	return s
}

func (s *SelectDiff) Resize() {
	w, h := eui.GetUi().Size()
	rect := eui.NewRect([]int{0, 0, w, h})
	hTop := int(float64(rect.GetLowestSize()) * 0.05)
	s.topBar.Resize([]int{0, 0, w, hTop})
	s.frame.Resize([]int{0, hTop, w, h - hTop})
}
