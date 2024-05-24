package game

import (
	"time"
)

type GamesData map[Dim]map[Difficult][]GameData

func NewGamesData() *GamesData {
	gDim := make(GamesData)
	dims := func() (result []Dim) {
		for i := 2; i <= 3; i++ {
			result = append(result, NewDim(i, i))
		}
		return result
	}()
	diffs := func() (result []Difficult) {
		for i := 0; i < 4; i++ {
			result = append(result, NewDiff(Difficult(i)))
		}
		return result
	}()
	for _, dim := range dims {
		gDim[dim] = make(map[Difficult][]GameData)
		for _, v := range diffs {
			gDim[dim][v] = make([]GameData, 0)
		}
	}
	return &gDim
}

func (g GamesData) AddGameResult(dim Dim, diff Difficult, result time.Duration) {
	gm := NewGameData(len(g[dim][diff]))
	gm.SetScore(result)
	g[dim][diff] = append(g[dim][diff], *gm)
}

func (g GamesData) GetLastBest(dim Dim, diff Difficult) (result string) {
	var best time.Duration
	for _, v := range g[dim][diff] {
		if v.Score() < best || best == 0 {
			best = v.Score()
		}
	}
	value := "0"
	if len(g[dim][diff]) > 0 {
		value = g[dim][diff][len(g[dim][diff])-1].String()
	}
	result += "Best:" + best.Round(time.Millisecond).String() + " Last:" + value
	return result
}

func (g GamesData) String() (result string) {
	for dim, diffs := range g {
		for diff, games := range diffs {
			result += dim.String()
			result += diff.String()
			for _, v := range games {
				result += v.String()
			}
			result += "\n"
		}
	}
	return result
}
