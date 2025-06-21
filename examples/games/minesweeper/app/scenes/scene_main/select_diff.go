package scene_main

import (
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/minesweeper/app/scenes/scene_game"
)

type SelectDiff struct {
	*eui.Scene
	frame                          *eui.Container
	topBar                         *eui.TopBar
	comboRow, comboCol, comboMines *eui.ComboBox
	btnExec                        *eui.Button
	percent, row, column           int
	kId                            int64
}

func NewSelectDiff(title string) *SelectDiff {
	s := &SelectDiff{Scene: eui.NewScene(eui.NewAbsoluteLayout())}
	s.topBar = eui.NewTopBar(title, nil)
	s.Add(s.topBar)
	s.frame = eui.NewContainer(eui.NewVBoxLayout(1))
	lblTitle := eui.NewText("Настрой сложность")
	s.frame.Add(lblTitle)
	s.column, s.row, s.percent = 5, 5, 15
	var data []interface{}
	for i := 5; i <= 50; i += 5 {
		data = append(data, i)
	}
	s.comboCol = eui.NewComboBox("Сколько рядов", data, 0, func(combo *eui.ComboBox) {
		s.row = combo.Value().(int)
	})
	s.frame.Add(s.comboCol)
	s.comboRow = eui.NewComboBox("Сколько столбиков", data, 0, func(combo *eui.ComboBox) {
		s.column = combo.Value().(int)
	})
	s.frame.Add(s.comboRow)
	s.comboMines = eui.NewComboBox("Сколько % мин", func() (arr []interface{}) {
		for i := 10; i <= 30; i++ {
			arr = append(arr, i)
		}
		return arr
	}(), 5, func(combo *eui.ComboBox) {
		s.percent = combo.Value().(int)
	})
	s.frame.Add(s.comboMines)
	s.btnExec = eui.NewButton("Запустить игру", func(b *eui.Button) {
		s.runGame()
	})
	s.frame.Add(s.btnExec)
	s.Add(s.frame)
	s.kId = eui.GetUi().GetInputKeyboard().Connect(s.UpdateInput)
	return s
}

func (s *SelectDiff) UpdateInput(ev eui.Event) {
	kd := ev.Value.(eui.KeyboardData)
	if kd.IsReleased(ebiten.KeySpace) {
		s.runGame()
	}
}

func (s *SelectDiff) runGame() {
	mines := s.percent * (s.row * s.column) / 100
	str := "Игра на " + strconv.Itoa(s.column) + " столбиков" + strconv.Itoa(s.row) + " рядов " + strconv.Itoa(mines) + " мин"
	game := scene_game.NewSceneGame(str, s.row, s.column, mines)
	eui.GetUi().Push(game)
	eui.GetUi().GetInputKeyboard().Disconnect(s.kId)
	log.Println("run game", s.row, s.column, mines)
}

func (s *SelectDiff) Resize() {
	w, h := eui.GetUi().Size()
	rect := eui.NewRect([]int{0, 0, w, h})
	s.SetRect(rect)
	hTop := int(float64(rect.GetLowestSize()) * 0.05)
	s.topBar.Resize(eui.NewRect([]int{0, 0, w, hTop}))
	s.frame.Resize(eui.NewRect([]int{0, hTop, w, h - hTop}))
}
