package mem

import "fmt"

type GamesData []GameData

func NewGamesData() *GamesData { return &GamesData{} }

func (g GamesData) Add(dt GameData) *GamesData { g = append(g, dt); return &g }
func (g GamesData) Max() (max int) {
	for _, v := range g {
		if v.level > max {
			max = v.level
		}
	}
	return max
}
func (g GamesData) Score() (score int) {
	for _, v := range g {
		score += v.score
	}
	return score
}
func (g GamesData) Size() (max int) { return len(g) }

func (g GamesData) Levels() (result []int) {
	for _, v := range g {
		result = append(result, v.level)
	}
	return result
}

func (g GamesData) String() string {
	return fmt.Sprintf("games:%v max:%v score:%v", g.Size(), g.Max(), g.Score())
}
