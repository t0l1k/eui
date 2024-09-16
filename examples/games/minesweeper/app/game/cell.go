package game

import (
	"strconv"

	"github.com/t0l1k/eui"
)

const (
	cellClosed       = " "
	cellFlagged      = "F"
	cellQuestioned   = "Q"
	cellFirstMined   = "f"
	cellSaved        = "v"
	cellBlown        = "b"
	cellWrongFlagged = "w"
	cellMined        = "*"
	cellZero         = "0"
)

type cell struct {
	State *cellState
	count byte
	mined bool
	pos   eui.PointInt
}

func newCell(pos eui.PointInt) *cell {
	return &cell{
		State: newCellState(closed, pos),
		count: cellEmpty,
		mined: false,
		pos:   pos,
	}
}

func (c *cell) reset() { c.State.SetValue(newCellData(closed, cellClosed, c.pos)) }

func (c *cell) open() {
	if c.State.Value() == closed || c.State.Value() == questioned {
		c.State.SetValue(newCellData(opened, c.String(), c.pos))
	}
}

func (c *cell) mark() {
	switch c.State.Value() {
	case closed:
		c.State.SetValue(newCellData(flagged, cellFlagged, c.pos))
	case flagged:
		c.State.SetValue(newCellData(questioned, cellQuestioned, c.pos))
	case questioned:
		c.State.SetValue(newCellData(closed, cellClosed, c.pos))
	}
}

func (c *cell) Pos() (int, int) { return c.pos.X, c.pos.Y }

func (c *cell) String() string {
	var str string
	switch c.State.Value() {
	case closed:
		str += "."
	case flagged:
		str += cellFlagged
	case questioned:
		str += cellQuestioned
	case firstMined:
		str += cellFirstMined
	case saved:
		str += cellSaved
	case blown:
		str += cellBlown
	case wrongFlagged:
		str += cellWrongFlagged
	case opened:
		if c.count == cellEmpty && c.mined {
			str += cellMined
		} else if c.count != cellEmpty && !c.mined {
			switch c.count {
			case 0:
				str += cellZero
			default:
				str += strconv.Itoa(int(c.count))
			}
		} else {
			str += "!"
		}
	default:
		str += "?"
	}
	return str
}
