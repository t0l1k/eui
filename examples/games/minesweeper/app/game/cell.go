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
	State *eui.Signal
	count byte
	mined bool
	pos   eui.PointInt
}

func newCell(pos eui.PointInt) *cell {
	v := newCellData(closed, cellClosed, pos)
	c := &cell{
		State: eui.NewSignal(),
		count: cellEmpty,
		mined: false,
		pos:   pos,
	}
	c.State.Emit(v)
	return c
}

func (c *cell) reset() { c.State.Emit(newCellData(closed, cellClosed, c.pos)) }

func (c *cell) open() {
	cd := c.State.Value().(*cellData)
	if cd.State() == closed || cd.State() == questioned {
		c.State.Emit(newCellData(opened, c.String(), c.pos))
	}
}

func (c *cell) mark() {
	cd := c.State.Value().(*cellData)
	switch cd.State() {
	case closed:
		c.State.Emit(newCellData(flagged, cellFlagged, c.pos))
	case flagged:
		c.State.Emit(newCellData(questioned, cellQuestioned, c.pos))
	case questioned:
		c.State.Emit(newCellData(closed, cellClosed, c.pos))
	}
}

func (c *cell) Pos() (int, int) { return c.pos.X, c.pos.Y }

func (c *cell) String() string {
	var str string
	cd := c.State.Value().(*cellData)
	switch cd.State() {
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
