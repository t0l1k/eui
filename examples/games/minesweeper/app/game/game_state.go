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

// Умею оповестить подписчиков о смене состояния игры
type gameState struct {
	*eui.Signal
}

func newGameState() *gameState {
	s := &gameState{Signal: eui.NewSignal()}
	s.Emit(GameStart)
	return s
}
