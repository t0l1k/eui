package scene_main

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/minesweeper/app/scenes/scene_game"
)

type SceneSelectGame struct {
	*eui.Scene
	topBar *eui.TopBar
	frame  *eui.Container
	sDiff  map[string][]int
}

func NewSceneSelectGame() *SceneSelectGame {
	s := &SceneSelectGame{Scene: eui.NewScene(eui.NewAbsoluteLayout())}
	s.topBar = eui.NewTopBar("Игра Сапёр", nil)
	s.topBar.SetUseStopwatch()
	s.Add(s.topBar)
	s.frame = eui.NewContainer(eui.NewVBoxLayout(1))
	lblTitle := eui.NewText("Выбери сложность")
	s.frame.Add(lblTitle)
	s.sDiff = make(map[string][]int)
	keys := []string{"Игра новичку", "Игра Легко", "Игра Средне", "Игра Сложно", "Настроить сложность"}
	v := [][]int{{5, 5, 5}, {8, 8, 10}, {16, 16, 40}, {40, 16, 99}, {5, 5, 5}}
	for i, k := range keys {
		s.sDiff[k] = v[i]
		button := eui.NewButton(k, s.selectGameLogic)
		s.frame.Add(button)
	}
	s.Add(s.frame)
	return s
}

func (s *SceneSelectGame) selectGameLogic(b *eui.Button) {
	if b.GetText() == "Настроить сложность" {
		game := NewSelectDiff(b.GetText())
		eui.GetUi().Push(game)
	} else {
		r, c, m := s.sDiff[b.GetText()][0], s.sDiff[b.GetText()][1], s.sDiff[b.GetText()][2]
		game := scene_game.NewSceneGame(b.GetText(), r, c, m)
		eui.GetUi().Push(game)
	}
}

func (s *SceneSelectGame) Resize() {
	w, h := eui.GetUi().Size()
	rect := eui.NewRect([]int{0, 0, w, h})
	s.SetRect(rect)
	hT := int(float64(rect.GetLowestSize()) * 0.05)
	s.topBar.Resize(eui.NewRect([]int{0, 0, w, hT}))
	s.frame.Resize(eui.NewRect([]int{hT / 2, hT + hT/2, w - hT, h - hT*2}))
}
