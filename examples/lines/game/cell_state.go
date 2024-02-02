package game

import "github.com/t0l1k/eui"

type CellStateType int

const (
	CellEmpty CellStateType = iota
	CellFilledNext
	CellFilled
	CellMarkedForMove
	CellFilledAfterMove
	CellMarkedForDelete
)

func (c CellStateType) String() (res string) {
	switch c {
	case CellEmpty:
		res = "empty"
	case CellFilledNext:
		res = "filled next"
	case CellFilled:
		res = "filled next"
	case CellMarkedForMove:
		res = "marked for move"
	case CellMarkedForDelete:
		res = "marked for del"
	}
	return res
}

type CellData struct {
	State CellStateType
	Color BallColor
	Pos   *eui.PointInt
}

func (c CellData) String() string {
	return "cell data:" + c.State.String() + "," + c.Color.String() + "," + c.Pos.String()
}

func NewCellData(state CellStateType, color BallColor, pos *eui.PointInt) *CellData {
	return &CellData{State: state, Color: color, Pos: pos}
}

type CellState struct{ eui.SubjectBase }

func NewCellState(state CellStateType, color BallColor, pos *eui.PointInt) *CellState {
	c := &CellState{}
	c.SetValue(NewCellData(state, color, pos))
	return c
}
