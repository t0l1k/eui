package app

import (
	"strconv"

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
				size, _, diff := v.GetData()
				s.topBar.SetTitle("Sudoku " + diff.String())
				s.dialogSelect.Visible(false)
				s.board.Setup(size, diff)
				s.board.Visible(true)
				s.bottomBar.Setup(size)
				s.bottomBar.Visible(true)
				s.bottomBar.UpdateNrs(s.board.field.ValuesCount())
			}
		}
	})
	s.Add(s.dialogSelect)
	s.dialogSelect.Visible(true)
	s.board = NewBoard(func(btn *eui.Button) {
		for i := range s.board.layoutCells.GetContainer() {
			icon := s.board.layoutCells.GetContainer()[i].(*CellIcon)
			if icon.btn == btn {
				x, y := s.board.field.Pos(i)
				if s.bottomBar.IsActDel() {
					s.board.field.ResetCell(x, y)
				} else if s.bottomBar.IsActNotes() {
				} else if s.bottomBar.IsActUndo() {
				} else {
					s.board.field.Add(x, y, s.board.GetHighlightValue())
					s.board.Highlight(strconv.Itoa(s.board.GetHighlightValue()))
					s.bottomBar.UpdateNrs(s.board.field.ValuesCount())
				}
			}
		}

	})
	s.Add(s.board)
	s.bottomBar = NewBottomBar(func(btn *eui.Button) {
		if s.bottomBar.SetAct(btn) {
			s.board.Highlight(btn.GetText())
		} else {
			s.board.Highlight("0")
		}
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
	h1 := h - hT*2
	s.board.Resize([]int{0, hT, w, h1})
	s.bottomBar.Resize([]int{(w - (h1)) / 2, h - hT, h1, hT})
}
