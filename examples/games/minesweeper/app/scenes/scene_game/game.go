package scene_game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/minesweeper/app/game"
)

type Game struct {
	*eui.Container
	field *game.MinedField
	timer *eui.Stopwatch
}

func newGame(r, c, m int) *Game {
	g := &Game{Container: eui.NewContainer(eui.NewGridLayout(float64(r), float64(c), 1))}
	g.field = game.NewMinedField(r, c, m)
	g.field.State.Connect(g.UpdateData)
	g.timer = eui.NewStopwatch()
	g.New()
	return g
}

func (g *Game) setupBoard() {
	g.field.New()
	var firstStart bool
	if g.Childrens() == nil {
		firstStart = true
	}
	for i := 0; i < len(g.field.GetField()); i++ {
		var btn *game.CellIcon
		if firstStart {
			btn = game.NewCellIcon(g.field, g.gameLogic)
			g.Add(btn)
		} else {
			btn = g.Childrens()[i].(*game.CellIcon)
			btn.Setup(g.field, g.gameLogic)
		}
		x, y := g.field.GetPos(i)
		cell := g.field.GetCell(x, y)
		cell.State.Connect(btn.UpdateData)
	}
}

func (g *Game) New() {
	g.timer.Reset()
	g.setupBoard()
}

func (g *Game) Reset() {
	g.timer.Reset()
	g.field.Reset()
	g.timer.Start()
}

func (g *Game) gameLogic(b *eui.Button) {
	switch g.field.State.Value() {
	case game.GameStart:
		for i, v := range g.Childrens() {
			if b == v.(*game.CellIcon).Btn {
				x, y := g.field.GetPos(i)
				if b.IsMouseDownLeft() {
					g.field.Shuffle(x, y)
					g.field.Open(x, y)
					g.timer.Start()
					break
				}
			}
		}
	case game.GamePlay:
		g.field.SaveGame()
		for i, v := range g.Childrens() {
			if b == v.(*game.CellIcon).Btn {
				x, y := g.field.GetPos(i)
				if b.IsMouseDownLeft() && !g.field.IsCellOpen(i) {
					g.field.Open(x, y)
				} else if b.IsMouseDownLeft() && g.field.IsCellOpen(i) {
					g.field.AutoMarkFlags(x, y)
				} else if b.IsMouseDownRight() {
					g.field.MarkFlag(x, y)
				}
				break
			}
		}
	}
}

func (g *Game) Update(dt int) {
	switch g.field.State.Value() {
	case game.GameStart, game.GamePlay:
		if !ebiten.IsFocused() && g.field.State.Value() == game.GamePlay {
			g.field.State.Emit(game.GamePause)
			g.timer.Stop()
			for _, cell := range g.Childrens() {
				cell.(*game.CellIcon).Visible(false)
			}
		}
	case game.GamePause:
		if ebiten.IsFocused() {
			g.field.State.Emit(game.GamePlay)
			g.timer.Start()
			for _, cell := range g.Childrens() {
				cell.(*game.CellIcon).Visible(true)
			}
		}
	}
	g.Container.Update(dt)
}

func (g *Game) UpdateData(value string) {
	switch value {
	case game.GameStart, game.GamePlay:
		for _, cell := range g.Childrens() {
			if cell.(*game.CellIcon).Btn.IsDisabled() {
				cell.(*game.CellIcon).Btn.Enable()
			}
		}
	case game.GameWin, game.GameOver:
		g.timer.Stop()
		for _, cell := range g.Childrens() {
			if !cell.(*game.CellIcon).Btn.IsDisabled() {
				cell.(*game.CellIcon).Btn.Disable()
			}
		}
	}
}
