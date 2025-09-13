package app

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/sudoku/game"
	"golang.org/x/image/colornames"
)

type SceneSudoku struct {
	*eui.Scene
	topBar       *eui.Topbar
	dialogSelect *DialogSelect
	board        *Board
	bottomBar    *BottomBar
}

func NewSceneSudoku() *SceneSudoku {
	gamesData := game.NewGamesData()
	s := &SceneSudoku{Scene: eui.NewScene(eui.NewLayoutVerticalPercent([]int{5, 95}, 5))}
	s.topBar = eui.NewTopBar(Title, func(b *eui.Button) {
		s.board.Hide()
		s.bottomBar.Hide()
		s.topBar.SetTitle(Title)
		s.topBar.SetShowTitle(false)
		s.topBar.SetShowStoppwatch(false)
		s.dialogSelect.Show()
	}).SetUseStopwatch()
	s.dialogSelect = NewDialogSelect(gamesData, func(b *eui.Button) {
		for _, v := range s.dialogSelect.btnsDiff {
			newVar, newVar1 := v.GetData()
			log.Println(b.Rect(), v.Rect(), b.Text(), newVar, newVar1, v.State(), b.State(), v.Rect().InRect(b.Rect().Pos()))
			x, y := ebiten.CursorPosition()
			pos := eui.NewPoint(x, y)
			if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) && pos.In(v.Rect()) {
				dim, diff := v.GetData()
				s.topBar.SetTitle("Sudoku " + dim.String() + diff.String())
				s.dialogSelect.Hide()
				s.board.Setup(dim, diff)
				s.topBar.SetShowTitle(false)
				s.bottomBar.Setup(s.board)
				s.bottomBar.UpdateNrs(s.board.game.ValuesCount())
				s.bottomBar.ShowNotes(s.board.IsShowNotes())
				s.topBar.SetShowStoppwatch(false)
				s.board.Show()
				s.bottomBar.Show()
			}
		}
	})
	s.dialogSelect.Show()
	s.board = NewBoard(func(btn *eui.Button) {
		for i := range s.board.layoutCells.Children() {
			icon, ok := s.board.layoutCells.Children()[i].(*CellIcon)
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
						s.dialogSelect.history.Reset()
						s.dialogSelect.history.SetupListViewText(gamesData.GamesPlayed(), 30, 1, colornames.Aqua, colornames.Black)
						for _, v := range s.dialogSelect.btnsDiff {
							if v.diff.Eq(s.board.diff) && v.dim.Eq(s.board.dim) {
								value := gamesData.GetLastBest(s.board.dim, s.board.diff)
								v.SetScore(value)
							}
						}
						s.board.Hide()
						s.bottomBar.Hide()
						s.dialogSelect.Show()
						s.topBar.SetTitle(Title)
						s.topBar.SetShowTitle(true)
						s.topBar.SetShowStoppwatch(true)
					}
				}
			}
		}
	})
	s.bottomBar = NewBottomBar(func(btn *eui.Button) {
		if s.bottomBar.SetAct(btn) {
			s.board.Highlight(btn.Text())
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

	contGame := eui.NewContainer(eui.NewLayoutHorizontalPercent([]int{70, 30}, 1))
	contGame.Add(s.board)
	contGame.Add(s.bottomBar)

	contBoard := eui.NewContainer(eui.NewStackLayout(5))
	contBoard.Add(s.dialogSelect)
	contBoard.Add(contGame)

	s.Add(s.topBar)
	s.Add(contBoard)
	return s
}
