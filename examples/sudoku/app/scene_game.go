package app

import (
	"log"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/sudoku/game"
)

type SceneSudoku struct {
	eui.SceneBase
	topBar       *eui.TopBar
	dialogSelect *DialogSelect
	board        *Board
	bottomBar    *BottomBar
}

func NewSceneSudoku() *SceneSudoku {
	gamesData := game.NewGamesData()
	s := &SceneSudoku{}
	s.topBar = eui.NewTopBar("Sudoku", func(b *eui.Button) {
		s.dialogSelect.Visible(true)
		s.board.Visible(false)
		s.bottomBar.Visible(false)
	})
	s.topBar.SetShowStopwatch()
	s.topBar.SetTitleCoverArea(0.85)
	s.Add(s.topBar)
	s.dialogSelect = NewDialogSelect(gamesData, func(b *eui.Button) {
		for _, v := range s.dialogSelect.btnsDiff {
			if v.btn.IsPressed() {
				dim, diff := v.GetData()
				s.topBar.SetTitle("Sudoku " + dim.String() + diff.String())
				s.dialogSelect.Visible(false)
				s.board.Setup(dim, diff)
				s.board.Visible(true)
				s.bottomBar.Setup(dim)
				s.bottomBar.Visible(true)
				s.bottomBar.UpdateNrs(s.board.game.ValuesCount())
				s.bottomBar.ShowNotes(s.board.IsShowNotes())
			}
		}
	})
	s.Add(s.dialogSelect)
	s.dialogSelect.Visible(true)
	s.board = NewBoard(func(btn *eui.Button) {
		for i := range s.board.layoutCells.GetContainer() {
			icon := s.board.layoutCells.GetContainer()[i].(*CellIcon)
			if icon.btn == btn {
				x, y := s.board.game.Pos(i)
				if s.bottomBar.IsActDel() {
					s.board.game.ResetCell(x, y)
					log.Println("Set Act Del", x, y)
				} else {
					if !s.board.isWin {
						s.board.Move(x, y)
						s.bottomBar.UpdateNrs(s.board.game.ValuesCount())
						s.bottomBar.UpdateUndoBtn(s.board.MoveCount())
					}
					if s.board.isWin {
						gamesData.AddGameResult(s.board.dim, s.board.diff, s.board.sw.Duration())
						for _, v := range s.dialogSelect.btnsDiff {
							if v.diff.Eq(s.board.diff) && v.dim.Eq(s.board.dim) {
								value := gamesData.GetLastBest(s.board.dim, s.board.diff)
								v.SetScore(value)
							}
						}
						s.board.Visible(false)
						s.bottomBar.Visible(false)
						s.dialogSelect.Visible(true)
					}
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
			if s.bottomBar.IsActNotes() {
				s.board.ShowNotes(true)
				log.Println("Set Act Notes", s.board.IsShowNotes())
			} else {
				s.board.ShowNotes(false)
				log.Println("Set Act Notes", s.board.IsShowNotes())
			}
			if s.bottomBar.IsActUndo() {
				s.board.Undo()
				s.bottomBar.UpdateNrs(s.board.game.ValuesCount())
				s.bottomBar.UpdateUndoBtn(s.board.MoveCount())
				log.Println("Set Act Undo")
			}
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
	s.board.Resize([]int{(w - (h1)) / 2, hT, h1, h1})
	s.bottomBar.Resize([]int{(w - (h1)) / 2, h - hT, h1, hT})
}
