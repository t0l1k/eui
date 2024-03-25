package game

import (
	"fmt"
	"math/rand"
)

type Difficult int

const (
	Easy Difficult = iota
	Normal
	Hard
)

func (d Difficult) Size() int { return int(Hard) }
func (d Difficult) Percent(dim int) (moves int) {
	var (
		percentMin, percentMax, size int
	)
	size = dim * dim
	switch d {
	case Easy:
		percentMin, percentMax = 20, 35
	case Normal:
		percentMin, percentMax = 35, 50
	case Hard:
		percentMin, percentMax = 50, 75
	}
	n := rand.Intn(percentMax-percentMin) + percentMin
	moves = n * (size * size) / 100
	fmt.Printf("Сложность %v size:%v min:%v max:%v n:%v moves:%v\n", d, size, percentMin, percentMax, n, moves)
	return moves
}
func (d Difficult) String() (res string) {
	switch d {
	case Easy:
		res = "Легко"
	case Normal:
		res = "Нормально"
	case Hard:
		res = "Сложно"
	}
	return res
}
