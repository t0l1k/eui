package app

import (
	"fmt"

	"github.com/t0l1k/eui"
)

type SceneSudoku struct {
	eui.SceneBase
	topBar       *eui.TopBar
	dialogSelect *DialogSelect
	board        *Board
}

func NewSceneSudoku() *SceneSudoku {
	s := &SceneSudoku{}
	s.topBar = eui.NewTopBar("Sudoku", func(b *eui.Button) {
		s.dialogSelect.Visible(true)
	})
	s.topBar.SetShowStopwatch()
	s.topBar.SetTitleCoverArea(0.85)
	s.Add(s.topBar)
	s.dialogSelect = NewDialogSelect(func(b *eui.Button) {
		for _, v := range s.dialogSelect.btnsDiff {
			if v.btn.IsPressed() {
				a, b, c := v.GetData()
				fmt.Println("pressed:", v.title.GetText(), a, b, c)
				s.topBar.SetTitle("Sudoku " + c.String())
				s.dialogSelect.Visible(false)
				s.board.Setup(a)
			}
		}
	})
	s.Add(s.dialogSelect)
	s.board = NewBoard()
	s.Add(s.board)
	s.Resize()
	return s
}

func (s *SceneSudoku) Resize() {
	w, h := eui.GetUi().Size()
	rect := eui.NewRect([]int{0, 0, w, h})
	hT := int(float64(rect.GetLowestSize()) * 0.1)
	s.topBar.Resize([]int{0, 0, w, hT})
	s.dialogSelect.Resize([]int{hT / 2, hT + hT/2, w - hT, h - hT*2})
	s.board.Resize([]int{0, hT, w, h - hT})
}
