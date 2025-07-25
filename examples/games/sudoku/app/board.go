package app

import (
	"fmt"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/sudoku/game"
	"golang.org/x/image/colornames"
)

type Board struct {
	*eui.Container
	layoutCells      *eui.Container
	diff             game.Difficult
	dim              game.Dim
	game             *game.Game
	grid             *eui.GridView
	fn               func(*eui.Button)
	showNotes, isWin bool
	highlight        int
	sw               *eui.Stopwatch
	spacing          int
}

func NewBoard(fn func(b *eui.Button)) *Board {
	b := &Board{Container: eui.NewContainer(eui.NewAbsoluteLayout())}
	b.layoutCells = eui.NewContainer(eui.NewSquareGridLayout(2, 2, 1))
	b.Add(b.layoutCells)
	b.fn = fn
	b.grid = eui.NewGridView(2, 2)
	b.grid.Visible(false)
	b.grid.DrawRect = true
	b.grid.Fg(colornames.Red)
	b.grid.Bg(colornames.Black)
	b.Add(b.grid)
	b.sw = eui.NewStopwatch()
	b.spacing = 5
	return b
}

func (b *Board) Setup(dim game.Dim, diff game.Difficult) {
	b.dim = dim
	b.diff = diff
	b.game = game.NewGame(b.dim)
	b.game.Load(diff)
	b.layoutCells.ResetContainer()
	for y := 0; y < b.dim.Size(); y++ {
		for x := 0; x < b.dim.Size(); x++ {
			btn := NewCellIcon(b.dim, b.game.Cell(x, y), b.fn, colornames.Silver, colornames.Black)
			b.game.Cell(x, y).Connect(btn.UpdateData)
			b.layoutCells.Add(btn)
		}
	}
	b.grid.Set(float64(b.dim.H), float64(b.dim.W))
	b.layoutCells.SetLayout(eui.NewSquareGridLayout(float64(b.dim.Size()), float64(b.dim.Size()), float64(b.spacing)))
	b.ShowNotes(true)
	b.isWin = false
	b.sw.Reset()
	if !(b.diff.String() == game.Manual.String()) {
		b.sw.Start()
	}
}

func (b *Board) GetDiffStr() string {
	return b.diff.String() + "(" + strconv.Itoa(b.game.GetPercent()) + "%)"
}

func (b *Board) IsShowNotes() bool { return b.showNotes }
func (b *Board) ShowNotes(value bool) {
	b.showNotes = value
	b.layoutCells.Traverse(func(d eui.Drawabler) {
		c, ok := d.(*CellIcon)
		if ok {
			c.ShowNotes(value)
		}
	}, false)
}

func (b *Board) Undo() {
	b.game.Undo()
	b.game.UpdateAllFieldNotes()
}

func (b *Board) MoveCount() int { return b.game.MovesCount() }

func (b *Board) Move(x, y int) {
	if !b.game.MakeMove(x, y, b.GetHighlightValue()) {
		if b.isWin = b.game.IsWin(); b.isWin {
			b.sw.Stop()
			fmt.Println("Sudoku field collected game completed", b.sw, b.dim, b.diff, b.isWin)
		}
	}
	b.Highlight(strconv.Itoa(b.GetHighlightValue()))
}

func (b *Board) GetHighlightValue() int { return b.highlight }
func (b *Board) Highlight(value string) {
	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic(err)
	}
	b.highlight = int(n)

	b.layoutCells.Traverse(func(d eui.Drawabler) {
		c, ok := d.(*CellIcon)
		if ok {
			c.Highlight(b.highlight)
		}
	}, false)
}

func (d *Board) Draw(surface *ebiten.Image) {
	if d.IsDirty() {
		d.Layout()
	}
	d.Traverse(func(d eui.Drawabler) { d.Draw(surface) }, false)
}

func (d *Board) Visible(value bool) {
	d.Drawable.Visible(value)
	d.Traverse(func(c eui.Drawabler) { c.Visible(value); c.MarkDirty() }, false)
	d.MarkDirty()
}

func (b *Board) Resize(rect eui.Rect[int]) {
	b.SetRect(rect)
	x, y, w, h := b.Rect().GetRect()
	b.layoutCells.Resize(eui.NewRect([]int{x + int(b.spacing/2), y + int(b.spacing/2), w, h}))
	x1, y1, w1, h1 := b.Rect().GetRect()
	b.grid.Resize(eui.NewRect([]int{x1 - int(b.spacing/2), y1 - int(b.spacing/2), w1, h1}))
	b.grid.SetStrokewidth(float64(b.spacing))
}
