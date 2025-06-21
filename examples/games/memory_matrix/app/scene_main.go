package app

import (
	"image/color"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/memory_matrix/mem"
	"golang.org/x/image/colornames"
)

type SceneMain struct {
	*eui.Scene
	topBar    *eui.TopBar
	board     *BoardMem
	lblStatus *eui.Text
}

func NewSceneMain() *SceneMain {
	s := &SceneMain{Scene: eui.NewScene(eui.NewAbsoluteLayout())}
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
		for i, v := range s.board.Childrens() {
			switch vv := v.(type) {
			case *eui.Button:
				if vv == btn {
					switch s.board.Game().Stage() {
					case mem.Recollection:
						if s.board.game.Move(i) {
							v.(*eui.Button).Bg(colornames.Aqua)
						} else {
							v.(*eui.Button).Bg(colornames.Orange)
						}
					}
				}
			}
		}
	})
	s.Add(s.board)
	s.board.varMsg.Connect(func(data string) { s.lblStatus.SetText(data) })
	s.board.varColor.Connect(func(arr []color.Color) {
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
	s.topBar.Resize(eui.NewRect([]int{0, 0, w0, hTop}))
	s.board.Resize(eui.NewRect([]int{hTop, hTop * 2, w0 - hTop*2, h0 - hTop*4}))
	s.lblStatus.Resize(eui.NewRect([]int{0, h0 - hTop, w0, hTop}))
}
