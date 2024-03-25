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
	bottomBar    *BottomBar
}

func NewSceneSudoku() *SceneSudoku {
	s := &SceneSudoku{}
	s.topBar = eui.NewTopBar("Sudoku", func(b *eui.Button) {
		s.dialogSelect.Visible(true)
		s.board.Visible(false)
		s.bottomBar.Visible(false)
	})
	s.topBar.SetShowStopwatch()
	s.topBar.SetTitleCoverArea(0.85)
	s.Add(s.topBar)
	s.dialogSelect = NewDialogSelect(func(b *eui.Button) {
		for _, v := range s.dialogSelect.btnsDiff {
			if v.btn.IsPressed() {
				size, score, diff := v.GetData()
				fmt.Println("pressed selected:", v.title.GetText(), size, score, diff)
				s.topBar.SetTitle("Sudoku " + diff.String())
				s.dialogSelect.Visible(false)
				s.board.Setup(size, diff)
				s.board.Visible(true)
				s.bottomBar.Setup(size)
				s.bottomBar.Visible(true)
			}
		}
	})
	s.Add(s.dialogSelect)
	s.dialogSelect.Visible(true)
	s.board = NewBoard(func(btn *eui.Button) {
		for i := range s.board.layout.GetContainer() {
			icon := s.board.layout.GetContainer()[i].(*CellIcon)
			if icon.btn == btn {
				cell := s.board.field.GetCells()[i]
				fmt.Println("pressed", cell.Value(), cell)
			}
		}

	})
	s.Add(s.board)
	s.bottomBar = NewBottomBar(func(btn *eui.Button) {
		fmt.Println("pressed", btn.GetText())
	})
	s.Add(s.bottomBar)
	s.Resize()
	return s
}

func (s *SceneSudoku) Resize() {
	w, h := eui.GetUi().Size()
	rect := eui.NewRect([]int{0, 0, w, h})
	hT := int(float64(rect.GetLowestSize()) * 0.1)
	s.topBar.Resize([]int{0, 0, w, hT})
	s.dialogSelect.Resize([]int{hT / 2, hT + hT/2, w - hT, h - hT*2})
	s.board.Resize([]int{0, hT, w, h - hT*3})
	s.bottomBar.Resize([]int{0, h - hT*2, w, hT * 2})
}
