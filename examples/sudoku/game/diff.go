package game

import (
	"log"
	"math/rand"
)

type Difficult int

const (
	Easy Difficult = iota
	Normal
	Hard
	Extreme
)

func NewDiff(diff Difficult) Difficult { return diff }

func (d Difficult) Eq(other Difficult) bool { return d == other }
func (d Difficult) Size() int               { return int(Extreme) }
func (d Difficult) Percent(size int) (moves, percent int) {
	var (
		percentMin, percentMax int
	)
	switch d {
	case Easy:
		percentMin, percentMax = 20, 35
	case Normal:
		percentMin, percentMax = 35, 50
	case Hard:
		percentMin, percentMax = 50, 65
	case Extreme:
		percentMin, percentMax = 65, 80
	}
	percent = rand.Intn(percentMax-percentMin) + percentMin
	moves = percent * (size * size) / 100
	log.Printf("Сложность %v size:%v min:%v max:%v n:%v moves:%v\n", d, size, percentMin, percentMax, percent, moves)
	return moves, percent
}

func (d Difficult) String() (res string) {
	switch d {
	case Easy:
		res = "Легко"
	case Normal:
		res = "Нормально"
	case Hard:
		res = "Сложно"
	case Extreme:
		res = "Экстремально"
	}
	return res
}
