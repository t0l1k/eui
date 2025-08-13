package scene_main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/minesweeper/app/scenes/scene_game"
)

type SelectDiff struct {
	*eui.Scene
	frame                *eui.Container
	topBar               *eui.TopBar
	cRow, cCol, cPerc    *eui.SpinBox[int]
	btnExec              *eui.Button
	percent, row, column int
	kId                  int64
}

func NewSelectDiff(title string) *SelectDiff {
	s := &SelectDiff{Scene: eui.NewScene(eui.NewAbsoluteLayout())}
	s.topBar = eui.NewTopBar(title, nil)
	s.Add(s.topBar)
	s.frame = eui.NewContainer(eui.NewVBoxLayout(1))
	lblTitle := eui.NewText("Настрой сложность")
	s.frame.Add(lblTitle)
	s.column, s.row, s.percent = 5, 5, 15

	gn := func(a, b int) []int {
		arr := make([]int, 0)
		for i := a; i < b; i++ {
			arr = append(arr, i)
		}
		return arr
	}

	s.cRow = eui.NewSpinBox(
		"Сколько рядов",
		gn(5, 51),
		0,
		func(i int) string { return fmt.Sprintf("%v", i) },
		func(a, b int) int { return a - b },
		func(a, b int) bool { return a == b },
		true,
		2,
	)
	s.cRow.SelectedValue.Connect(func(value int) {
		s.row = value
	})
	s.frame.Add(s.cRow)
	s.cCol = eui.NewSpinBox(
		"Сколько столбиков",
		gn(5, 51),
		0,
		func(i int) string { return fmt.Sprintf("%v", i) },
		func(a, b int) int { return a - b },
		func(a, b int) bool { return a == b },
		true,
		2,
	)
	s.cCol.SelectedValue.Connect(func(value int) {
		s.column = value
	})
	s.frame.Add(s.cCol)
	s.cPerc = eui.NewSpinBox(
		"Сколько % мин", gn(15, 31),
		0,
		func(i int) string { return fmt.Sprintf("%v", i) },
		func(a, b int) int { return a - b },
		func(a, b int) bool { return a == b },
		true,
		2,
	)
	s.cPerc.SelectedValue.Connect(func(value int) {
		s.percent = value
	})
	s.frame.Add(s.cPerc)

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
	log.Println("run game", s.row, s.column, mines, s.percent)
}

func (s *SelectDiff) SetRect(rect eui.Rect[int]) {
	w, h := rect.Size()
	s.Scene.SetRect(rect)
	hTop := int(float64(rect.GetLowestSize()) * 0.05)
	s.topBar.SetRect(eui.NewRect([]int{0, 0, w, hTop}))
	s.frame.SetRect(eui.NewRect([]int{0, hTop, w, h - hTop}))
}
