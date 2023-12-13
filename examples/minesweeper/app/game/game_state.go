package game

import (
	"github.com/t0l1k/eui"
)

const (
	GameStart = "start"
	GamePlay  = "play"
	GamePause = "pause"
	GameWin   = "win"
	GameOver  = "game over"
)

type gameState struct {
	eui.SubjectBase
}

func newGameState() *gameState {
	s := &gameState{}
	s.SetValue(GameStart)
	return s
}
