package app

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/sudoku/game"
)

type Board struct {
	eui.DrawableBase
	dim    int
	field  *game.Field
	layout *eui.GridLayoutRightDown
	grid   *eui.GridView
}

func NewBoard() *Board {
	b := &Board{}
	b.dim = 2
	b.layout = eui.NewGridLayoutRightDown(b.dim, b.dim)
	b.grid = eui.NewGridView(2, 2)
	b.grid.Visible(false)
	b.grid.DrawRect = true
	b.grid.Fg(eui.Red)
	b.grid.Bg(eui.Black)
	b.Add(b.grid)
	return b
}

func (b *Board) Setup(dim int) {
	b.dim = dim
	size := b.dim * b.dim
	b.layout.ResetContainerBase()
	b.field = game.NewField(b.dim)
	go b.field.New()
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			btn := eui.NewButton(" ", b.buttonsLogic)
			b.field.GetCells()[b.field.Idx(x, y)].Attach(btn)
			b.layout.Add(btn)
		}
	}
	b.grid.Set(dim, dim)
	b.layout.SetDim(size, size)
}

func (b *Board) Visible(value bool) {
	for _, v := range b.GetContainer() {
		switch vT := v.(type) {
		case *eui.GridView:
			vT.Visible(value)
		}
	}
	for _, v := range b.layout.GetContainer() {
		switch vT := v.(type) {
		case *eui.Button:
			vT.Visible(value)
			if value {
				vT.Enable()
			} else {
				vT.Disable()
			}
		}
	}
}

func (b *Board) buttonsLogic(btn *eui.Button) {
	fmt.Println("pressed")
}

func (b *Board) Update(dt int) {
	for _, v := range b.layout.GetContainer() {
		v.Update(dt)
	}
	for _, v := range b.GetContainer() {
		v.Update(dt)
	}
}

func (b *Board) Draw(surface *ebiten.Image) {
	// if !b.IsVisible() {
	// 	return
	// }
	for _, v := range b.layout.GetContainer() {
		v.Draw(surface)
	}
	for _, v := range b.GetContainer() {
		v.Draw(surface)
	}
}

func (b *Board) Resize(rect []int) {
	b.layout.Resize(rect)
	b.grid.Resize(rect)
	margin := float64(b.layout.GetRect().GetLowestSize()) * 0.005
	b.layout.SetCellMargin(int(margin))
	b.grid.SetStrokewidth(int(margin) * 2)
}
