package game

import (
	"github.com/t0l1k/eui"
)

const (
	closed       string = "closed"        // начальное состояние
	flagged      string = "flagged"       // отмечена флагом
	questioned   string = "questioned"    // отмечена вопросом
	saved        string = "saved"         // если помечена флагом
	blown        string = "blown"         // есть мина, не отмечена флагом
	wrongFlagged string = "wrong flagged" // помечена флагом, без мины
	opened       string = "opened"        // открыта
	firstMined   string = "first mined"   // взорваная первой
	cellEmpty           = 9
)

type cellData struct {
	pos          *eui.PointInt
	state, value string
}

func newCellData(state, value string, pos *eui.PointInt) *cellData {
	return &cellData{pos: pos, state: state, value: value}
}

func (cd *cellData) String() string { return "cell[" + cd.value + "]at:" + cd.pos.String() + cd.state }

// Умею оповестить подписчиков о смене состояния в ячейке
type cellState struct {
	listener []eui.Observer
	data     *cellData
}

func newCellState(state string, pos *eui.PointInt) *cellState {
	s := &cellState{}
	s.SetValue(newCellData(state, cellClosed, pos))
	return s
}

func (c cellState) Value() string { return c.data.state }

func (s *cellState) SetValue(value *cellData) {
	if s.data == value {
		return
	}
	s.data = value
	s.Notify()
}

func (s *cellState) Attach(o eui.Observer) { s.listener = append(s.listener, o) }

func (s *cellState) Detach(o eui.Observer) {
	for i, observer := range s.listener {
		if observer == o {
			s.listener = append(s.listener[:i], s.listener[i+1:]...)
			break
		}
	}
}

func (s *cellState) Notify() {
	for _, observer := range s.listener {
		observer.UpdateData(s.data)
	}
}
