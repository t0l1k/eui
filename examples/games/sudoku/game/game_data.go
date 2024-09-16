package game

import (
	"fmt"
	"time"
)

type GameData struct {
	count int
	score time.Duration
	tm    time.Time
}

func NewGameData(count int) *GameData {
	return &GameData{count: count + 1, tm: time.Now()}
}

func (g *GameData) Score() time.Duration         { return g.score }
func (g *GameData) SetScore(value time.Duration) { g.score = value }
func (g *GameData) Eq(other GameData) bool {
	return g.count == other.count && g.score == other.score && g.tm == other.tm
}

func (g GameData) String() (result string) {
	return fmt.Sprintf("#%v %v %v", g.count, g.score.Round(time.Millisecond).String(), g.tm.Format("15:04:05"))
}

type GameByTime []GameData

func (ds GameByTime) Len() int           { return len(ds) }
func (ds GameByTime) Less(i, j int) bool { return ds[i].tm.After(ds[j].tm) }
func (ds GameByTime) Swap(i, j int)      { ds[i], ds[j] = ds[j], ds[i] }
