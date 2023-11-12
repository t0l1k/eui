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
	g.field = game.NewMinedField(r, c, m)
	g.timer = eui.NewStopwatch()
	for i := 0; i < len(g.field.GetField()); i++ {
		btn := eui.NewButton("", g.gameLogic)
		g.layout.Add(btn)
	}
	g.Add(g.layout)
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
			g.layout.Container[idx].(*eui.Button).SetText(" ")
			g.layout.Container[idx].(*eui.Button).Bg(eui.Gray)
			g.layout.Container[idx].(*eui.Button).Fg(eui.Red)
		case "0":
			g.layout.Container[idx].(*eui.Button).SetText(" ")
			g.layout.Container[idx].(*eui.Button).Bg(eui.Silver)
			g.layout.Container[idx].(*eui.Button).Fg(eui.Red)
		case "F":
			g.layout.Container[idx].(*eui.Button).SetText("F")
		case "v":
			g.layout.Container[idx].(*eui.Button).SetText("v")
			g.layout.Container[idx].(*eui.Button).Bg(eui.Silver)
			g.layout.Container[idx].(*eui.Button).Fg(eui.Aqua)
		case "Q":
			g.layout.Container[idx].(*eui.Button).SetText("Q")
		case "f":
			g.layout.Container[idx].(*eui.Button).SetText("f")
			g.layout.Container[idx].(*eui.Button).Bg(eui.Red)
			g.layout.Container[idx].(*eui.Button).Fg(eui.Black)
		case "b":
			g.layout.Container[idx].(*eui.Button).SetText("b")
			g.layout.Container[idx].(*eui.Button).Bg(eui.Red)
			g.layout.Container[idx].(*eui.Button).Fg(eui.Gray)
		case "w":
			g.layout.Container[idx].(*eui.Button).SetText("w")
			g.layout.Container[idx].(*eui.Button).Bg(eui.Red)
			g.layout.Container[idx].(*eui.Button).Fg(eui.Gray)

		case "1":
			g.layout.Container[idx].(*eui.Button).SetText("1")
			g.layout.Container[idx].(*eui.Button).Bg(eui.Silver)
			g.layout.Container[idx].(*eui.Button).Fg(eui.Blue)
		case "2":
			g.layout.Container[idx].(*eui.Button).SetText("2")
			g.layout.Container[idx].(*eui.Button).Bg(eui.Silver)
			g.layout.Container[idx].(*eui.Button).Fg(eui.Orange)
		case "3":
			g.layout.Container[idx].(*eui.Button).SetText("3")
			g.layout.Container[idx].(*eui.Button).Bg(eui.Silver)
			g.layout.Container[idx].(*eui.Button).Fg(eui.Green)
		case "4":
			g.layout.Container[idx].(*eui.Button).SetText("4")
			g.layout.Container[idx].(*eui.Button).Bg(eui.Silver)
			g.layout.Container[idx].(*eui.Button).Fg(eui.Aqua)
		case "5":
			g.layout.Container[idx].(*eui.Button).SetText("5")
			g.layout.Container[idx].(*eui.Button).Bg(eui.Silver)
			g.layout.Container[idx].(*eui.Button).Fg(eui.Navy)
		case "6":
			g.layout.Container[idx].(*eui.Button).SetText("6")
			g.layout.Container[idx].(*eui.Button).Bg(eui.Silver)
			g.layout.Container[idx].(*eui.Button).Fg(eui.Fuchsia)
		case "7":
			g.layout.Container[idx].(*eui.Button).SetText("7")
			g.layout.Container[idx].(*eui.Button).Bg(eui.Silver)
			g.layout.Container[idx].(*eui.Button).Fg(eui.Purple)
		case "8":
			g.layout.Container[idx].(*eui.Button).SetText("8")
			g.layout.Container[idx].(*eui.Button).Bg(eui.Silver)
			g.layout.Container[idx].(*eui.Button).Fg(eui.Black)
		}
	}
}

func (g *Game) gameLogic(b *eui.Button) {
	defer g.redraw()
	switch g.field.GetState() {
	case game.GameStart:
		for i, v := range g.layout.Container {
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
		for i, v := range g.layout.Container {
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
		for _, cell := range g.layout.Container {
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
		for _, cell := range g.layout.Container {
			if !cell.(*eui.Button).IsDisabled() {
				cell.(*eui.Button).Disable()
			}
		}
	}
	for _, cell := range g.layout.Container {
		cell.Update(dt)
	}
}

func (g *Game) Draw(surface *ebiten.Image) {
	for _, cell := range g.layout.Container {
		cell.Draw(surface)
	}
}
