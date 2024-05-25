package app

import (
	"fmt"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/sudoku/game"
)

type Board struct {
	eui.DrawableBase
	diff                   game.Difficult
	dim                    game.Dim
	game                   *game.Game
	layoutCells            *eui.GridLayoutRightDown
	grid                   *eui.GridView
	fn                     func(*eui.Button)
	show, showNotes, isWin bool
	highlight              int
	sw                     *eui.Stopwatch
}

func NewBoard(fn func(b *eui.Button)) *Board {
	b := &Board{}
	b.fn = fn
	b.layoutCells = eui.NewGridLayoutRightDown(2, 2)
	b.grid = eui.NewGridView(2, 2)
	b.grid.Visible(false)
	b.grid.DrawRect = true
	b.grid.Fg(eui.Red)
	b.grid.Bg(eui.Black)
	b.Add(b.grid)
	b.sw = eui.NewStopwatch()
	return b
}

func (b *Board) Setup(dim game.Dim, diff game.Difficult) {
	b.dim = dim
	b.diff = diff
	b.game = game.NewGame(b.dim)
	b.game.Load(diff)
	b.layoutCells.ResetContainerBase()
	for y := 0; y < b.dim.Size(); y++ {
		for x := 0; x < b.dim.Size(); x++ {
			btn := NewCellIcon(b.dim, b.game.Cell(x, y), b.fn, eui.Silver, eui.Black)
			b.game.Cell(x, y).Attach(btn)
			b.layoutCells.Add(btn)
		}
	}
	b.grid.Set(float64(b.dim.H), float64(b.dim.W))
	b.layoutCells.SetDim(float64(b.dim.Size()), float64(b.dim.Size()))
	b.ShowNotes(true)
	b.isWin = false
	b.sw.Reset()
	b.sw.Start()
}

func (b *Board) IsShowNotes() bool { return b.showNotes }
func (b *Board) ShowNotes(value bool) {
	b.showNotes = value
	for _, v := range b.layoutCells.GetContainer() {
		v.(*CellIcon).ShowNotes(b.showNotes)
	}
}

func (b *Board) Undo() {
	b.game.Undo()
	b.game.UpdateAllFieldNotes()
}

func (b *Board) MoveCount() int { return b.game.MovesCount() }

func (b *Board) Move(x, y int) {
	if !b.game.MakeMove(x, y, b.GetHighlightValue()) {
		b.isWin = b.game.IsWin()
		b.sw.Stop()
		fmt.Println("Sudoku field collected game completed", b.sw, b.dim, b.diff, b.isWin)
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
	for _, v := range b.layoutCells.GetContainer() {
		v.(*CellIcon).Highlight(b.highlight)
	}
}

func (b *Board) Update(dt int) {
	if !b.IsVisible() {
		return
	}
	for _, v := range b.layoutCells.GetContainer() {
		v.Update(dt)
	}
	for _, v := range b.GetContainer() {
		v.Update(dt)
	}
}

func (b *Board) Draw(surface *ebiten.Image) {
	if !b.IsVisible() {
		return
	}
	for _, v := range b.layoutCells.GetContainer() {
		v.Draw(surface)
	}
	for _, v := range b.GetContainer() {
		v.Draw(surface)
	}
}

func (b *Board) IsVisible() bool    { return b.show }
func (b *Board) Visible(value bool) { b.show = value; b.grid.Visible(value) }

func (b *Board) Resize(rect []int) {
	b.Rect(eui.NewRect(rect))
	x, y := b.GetRect().Pos()
	w, h := b.GetRect().Size()
	b.layoutCells.Resize([]int{x, y, w, h})
	b.grid.Resize([]int{x, y, w, h})
	margin := float64(b.layoutCells.GetRect().GetLowestSize()) * 0.005
	b.layoutCells.SetCellMargin(margin)
	b.grid.SetStrokewidth(margin * 2)
}
