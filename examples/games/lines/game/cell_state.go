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
	return []string{
		"empty",
		"filled next",
		"filled next",
		"marked for move",
		"marked for del",
	}[c]

}

type CellData struct {
	State CellStateType
	Color BallColor
	Pos   eui.Point[int]
}

func (c CellData) String() string {
	return "cell data:" + c.State.String() + "," + c.Color.String() + "," + c.Pos.String()
}

func NewCellData(state CellStateType, color BallColor, pos eui.Point[int]) *CellData {
	return &CellData{State: state, Color: color, Pos: pos}
}

type CellState struct{ *eui.Signal[*CellData] }

func NewCellState(state CellStateType, color BallColor, pos eui.Point[int]) *CellState {
	c := &CellState{Signal: eui.NewSignal[*CellData]()}
	c.Emit(NewCellData(state, color, pos))
	return c
}
