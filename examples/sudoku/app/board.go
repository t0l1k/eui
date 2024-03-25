package app

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/sudoku/game"
	"github.com/t0l1k/eui/examples/sudoku/game/dynamo"
)

type Board struct {
	eui.DrawableBase
	dim    int
	diff   game.Difficult
	field  *game.Field
	layout *eui.GridLayoutRightDown
	grid   *eui.GridView
	fn     func(*eui.Button)
	show   bool
}

func NewBoard(fn func(b *eui.Button)) *Board {
	b := &Board{}
	b.fn = fn
	b.layout = eui.NewGridLayoutRightDown(2, 2)
	b.grid = eui.NewGridView(2, 2)
	b.grid.Visible(false)
	b.grid.DrawRect = true
	b.grid.Fg(eui.Red)
	b.grid.Bg(eui.Black)
	b.Add(b.grid)
	return b
}

func (b *Board) Setup(dim int, diff game.Difficult) {
	b.dim = dim
	size := b.dim * b.dim
	b.diff = diff
	gen := dynamo.NewGenSudokuField(dim, diff)
	b.field = game.NewField(b.dim)
	b.field.Load(gen.GetField())
	b.layout.ResetContainerBase()
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			idx := y*size + x
			btn := NewCellIcon(b.field.GetCells()[idx], b.fn, eui.Silver, eui.Black)
			b.field.GetCells()[idx].Attach(btn)
			b.layout.Add(btn)
		}
	}
	b.grid.Set(float64(dim), float64(dim))
	b.layout.SetDim(float64(size), float64(size))
}

func (b *Board) Update(dt int) {
	if !b.IsVisible() {
		return
	}
	for _, v := range b.layout.GetContainer() {
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
	for _, v := range b.layout.GetContainer() {
		v.Draw(surface)
	}
	for _, v := range b.GetContainer() {
		v.Draw(surface)
	}
}

func (b *Board) IsVisible() bool    { return b.show }
func (b *Board) Visible(value bool) { b.show = value }

func (b *Board) Resize(rect []int) {
	b.layout.Resize(rect)
	b.grid.Resize(rect)
	margin := float64(b.layout.GetRect().GetLowestSize()) * 0.005
	b.layout.SetCellMargin(margin)
	b.grid.SetStrokewidth(margin * 2)
}
