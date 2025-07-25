package scene

import (
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/lines/app"
)

type SceneGame struct {
	*eui.Scene
	topBar *eui.TopBar
	board  *Board
	dialog *Dialog
}

func NewSceneGame() *SceneGame {
	s := &SceneGame{Scene: eui.NewScene(eui.NewAbsoluteLayout())}
	s.topBar = eui.NewTopBar(app.GameTitle, func(b *eui.Button) {
		s.dialog.SetTitle("Выбор новой игры")
		s.dialog.SetHidden(false)
		s.board.SetHidden(true)
	})
	s.topBar.SetUseStopwatch()
	s.topBar.SetTitleCoverArea(0.5)
	s.Add(s.topBar)
	s.dialog = NewDialog("Запустить игру", func(dlg *eui.Button) {
		if dlg.GetText() == app.BNew {
			s.board.NewGame(s.dialog.diff)
		}
		s.dialog.SetHidden(true)
		s.board.SetHidden(false)
	})
	s.Add(s.dialog)
	s.board = NewBoard(s.dialog.diff)
	s.Add(s.board)
	return s
}

func (s *SceneGame) Entered() {
	s.Resize()
	s.dialog.comboSelGameDiff.SetValue(app.FieldSizeNormal)
	s.dialog.diff = s.dialog.comboSelGameDiff.Value().(int)
}

func (s *SceneGame) Resize() {
	w, h := eui.GetUi().Size()
	rect := eui.NewRect([]int{0, 0, w, h})
	hT := int(float64(rect.GetLowestSize()) * 0.1)
	s.topBar.Resize(eui.NewRect([]int{0, 0, w, hT}))
	s.board.Resize(eui.NewRect([]int{hT / 2, hT + hT/2, w - hT, h - hT*2}))
	w0, h0 := w, h
	x0, y0 := 0, 0
	w1 := w0 / 2
	h1 := h0 / 2
	x1 := x0 + (w0-w1)/2
	y1 := y0 + (h0-h1)/2
	s.dialog.Resize(eui.NewRect([]int{x1, y1, w1, h1}))
}
