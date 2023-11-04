package scene_game

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/minesweeper/app/game"
)

type Game struct {
	eui.ContainerBase
	field *game.MinedField
	timer *eui.Stopwatch
}

func newGame(r, c, m int) *Game {
	g := &Game{}
	g.field = game.NewMinedField(r, c, m)
	g.timer = eui.NewStopwatch()
	bg := eui.Gray
	fg := eui.Black
	for i := 0; i < len(g.field.GetField()); i++ {
		btn := eui.NewButton("", bg, fg, g.gameLogic)
		g.Add(btn)
	}
	return g
}

func (g *Game) New() {
	g.timer.Reset()
	g.field.New()
	g.redraw()
}

func (g *Game) Reset() {
	g.timer.Reset()
	g.field.Reset()
	g.redraw()
	g.timer.Start()
}

func (g *Game) redraw() {
	for idx, cell := range g.field.GetField() {
		switch cell.String() {
		case " ":
			g.Container[idx].(*eui.Button).SetText(" ")
			g.Container[idx].(*eui.Button).Bg(eui.Gray)
			g.Container[idx].(*eui.Button).Fg(eui.Red)
		case "0":
			g.Container[idx].(*eui.Button).SetText(" ")
			g.Container[idx].(*eui.Button).Bg(eui.Silver)
			g.Container[idx].(*eui.Button).Fg(eui.Red)
		case "F":
			g.Container[idx].(*eui.Button).SetText("F")
		case "v":
			g.Container[idx].(*eui.Button).SetText("v")
			g.Container[idx].(*eui.Button).Bg(eui.Silver)
			g.Container[idx].(*eui.Button).Fg(eui.Aqua)
		case "Q":
			g.Container[idx].(*eui.Button).SetText("Q")
		case "f":
			g.Container[idx].(*eui.Button).SetText("f")
			g.Container[idx].(*eui.Button).Bg(eui.Red)
			g.Container[idx].(*eui.Button).Fg(eui.Black)
		case "b":
			g.Container[idx].(*eui.Button).SetText("b")
			g.Container[idx].(*eui.Button).Bg(eui.Red)
			g.Container[idx].(*eui.Button).Fg(eui.Gray)
		case "w":
			g.Container[idx].(*eui.Button).SetText("w")
			g.Container[idx].(*eui.Button).Bg(eui.Red)
			g.Container[idx].(*eui.Button).Fg(eui.Gray)

		case "1":
			g.Container[idx].(*eui.Button).SetText("1")
			g.Container[idx].(*eui.Button).Bg(eui.Silver)
			g.Container[idx].(*eui.Button).Fg(eui.Blue)
		case "2":
			g.Container[idx].(*eui.Button).SetText("2")
			g.Container[idx].(*eui.Button).Bg(eui.Silver)
			g.Container[idx].(*eui.Button).Fg(eui.Orange)
		case "3":
			g.Container[idx].(*eui.Button).SetText("3")
			g.Container[idx].(*eui.Button).Bg(eui.Silver)
			g.Container[idx].(*eui.Button).Fg(eui.Green)
		case "4":
			g.Container[idx].(*eui.Button).SetText("4")
			g.Container[idx].(*eui.Button).Bg(eui.Silver)
			g.Container[idx].(*eui.Button).Fg(eui.Aqua)
		case "5":
			g.Container[idx].(*eui.Button).SetText("5")
			g.Container[idx].(*eui.Button).Bg(eui.Silver)
			g.Container[idx].(*eui.Button).Fg(eui.Navy)
		case "6":
			g.Container[idx].(*eui.Button).SetText("6")
			g.Container[idx].(*eui.Button).Bg(eui.Silver)
			g.Container[idx].(*eui.Button).Fg(eui.Fuchsia)
		case "7":
			g.Container[idx].(*eui.Button).SetText("7")
			g.Container[idx].(*eui.Button).Bg(eui.Silver)
			g.Container[idx].(*eui.Button).Fg(eui.Purple)
		case "8":
			g.Container[idx].(*eui.Button).SetText("8")
			g.Container[idx].(*eui.Button).Bg(eui.Silver)
			g.Container[idx].(*eui.Button).Fg(eui.Black)
		}
	}
}

func (g *Game) gameLogic(b *eui.Button) {
	defer g.redraw()
	switch g.field.GetState() {
	case game.GameStart:
		for i, v := range g.Container {
			if b == v {
				x, y := g.field.GetPos(i)
				if b.IsMouseDownLeft() {
					g.field.Shuffle(x, y)
					g.field.Open(x, y)
					g.timer.Start()
					log.Println("game: begin new game")
					break
				}
			}
		}
	case game.GamePlay:
		g.field.SaveGame()
		for i, v := range g.Container {
			if b == v {
				x, y := g.field.GetPos(i)
				if b.IsMouseDownLeft() && !g.field.IsCellOpen(i) {
					g.field.Open(x, y)
					log.Println("game: manual cell open at:", x, y)
				} else if b.IsMouseDownLeft() && g.field.IsCellOpen(i) {
					g.field.AutoMarkFlags(x, y)
				} else if b.IsMouseDownRight() {
					g.field.MarkFlag(x, y)
				}
				break
			}
		}
	case game.GamePause:
	case game.GameWin:
	case game.GameOver:
	}
}

func (g *Game) Update(dt int) {
	switch g.field.GetState() {
	case game.GameStart, game.GamePlay:
		for _, cell := range g.Container {
			if cell.(*eui.Button).IsDisabled() {
				cell.(*eui.Button).Enable()
			}
		}
	case game.GamePause:
		if !ebiten.IsFocused() {
			g.timer.Stop()
		}
	case game.GameWin, game.GameOver:
		g.timer.Stop()
		for _, cell := range g.Container {
			if !cell.(*eui.Button).IsDisabled() {
				cell.(*eui.Button).Disable()
			}
		}
	}
	for _, cell := range g.Container {
		cell.Update(dt)
	}
}

func (g *Game) Draw(surface *ebiten.Image) {
	for _, cell := range g.Container {
		cell.Draw(surface)
	}
}

func (g *Game) Resize(r []int) {
	rect := eui.NewRect(r)
	cellSize := g.getCellSize(rect)
	marginX := (rect.W - cellSize*g.field.GetRow()) / 2
	marginY := (rect.H - cellSize*g.field.GetColumn()) / 2
	for idx, cell := range g.Container {
		x := rect.X + cellSize*(idx%g.field.GetRow()) + marginX
		y := rect.Y + cellSize*(idx/g.field.GetRow()) + marginY
		r := []int{x, y, cellSize - 1, cellSize - 1}
		cell.(*eui.Button).Resize(r)
	}
}

func (g *Game) getCellSize(rect *eui.Rect) int {
	var size int
	r := g.field.GetRow() + 1
	c := g.field.GetColumn() + 1
	for r*size < rect.W && c*size < rect.H {
		size += 1
	}
	return size
}
