package game

import "strconv"

type position struct {
	x, y int
}

type cellState string

const (
	closed       cellState = "closed"        // начальное состояние
	flagged      cellState = "flagged"       // отмечена флагом
	questioned   cellState = "questioned"    // отмечена вопросом
	saved        cellState = "saved"         // если помечена флагом
	blown        cellState = "blown"         // есть мина, не отмечена флагом
	wrongFlagged cellState = "wrong flagged" // помечена флагом, без мины
	opened       cellState = "opened"        // открыта
	firstMined   cellState = "first mined"   // взорваная первой
	cellEmpty              = 9
)

type cell struct {
	state cellState
	count byte
	mined bool
	position
}

func newCell(p position) *cell {
	return &cell{
		state:    closed,
		count:    cellEmpty,
		mined:    false,
		position: p,
	}
}

func (c *cell) reset() {
	c.state = closed
}

func (c *cell) open() {
	if c.state == closed || c.state == questioned {
		c.state = opened
	}
}

func (c *cell) mark() {
	switch c.state {
	case closed:
		c.state = flagged
	case flagged:
		c.state = questioned
	case questioned:
		c.state = closed
	}
}

func (c *cell) Pos() (int, int) {
	return c.position.x, c.position.y
}

func (c *cell) String() string {
	var str string
	switch c.state {
	case closed:
		str += " "
	case flagged:
		str += "F"
	case questioned:
		str += "Q"
	case firstMined:
		str += "f"
	case saved:
		str += "v"
	case blown:
		str += "b"
	case wrongFlagged:
		str += "w"
	case opened:
		if c.count == cellEmpty && c.mined {
			str += "*"
		} else if c.count != cellEmpty && !c.mined {
			switch c.count {
			case 0:
				str += "0"
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
