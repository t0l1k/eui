package app

import (
	"image/color"

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
		if btn.IsPressed() {
			switch s.board.Game().Stage() {
			case mem.Preparation:
				s.board.game.Move(0)
				s.board.showTimer.On()
				s.board.SetupShow()
			case mem.Restart:
				s.board.game.NextLevel()
				s.board.SetupPreparation()
			}
		}
		for i, v := range s.board.layout.GetContainer() {
			switch vv := v.(type) {
			case *eui.Button:
				if vv == btn {
					switch s.board.Game().Stage() {
					case mem.Recollection:
						if s.board.game.Move(i) {
							v.(*eui.Button).Bg(colors.Aqua)
						} else {
							v.(*eui.Button).Bg(colors.Orange)
						}
					}
				}
			}
		}
	})
	s.Add(s.board)
	s.board.varMsg.Connect(func(data any) { s.lblStatus.SetText(data.(string)) })
	s.board.varColor.Connect(func(data any) {
		arr := data.([]color.Color)
		s.lblStatus.Bg(arr[0])
		s.lblStatus.Fg(arr[1])
	})
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
