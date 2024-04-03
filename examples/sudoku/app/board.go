package app

import (
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/sudoku/game"
	"github.com/t0l1k/eui/examples/sudoku/game/recur"
)

type Board struct {
	eui.DrawableBase
	diff        game.Difficult
	field       *game.Field
	layoutCells *eui.GridLayoutRightDown
	grid        *eui.GridView
	fn          func(*eui.Button)
	show        bool
	highlight   int
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
	return b
}

func (b *Board) Setup(dim int, diff game.Difficult) {
	var size = dim * dim
	b.diff = diff
	b.field = game.NewField(dim)
	b.field.Load(recur.LoadSudokuField(dim, diff))
	b.layoutCells.ResetContainerBase()
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			idx := y*size + x
			btn := NewCellIcon(b.field.GetCells()[idx], b.fn, eui.Silver, eui.Black)
			b.field.GetCells()[idx].Attach(btn)
			b.layoutCells.Add(btn)
		}
	}
	b.grid.Set(float64(dim), float64(dim))
	b.layoutCells.SetDim(float64(size), float64(size))
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
	b.layoutCells.Resize(rect)
	b.grid.Resize(rect)
	margin := float64(b.layoutCells.GetRect().GetLowestSize()) * 0.005
	b.layoutCells.SetCellMargin(margin)
	b.grid.SetStrokewidth(margin * 2)
}
