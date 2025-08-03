package app

import (
	"log"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/sudoku/game"
	"golang.org/x/image/colornames"
)

type SceneSudoku struct {
	*eui.Scene
	topBar       *eui.TopBar
	dialogSelect *DialogSelect
	board        *Board
	bottomBar    *BottomBar
}

func NewSceneSudoku() *SceneSudoku {
	gamesData := game.NewGamesData()
	s := &SceneSudoku{Scene: eui.NewScene(eui.NewAbsoluteLayout())}
	s.topBar = eui.NewTopBar(Title, func(b *eui.Button) {
		s.dialogSelect.SetHidden(false)
		s.board.SetHidden(true)
		s.bottomBar.SetHidden(true)
		s.topBar.SetTitle(Title)
		s.topBar.SetShowTitle(false)
		s.topBar.SetShowStoppwatch(false)
	})
	s.topBar.SetUseStopwatch()
	s.topBar.SetTitleCoverArea(0.85)
	s.Add(s.topBar)
	s.dialogSelect = NewDialogSelect(gamesData, func(b *eui.Button) {
		for _, v := range s.dialogSelect.btnsDiff {
			if v.btn.IsPressed() {
				dim, diff := v.GetData()
				s.topBar.SetTitle("Sudoku " + dim.String() + diff.String())
				s.dialogSelect.SetHidden(true)
				s.board.Setup(dim, diff)
				s.topBar.SetShowTitle(false)
				s.board.SetHidden(false)
				s.bottomBar.Setup(s.board)
				s.bottomBar.SetHidden(false)
				s.bottomBar.UpdateNrs(s.board.game.ValuesCount())
				s.bottomBar.ShowNotes(s.board.IsShowNotes())
				s.topBar.SetShowStoppwatch(false)
			}
		}
	})
	s.Add(s.dialogSelect)
	s.dialogSelect.SetHidden(false)
	s.board = NewBoard(func(btn *eui.Button) {
		for i := range s.board.layoutCells.Childrens() {
			icon, ok := s.board.layoutCells.Childrens()[i].(*CellIcon)
			if !ok {
				continue
			}
			if icon.btn == btn {
				x, y := s.board.game.Pos(i)
				if s.bottomBar.IsActDel() {
					s.board.game.ResetCell(x, y)
					s.bottomBar.UpdateNrs(s.board.game.ValuesCount())
					log.Println("Set Act Del", x, y)
				} else {
					if !s.board.isWin {
						s.board.Move(x, y)
						s.bottomBar.UpdateUndoBtn(s.board.MoveCount())
						s.bottomBar.UpdateNrs(s.board.game.ValuesCount())
					}
					if s.board.isWin {
						gamesData.AddGameResult(s.board.dim, s.board.diff, s.board.sw.Duration())
						margin := s.dialogSelect.margin
						s.dialogSelect.history.Reset()
						s.dialogSelect.history.SetupListViewText(gamesData.GamesPlayed(), margin, 1, colornames.Aqua, colornames.Black)
						for _, v := range s.dialogSelect.btnsDiff {
							if v.diff.Eq(s.board.diff) && v.dim.Eq(s.board.dim) {
								value := gamesData.GetLastBest(s.board.dim, s.board.diff)
								v.SetScore(value)
							}
						}
						s.board.SetHidden(true)
						s.bottomBar.SetHidden(true)
						s.dialogSelect.SetHidden(false)
						s.topBar.SetTitle(Title)
						s.topBar.SetShowTitle(true)
						s.topBar.SetShowStoppwatch(true)
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
			if s.bottomBar.IsActAccept() {
				s.board.game.MarkReadOnly()
				s.board.sw.Start()
				log.Println("Set Act Accept in hand game", s.board.GetDiffStr(), s.bottomBar.varDiff.Value())
			}
		}
	})
	s.Add(s.bottomBar)
	return s
}

func (s *SceneSudoku) SetRect(rect eui.Rect[int]) {
	w0, h0 := rect.Size()
	s.Scene.SetRect(rect)
	hT := int(float64(rect.GetLowestSize()) * 0.1)
	s.topBar.SetRect(eui.NewRect([]int{0, 0, w0, hT}))
	s.dialogSelect.SetRect(eui.NewRect([]int{hT / 2, hT + hT/2, w0 - hT, h0 - hT*2}))
	w1 := (w0 - hT) / 3
	s.board.SetRect(eui.NewRect([]int{hT, 0, w1 * 2, h0}))
	s.bottomBar.SetRect(eui.NewRect([]int{hT + w1*2, 0, w1, h0}))
}
