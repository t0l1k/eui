package app

import (
	"fmt"

	"github.com/t0l1k/eui"
)

type SceneMain struct {
	eui.SceneBase
	topBar    *eui.TopBar
	board     *BoardSol15
	bottomBar *BottomBar
}

func NewSceneMain() *SceneMain {
	s := &SceneMain{}
	s.topBar = eui.NewTopBar(title, nil)
	s.topBar.SetUseStopwatch()
	s.topBar.SetShowStoppwatch(true)
	s.topBar.SetTitleCoverArea(0.9)
	s.Add(s.topBar)
	s.board = NewBoardSol15(func(btn *eui.Button) {
		for k, cells := range *s.board.game {
			for i, cell := range cells {
				card := cell.GetCard()
				if card == nil {
					continue
				}
				if card.StringShort() == btn.GetText() {
					count := 3 - i
					for count >= 0 {
						s.board.MakeMove(k)
						count--
					}
					s.bottomBar.updateMoveCount()
				}
			}
		}
	})
	s.board.Setup(true)
	s.Add(s.board)
	s.bottomBar = NewBottomBar(func(btn *eui.Button) {
		fmt.Println("pressed:", btn.GetText())
		switch btn.GetText() {
		case actNew:
			s.board.Setup(true)
			s.bottomBar.updateMoveCount()
			sb := eui.NewSnackBar("Новый рассклад!").Show(2000)
			s.Add(sb)
		case actReset:
			s.board.sw.Stop()
			s.board.Setup(false)
			s.bottomBar.updateMoveCount()
			sb := eui.NewSnackBar("Повторить собирать рассклад!").Show(1000)
			s.Add(sb)
		case actBackwardMove:
			if s.board.moveIdx > 0 {
				s.board.moveIdx--
			}
			s.board.game.SetDeck(s.board.historyOfMoves[s.board.moveIdx])
			fmt.Println("idx:", s.board.moveIdx)
		case actForwardMove:
			if s.board.moveIdx < len(s.board.historyOfMoves)-1 {
				s.board.moveIdx++
			}
			s.board.game.SetDeck(s.board.historyOfMoves[s.board.moveIdx])
			fmt.Println("idx:", s.board.moveIdx)
		}
	})
	s.bottomBar.Setup(s.board)
	s.Add(s.bottomBar)
	s.Resize()
	return s
}

func (s *SceneMain) Resize() {
	w0, h0 := eui.GetUi().Size()
	rect := eui.NewRect([]int{0, 0, w0, h0})
	Htop := int(float64(rect.GetLowestSize()) * 0.05)
	s.topBar.Resize([]int{0, 0, w0, Htop})
	s.board.Resize([]int{0, Htop, w0, h0 - Htop*4})
	s.bottomBar.Resize([]int{0, h0 - Htop, w0, Htop})
}
