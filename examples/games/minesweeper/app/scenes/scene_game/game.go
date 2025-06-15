package scene_game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/minesweeper/app/game"
)

type Game struct {
	eui.DrawableBase
	field  *game.MinedField
	timer  *eui.Stopwatch
	layout *eui.GridLayoutRightDown
}

func newGame(r, c, m int) *Game {
	g := &Game{}
	g.layout = eui.NewGridLayoutRightDown(float64(r), float64(c))
	g.layout.SetCellMargin(1)
	g.field = game.NewMinedField(r, c, m)
	g.field.State.Connect(g.UpdateData)
	g.timer = eui.NewStopwatch()
	g.New()
	return g
}

func (g *Game) setupBoard() {
	g.field.New()
	var firstStart bool
	if g.layout.GetContainer() == nil {
		firstStart = true
	}
	for i := 0; i < len(g.field.GetField()); i++ {
		var btn *game.CellIcon
		if firstStart {
			btn = game.NewCellIcon(g.field, g.gameLogic)
		} else {
			btn = g.layout.GetContainer()[i].(*game.CellIcon)
			btn.Setup(g.field, g.gameLogic)
		}
		x, y := g.field.GetPos(i)
		cell := g.field.GetCell(x, y)
		cell.State.Connect(btn.UpdateData)
		g.layout.Add(btn)
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
		for i, v := range g.layout.GetContainer() {
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
		for i, v := range g.layout.GetContainer() {
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
			for _, cell := range g.layout.GetContainer() {
				cell.(*game.CellIcon).Visible(false)
			}
		}
	case game.GamePause:
		if ebiten.IsFocused() {
			g.field.State.Emit(game.GamePlay)
			g.timer.Start()
			for _, cell := range g.layout.GetContainer() {
				cell.(*game.CellIcon).Visible(true)
			}
		}
	}
	for _, cell := range g.layout.GetContainer() {
		cell.Update(dt)
	}
}

func (g *Game) UpdateData(value interface{}) {
	switch v := value.(type) {
	case string:
		switch v {
		case game.GameStart, game.GamePlay:
			for _, cell := range g.layout.GetContainer() {
				if cell.(*game.CellIcon).Btn.IsDisabled() {
					cell.(*game.CellIcon).Btn.Enable()
				}
			}
		case game.GameWin, game.GameOver:
			g.timer.Stop()
			for _, cell := range g.layout.GetContainer() {
				if !cell.(*game.CellIcon).Btn.IsDisabled() {
					cell.(*game.CellIcon).Btn.Disable()
				}
			}
		}
	}
}

func (g *Game) Draw(surface *ebiten.Image) {
	for _, cell := range g.layout.GetContainer() {
		cell.Draw(surface)
	}
}

func (g *Game) Resize(rect []int) {
	g.layout.Resize(rect)
}
