package mem

import (
	"strconv"
	"time"
)

type GameData struct {
	tmBeginShow, tmBeginResolve, tmEndResolve time.Time
	score, level                              float64
}

func NewGameData(level float64) *GameData { return &GameData{level: level} }
func (d *GameData) GetData() *GameData    { return d }
func (d *GameData) SetBeginShow()         { d.tmBeginShow = time.Now() }
func (d *GameData) SetBeginResolve()      { d.tmBeginResolve = time.Now() }
func (d *GameData) SetEndResolve()        { d.tmEndResolve = time.Now() }
func (d *GameData) SetScore(count float64) {
	dur1 := d.tmEndResolve.Sub(d.tmBeginResolve)
	if dur1.Seconds() < count {
		d.score += count * 10 * 2
	} else {
		d.score += count * 10
	}
}
func (d *GameData) String() string { return strconv.FormatFloat(d.score, 'f', 1, 64) }
