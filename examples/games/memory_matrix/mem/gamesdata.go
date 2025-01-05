package mem

import "fmt"

type GamesData []GameData

func NewGamesData() *GamesData { return &GamesData{} }

func (g GamesData) Add(dt GameData) *GamesData { g = append(g, dt); return &g }
func (g GamesData) IsEmpty() bool              { return len(g) == 0 }
func (g GamesData) Size() int                  { return len(g) }

func (g GamesData) Max() (max int) {
	for _, v := range g {
		if v.level > max {
			max = v.level
		}
	}
	return max
}

func (g GamesData) Average() (value float64) {
	if g.IsEmpty() {
		return float64(g.Size())
	}
	for _, v := range g {
		value += float64(v.level)
	}
	return value / float64(g.Size())
}

func (g GamesData) Score() (score int) {
	for _, v := range g {
		score += v.score
	}
	return score
}

func (g GamesData) Levels() (result []int) {
	for _, v := range g {
		result = append(result, v.level)
	}
	return result
}

func (g GamesData) String() string {
	return fmt.Sprintf("games:%v max:%v average:%.02v score:%v", g.Size(), g.Max(), g.Average(), g.Score())
}
