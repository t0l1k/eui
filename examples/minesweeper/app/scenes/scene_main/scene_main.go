package scene_main

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/minesweeper/app/scenes"
	"github.com/t0l1k/eui/examples/minesweeper/app/scenes/scene_game"
)

type SceneSelectGame struct {
	eui.SceneDefault
	topBar *scenes.TopBar
	frame  *eui.BoxLayout
	sDiff  []string
}

func NewSceneSelectGame() *SceneSelectGame {
	s := &SceneSelectGame{}

	s.topBar = scenes.NewTopBar("Игра Сапёр")
	s.Add(s.topBar)

	bg := eui.Green
	fg := eui.Yellow

	s.frame = eui.NewVLayout()
	lblTitle := eui.NewText("Выбери сложность", bg, fg)
	s.frame.Add(lblTitle)
	s.sDiff = []string{"Игра новичку", "Игра Легко", "Игра Средне", "Игра Сложно", "Настроить сложность"}
	for _, value := range s.sDiff {
		bg := eui.Gray
		fg := eui.Maroon
		button := eui.NewButton(value, bg, fg, s.selectGameLogic)
		s.frame.Add(button)
	}
	s.Add(s.frame)
	s.Resize()
	return s
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
