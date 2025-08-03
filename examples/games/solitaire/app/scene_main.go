package app

import (
	"fmt"
	"time"

	"github.com/t0l1k/eui"
)

type SceneMain struct {
	*eui.Scene
	topBar    *eui.TopBar
	boards    []Sols
	boardIdx  int
	current   Sols
	fn        func(*eui.Button)
	bottomBar *BottomBar
}

func NewSceneMain() *SceneMain {
	s := &SceneMain{Scene: eui.NewScene(eui.NewAbsoluteLayout())}
	s.topBar = eui.NewTopBar(title, nil)
	s.topBar.SetUseStopwatch()
	s.topBar.SetShowStoppwatch(true)
	s.topBar.SetTitleCoverArea(0.9)
	s.Add(s.topBar)
	s.fn = s.gameLogic
	var board Sols

	board = NewBoardSol15(s.fn)
	s.Add(board)
	board.Setup(true)
	s.boards = append(s.boards, board)

	// board = nil

	// board = NewBoardFreecell(s.fn)
	// s.Add(board)
	// board.Setup(true)
	// s.boards = append(s.boards, board)

	s.current = s.boards[s.boardIdx]
	// s.current.SetHidden(true)
	s.bottomBar = NewBottomBar(func(btn *eui.Button) {
		switch btn.GetText() {
		case actNextSol:
			s.current.SetHidden(true)
			s.current = nil
			s.boardIdx++
			if s.boardIdx >= len(s.boards) {
				s.boardIdx = 0
			}
			s.current = s.boards[s.boardIdx]
			s.current.Setup(true)
			s.current.SetHidden(false)
			s.Add(s.current)
			s.SetRect(s.Rect())
			fmt.Println("pressed:", btn.GetText(), s.boardIdx)
		case actNew:
			s.current.Setup(true)
			s.bottomBar.UpdateMoveCount()
			eui.NewSnackBar("Новый рассклад!").Show(2 * time.Second)
		case actReset:
			s.current.Stopwatch().Stop()
			s.current.Setup(false)
			s.bottomBar.UpdateMoveCount()
			eui.NewSnackBar("Повторить собирать рассклад!").Show(1 * time.Second)
		case actBackwardMove:
			if s.current.GetMoveNr() > 0 {
				s.current.SetMoveNr(s.current.GetMoveNr() - 1)
				s.bottomBar.UpdateMoveCount()
			}
			s.current.Game().SetDeck(s.current.GetHistory()[s.current.GetMoveNr()])
		case actForwardMove:
			if s.current.GetMoveNr() < len(s.current.GetHistory())-1 {
				s.current.SetMoveNr(s.current.GetMoveNr() + 1)
				s.bottomBar.UpdateMoveCount()
			}
			s.current.Game().SetDeck(s.current.GetHistory()[s.current.GetMoveNr()])
		}
	})
	s.bottomBar.Setup(s.current)
	s.Add(s.bottomBar)
	return s
}

func (s *SceneMain) gameLogic(btn *eui.Button) {
	for _, card := range s.current.Game().GetDeck() {
		if card == nil {
			continue
		}
		if card.StringShort() == btn.GetText() {
			column, idx := s.current.Game().Index(card)
			fmt.Println("sc mv:", column, idx, btn.GetText())
			for idx >= 0 {
				s.current.MakeMove(column)
				idx--
			}
			if s.bottomBar.UpdateMoveCount() {
				eui.NewSnackBar("Нет ходов").Show(1 * time.Second)
			}
		}
	}
}

func (s *SceneMain) SetRect(rect eui.Rect[int]) {
	w0, h0 := rect.Size()
	s.Scene.SetRect(rect)
	Htop := int(float64(rect.GetLowestSize()) * 0.05)
	s.topBar.SetRect(eui.NewRect([]int{0, 0, w0, Htop}))
	s.current.SetRect(eui.NewRect([]int{0, Htop, w0, h0 - Htop*4}))
	s.bottomBar.SetRect(eui.NewRect([]int{0, h0 - Htop, w0, Htop}))
}
