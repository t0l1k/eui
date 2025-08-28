package mem

import (
	"strconv"
)

type Game struct {
	field         *field
	level         Level
	dim           Dim
	stage         GameStage
	moveCount     int
	Win, GameOver bool
	gameData      *GameData
}

func NewGame(level Level) *Game {
	g := &Game{}
	g.level = level
	g.dim = GetDimForLevel(level)
	g.field = newField().Shuffle(int(level), g.dim)
	g.stage = NewGameStage()
	g.moveCount = 0
	g.gameData = NewGameData(int(level))
	return g
}

func (g *Game) NextLevel() {
	if g.Win {
		g.level++
	}
	if g.GameOver {
		if g.level > 1 {
			g.level--
		}
	}
	g.Win = false
	g.GameOver = false
	g.stage = NewGameStage()
	g.field = nil
	g.moveCount = 0
	g.dim = GetDimForLevel(g.level)
	g.field = newField().Shuffle(int(g.level), g.dim)
	g.gameData = NewGameData(int(g.level))
}

func (g *Game) Move(idx int) bool {
	switch g.stage {
	case Preparation:
		g.SetNextStage()
	case Show:
	case Recollection:
		if g.Cell(idx).IsEmpty() {
			g.Cell(idx).SetFail()
			g.GameOver = true
			g.SetNextStage()
			return false
		} else {
			if !g.Cell(idx).IsMarked() {
				g.moveCount++
				g.Cell(idx).SetMarked()
			}
			if g.moveCount == int(g.level) {
				g.Win = true
				g.SetNextStage()
			}
			return true
		}
	case Conclusion:
	}
	return true
}

func (g *Game) Stage() GameStage { return g.stage }
func (g *Game) SetNextStage() {
	g.stage = g.stage.Next()
	switch g.stage {
	case Show:
		g.gameData.SetBeginShow()
	case Recollection:
		g.gameData.SetBeginResolve()
	case Conclusion:
		g.gameData.SetEndResolve()
		if g.Win {
			g.gameData.SetScore(g.moveCount)
		} else {
			g.gameData.SetScore(0)
		}
	}
}

func (g *Game) Cell(idx int) *Cell       { return g.field.cell(idx) }
func (g *Game) IsCellEmpty(idx int) bool { return g.field.cell(idx).IsEmpty() }

func (g *Game) GameData() GameData     { return *g.gameData }
func (g *Game) Level() Level           { return g.level }
func (g *Game) Dim() Dim               { return g.dim }
func (g *Game) Idx(x, y int) int       { return g.dim.Idx(x, y) }
func (g *Game) Pos(idx int) (int, int) { return g.dim.Pos(idx) }

func (g *Game) String() (result string) {
	result += " Memory Matrix "
	result += strconv.Itoa(int(g.level)) + " "
	result += g.dim.String()
	result += " score:" + g.gameData.String()
	return result
}

func (g *Game) StringFull() (result string) {
	result += g.String() + "\n"
	result += g.stage.String() + "\n"
	for y := 0; y < g.dim.Height(); y++ {
		for x := 0; x < g.dim.Width(); x++ {
			result += g.field.cell(g.Idx(x, y)).String()
		}
		result += "\n"
	}
	return result
}
