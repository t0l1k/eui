package mem

import (
	"strconv"
	"time"
)

type GameData struct {
	tmBeginShow, tmBeginResolve, tmEndResolve time.Time
	score, level                              int
}

func NewGameData() *GameData           { return &GameData{} }
func (d *GameData) GetData() *GameData { return d }
func (d *GameData) SetBeginShow()      { d.tmBeginShow = time.Now() }
func (d *GameData) SetBeginResolve()   { d.tmBeginResolve = time.Now() }
func (d *GameData) SetEndResolve()     { d.tmEndResolve = time.Now() }
func (d *GameData) SetScore(count int) {
	d.level = count
	dur1 := d.tmEndResolve.Sub(d.tmBeginResolve)
	if int(dur1.Seconds()) < count {
		d.score += count * 10 * 2
	} else {
		d.score += count * 10
	}
}
func (d *GameData) String() string { return strconv.Itoa(d.score) }
