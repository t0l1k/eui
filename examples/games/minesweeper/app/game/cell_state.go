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
	pos          eui.PointInt
	state, value string
}

func newCellData(state, value string, pos eui.PointInt) *cellData {
	return &cellData{pos: pos, state: state, value: value}
}
func (c *cellData) State() string   { return c.state }
func (cd *cellData) String() string { return "cell[" + cd.value + "]at:" + cd.pos.String() + cd.state }
