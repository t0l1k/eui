package scene_game

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/minesweeper/app/game"
)

type Game struct {
	eui.View
	field  *game.MinedField
	timer  *eui.Stopwatch
	layout *eui.GridLayoutRightDown
}

func newGame(r, c, m int) *Game {
	g := &Game{}
	g.layout = eui.NewGridLayoutRightDown(r, c)
	g.layout.SetCellMargin(1)
	g.Add(g.layout)
	g.field = game.NewMinedField(r, c, m)
	g.field.State.Attach(g)
	g.timer = eui.NewStopwatch()
	g.New()
	return g
}

func (g *Game) setupBoard() {
	g.field.New()
	var firstStart bool
	if g.layout.Container == nil {
		firstStart = true
	}
	for i := 0; i < len(g.field.GetField()); i++ {
		var btn *game.CellIcon
		if firstStart {
			btn = game.NewCellIcon(g.field, g.gameLogic)
		} else {
			btn = g.layout.Container[i].(*game.CellIcon)
			btn.Setup(g.field, g.gameLogic)
		}
		x, y := g.field.GetPos(i)
		cell := g.field.GetCell(x, y)
		cell.State.Attach(btn)
		g.layout.Add(btn)
	}
}

func (g *Game) New() {
	g.timer.Reset()
	g.setupBoard()
	g.redraw()
}

func (g *Game) Reset() {
	g.timer.Reset()
	g.field.Reset()
	g.redraw()
	g.timer.Start()
}

func (g *Game) redraw() {
	log.Println(g.field)
}

func (g *Game) gameLogic(b *eui.Button) {
	defer g.redraw()
	switch g.field.State.Value() {
	case game.GameStart:
		for i, v := range g.layout.Container {
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
		for i, v := range g.layout.Container {
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
			g.field.State.SetValue(game.GamePause)
			g.timer.Stop()
			for _, cell := range g.layout.Container {
				cell.(*game.CellIcon).Visible(false)
			}
		}
	case game.GamePause:
		if ebiten.IsFocused() {
			g.field.State.SetValue(game.GamePlay)
			g.timer.Start()
			for _, cell := range g.layout.Container {
				cell.(*game.CellIcon).Visible(true)
			}
		}
	}
	for _, cell := range g.layout.Container {
		cell.Update(dt)
	}
}

func (g *Game) UpdateData(value interface{}) {
	switch v := value.(type) {
	case string:
		switch v {
		case game.GameStart, game.GamePlay:
			for _, cell := range g.layout.Container {
				if cell.(*game.CellIcon).Btn.IsDisabled() {
					cell.(*game.CellIcon).Btn.Enable()
				}
			}
		case game.GameWin, game.GameOver:
			g.timer.Stop()
			for _, cell := range g.layout.Container {
				if !cell.(*game.CellIcon).Btn.IsDisabled() {
					cell.(*game.CellIcon).Btn.Disable()
				}
			}
		}
	}
}

func (g *Game) Draw(surface *ebiten.Image) {
	for _, cell := range g.layout.Container {
		cell.Draw(surface)
	}
}
