package app

import (
	"github.com/hajimehoshi/ebiten/v2"
	ui "github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/minesweeper/mines"
	"github.com/t0l1k/eui/examples/minesweeper/res"
)

type Game struct {
	ui.ContainerDefault
	field  *mines.MinedField
	timer  *mines.Timer
	fEmoji bool
}

func newGame(r, c, m int) *Game {
	g := &Game{}
	g.field = mines.NewMinedField(r, c, m)
	g.timer = mines.NewTimer()
	rect := []int{0, 0, 1, 1}
	for i := 0; i < len(g.field.GetField()); i++ {
		btn := ui.NewButtonIcon([]*ebiten.Image{res.CellsUp[0], res.CellsUp[1]}, rect, g.gameLogic)
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
			g.Container[idx].(*ui.ButtonIcon).SetIconRelesed(res.CellsUp[0])
			g.Container[idx].(*ui.ButtonIcon).SetIconPressed(res.CellsUp[1])
		case "0":
			g.Container[idx].(*ui.ButtonIcon).SetIconRelesed(res.CellsUp[1])
			g.Container[idx].(*ui.ButtonIcon).SetIconPressed(res.CellsUp[1])
		case "F", "v":
			g.Container[idx].(*ui.ButtonIcon).SetIconRelesed(res.CellsUp[2])
			g.Container[idx].(*ui.ButtonIcon).SetIconPressed(res.CellsUp[2])
		case "Q":
			g.Container[idx].(*ui.ButtonIcon).SetIconRelesed(res.CellsUp[3])
			g.Container[idx].(*ui.ButtonIcon).SetIconPressed(res.CellsUp[3])
		case "f":
			g.Container[idx].(*ui.ButtonIcon).SetIconRelesed(res.CellsUp[6])
			g.Container[idx].(*ui.ButtonIcon).SetIconPressed(res.CellsUp[6])
		case "b":
			g.Container[idx].(*ui.ButtonIcon).SetIconRelesed(res.CellsUp[5])
			g.Container[idx].(*ui.ButtonIcon).SetIconPressed(res.CellsUp[5])
		case "w":
			g.Container[idx].(*ui.ButtonIcon).SetIconRelesed(res.CellsUp[7])
			g.Container[idx].(*ui.ButtonIcon).SetIconPressed(res.CellsUp[7])

		case "1":
			g.Container[idx].(*ui.ButtonIcon).SetIconRelesed(res.CellsDown[0])
			g.Container[idx].(*ui.ButtonIcon).SetIconPressed(res.CellsDown[0])
		case "2":
			g.Container[idx].(*ui.ButtonIcon).SetIconRelesed(res.CellsDown[1])
			g.Container[idx].(*ui.ButtonIcon).SetIconPressed(res.CellsDown[1])
		case "3":
			g.Container[idx].(*ui.ButtonIcon).SetIconRelesed(res.CellsDown[2])
			g.Container[idx].(*ui.ButtonIcon).SetIconPressed(res.CellsDown[2])
		case "4":
			g.Container[idx].(*ui.ButtonIcon).SetIconRelesed(res.CellsDown[3])
			g.Container[idx].(*ui.ButtonIcon).SetIconPressed(res.CellsDown[3])
		case "5":
			g.Container[idx].(*ui.ButtonIcon).SetIconRelesed(res.CellsDown[4])
			g.Container[idx].(*ui.ButtonIcon).SetIconPressed(res.CellsDown[4])
		case "6":
			g.Container[idx].(*ui.ButtonIcon).SetIconRelesed(res.CellsDown[5])
			g.Container[idx].(*ui.ButtonIcon).SetIconPressed(res.CellsDown[5])
		case "7":
			g.Container[idx].(*ui.ButtonIcon).SetIconRelesed(res.CellsDown[6])
			g.Container[idx].(*ui.ButtonIcon).SetIconPressed(res.CellsDown[6])
		case "8":
			g.Container[idx].(*ui.ButtonIcon).SetIconRelesed(res.CellsDown[7])
			g.Container[idx].(*ui.ButtonIcon).SetIconPressed(res.CellsDown[7])
		}
	}
}

func (g *Game) gameLogic(b *ui.ButtonIcon) {
	g.fEmoji = false
	switch g.field.GetState() {
	case mines.GameStart:
		for i, v := range g.Container {
			if b == v {
				x, y := g.field.GetPos(i)
				if b.IsMouseDownLeft() {
					g.field.Shuffle(x, y)
					g.field.Open(x, y)
					g.timer.Start()
					break
				}
			}
		}
	case mines.GamePlay:
		g.field.SaveGame()
		for i, v := range g.Container {
			if b == v {
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
	case mines.GamePause:
	case mines.GameWin:
	case mines.GameOver:
	}
	g.redraw()
}

func (g *Game) Update(dt int) {
	g.timer.Update(dt)
	switch g.field.GetState() {
	case mines.GameStart, mines.GamePlay:
		for _, cell := range g.Container {
			if cell.(*ui.ButtonIcon).IsDisabled() {
				cell.(*ui.ButtonIcon).Enable()
			}
			if cell.(*ui.ButtonIcon).IsMouseDownLeft() && !g.fEmoji {
				g.fEmoji = true
			}
		}
	case mines.GamePause:
		if !ebiten.IsFocused() {
			g.timer.Pause()
		}
	case mines.GameWin, mines.GameOver:
		g.timer.Stop()
		for _, cell := range g.Container {
			if !cell.(*ui.ButtonIcon).IsDisabled() {
				cell.(*ui.ButtonIcon).Disable()
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

func (g *Game) Resize() {
	w, h := ebiten.WindowSize()
	hTop := int(float64(h) * 0.05) // topbar height
	rect := ui.NewRect([]int{0, hTop, w, h - hTop})
	cellSize := g.getCellSize(rect)
	marginX := (rect.W - cellSize*g.field.GetRow()) / 2
	marginY := (rect.H - cellSize*g.field.GetColumn()) / 2
	for idx, cell := range g.Container {
		x := rect.X + cellSize*(idx%g.field.GetRow()) + marginX
		y := rect.Y + cellSize*(idx/g.field.GetRow()) + marginY
		r := []int{x, y, cellSize - 1, cellSize - 1}
		cell.(*ui.ButtonIcon).Resize(r)
	}
}

func (g *Game) getCellSize(rect *ui.Rect) int {
	var size int
	r := g.field.GetRow() + 1
	c := g.field.GetColumn() + 1
	for r*size < rect.W && c*size < rect.H {
		size += 1
	}
	return size
}

func (g *Game) Close() {
	for _, cell := range g.Container {
		cell.Close()
	}
}
