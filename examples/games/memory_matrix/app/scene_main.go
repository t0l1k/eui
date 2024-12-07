package app

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/colors"
	"github.com/t0l1k/eui/examples/games/memory_matrix/mem"
)

type SceneMain struct {
	eui.SceneBase
	topBar    *eui.TopBar
	board     *BoardMem
	lblStatus *eui.Text
}

func NewSceneMain() *SceneMain {
	s := &SceneMain{}
	s.topBar = eui.NewTopBar(title, nil)
	s.topBar.SetUseStopwatch()
	s.topBar.SetShowStoppwatch(true)
	s.topBar.SetTitleCoverArea(0.9)
	s.Add(s.topBar)
	s.lblStatus = eui.NewText("")
	s.Add(s.lblStatus)
	s.board = NewBoardMem(func(btn *eui.Button) {
		for i, v := range s.board.layout.GetContainer() {
			if v.(*eui.Button) == btn {
				switch s.board.Game().Stage() {
				case mem.Preparation:
					s.board.game.Move(i)
					s.board.showTimer.On()
					s.board.SetupShow()
				case mem.Show:
				case mem.Recollection:
					if s.board.game.Move(i) {
						v.(*eui.Button).Bg(colors.Aqua)
					} else {
						v.(*eui.Button).Bg(colors.Orange)
					}
				case mem.Conclusion:
				case mem.Restart:
					s.board.game.NextLevel()
					s.board.SetupPreparation()
				}
			}
		}
	})
	s.Add(s.board)
	s.board.varMsg.Attach(s.lblStatus)
	s.board.varColor.Attach(s.lblStatus)
	s.lblStatus.SetText(s.board.Game().String())
	s.Resize()
	return s
}

func (s *SceneMain) Resize() {
	w0, h0 := eui.GetUi().Size()
	hTop := int(float64(h0) * 0.05) // topbar height
	s.topBar.Resize([]int{0, 0, w0, hTop})
	s.board.Resize([]int{hTop, hTop * 2, w0 - hTop*2, h0 - hTop*4})
	s.lblStatus.Resize([]int{0, h0 - hTop, w0, hTop})
}
