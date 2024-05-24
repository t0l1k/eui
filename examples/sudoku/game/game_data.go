package game

import (
	"fmt"
	"time"
)

type GameData struct {
	count int
	score time.Duration
}

func NewGameData(count int) *GameData {
	return &GameData{count: count + 1}
}

func (g *GameData) Score() time.Duration         { return g.score }
func (g *GameData) SetScore(value time.Duration) { g.score = value }

func (g GameData) String() (result string) {
	return fmt.Sprintf("#%v %v", g.count, g.score.Round(time.Millisecond).String())
}
