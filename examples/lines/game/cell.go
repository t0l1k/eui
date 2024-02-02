package game

import (
	"fmt"

	"github.com/t0l1k/eui"
)

type Cell struct {
	State *CellState
	color BallColor
	pos   *eui.PointInt
}

func NewCell(x, y int) *Cell {
	pos := eui.NewPointInt(x, y)
	return &Cell{pos: pos,
		color: BallNoColor,
		State: NewCellState(CellEmpty, BallNoColor, pos),
	}
}

func (c *Cell) Pos() *eui.PointInt {
	return c.pos
}

func (c *Cell) Reset() {
	c.color = BallNoColor
	c.State.SetValue(NewCellData(CellEmpty, BallNoColor, c.pos))
}

func (c *Cell) Color() BallColor {
	return c.color
}

func (c *Cell) SetColor(color BallColor) {
	c.color = color
}

func (c *Cell) EqualColors(color BallColor) bool {
	return c.color == color
}

func (c *Cell) IsFilledNext() bool {
	return c.State.Value().(*CellData).State == CellFilledNext
}

func (c *Cell) SetFilledNext(color BallColor) {
	c.color = color
	c.State.SetValue(NewCellData(CellFilledNext, color, c.pos))
}

func (c *Cell) IsFilled() bool {
	return c.State.Value().(*CellData).State == CellFilled
}

func (c *Cell) SetFilled() {
	c.State.SetValue(NewCellData(CellFilled, c.color, c.pos))
}

func (c *Cell) IsMarkedForMove() bool {
	return c.State.Value().(*CellData).State == CellMarkedForMove
}

func (c *Cell) SetMarkedForMove() {
	c.State.SetValue(NewCellData(CellMarkedForMove, c.color, c.pos))
}

func (c *Cell) IsFilledAfterMove() bool {
	return c.State.Value().(*CellData).State == CellFilledAfterMove
}

func (c *Cell) SetFilledAfterMove(color BallColor) {
	c.SetColor(color)
	c.State.SetValue(NewCellData(CellFilledAfterMove, c.color, c.pos))
}

func (c *Cell) IsMarkedForDelete() bool {
	return c.State.Value().(*CellData).State == CellMarkedForDelete
}

func (c *Cell) SetMarkedForDelete() {
	c.State.SetValue(NewCellData(CellMarkedForDelete, c.color, c.pos))
}

func (c *Cell) IsEmpty() bool {
	return c.color == BallNoColor
}

func (c *Cell) String() (result string) {
	switch c.color {
	case BallNoColor:
		result += fmt.Sprintf("[%.03v]", "...")
	default:
		result += fmt.Sprintf("[%.03v]", c.color)
	}
	return result
}
