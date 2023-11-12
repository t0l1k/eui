package scene_main

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/minesweeper/app/scenes/scene_game"
)

type SceneSelectGame struct {
	eui.SceneBase
	topBar *eui.TopBar
	frame  *eui.BoxLayout
	sDiff  []string
}

func NewSceneSelectGame() *SceneSelectGame {
	s := &SceneSelectGame{}

	s.topBar = eui.NewTopBar("Игра Сапёр")
	s.Add(s.topBar)
	s.frame = eui.NewVLayout()
	lblTitle := eui.NewText("Выбери сложность")
	s.frame.Add(lblTitle)
	s.sDiff = []string{"Игра новичку", "Игра Легко", "Игра Средне", "Игра Сложно", "Настроить сложность"}
	for _, value := range s.sDiff {
		button := eui.NewButton(value, s.selectGameLogic)
		s.frame.Add(button)
	}
	s.Add(s.frame)
	return s
}

func (s *SceneSelectGame) Entered() {
	s.Resize()
}

func (s *SceneSelectGame) selectGameLogic(b *eui.Button) {
	switch b.GetText() {
	case s.sDiff[0]:
		game := scene_game.NewSceneGame(b.GetText(), 5, 5, 5)
		eui.GetUi().Push(game)
	case s.sDiff[1]:
		game := scene_game.NewSceneGame(b.GetText(), 8, 8, 10)
		eui.GetUi().Push(game)
	case s.sDiff[2]:
		game := scene_game.NewSceneGame(b.GetText(), 16, 16, 40)
		eui.GetUi().Push(game)
	case s.sDiff[3]:
		game := scene_game.NewSceneGame(b.GetText(), 40, 16, 99)
		eui.GetUi().Push(game)
	case s.sDiff[4]:
		game := NewSelectDiff(b.GetText())
		eui.GetUi().Push(game)
	}

}

func (s *SceneSelectGame) Resize() {
	w, h := eui.GetUi().Size()
	rect := eui.NewRect([]int{0, 0, w, h})
	hTop := int(float64(rect.GetLowestSize()) * 0.05)
	s.topBar.Resize([]int{0, 0, w, hTop})
	s.frame.Resize([]int{0, hTop, w, h - hTop})
}
