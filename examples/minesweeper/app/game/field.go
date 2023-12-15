package game

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/t0l1k/eui"
)

type MinedField struct {
	State                                                *gameState
	field                                                []*cell
	saved                                                []string
	row, column, totalMines, markedMines, totalOpenCells int
	firstMove                                            *eui.PointInt
}

func NewMinedField(r, c, m int) *MinedField {
	mf := &MinedField{
		row:        r,
		column:     c,
		totalMines: m,
	}
	mf.State = newGameState()
	return mf
}

func (mf *MinedField) New() {
	mf.State.SetValue(GameStart)
	mf.resetMarkedMines()
	if mf.field != nil {
		mf.field = mf.field[:0]
	}
	for y := 0; y < mf.column; y++ {
		for x := 0; x < mf.row; x++ {
			cell := newCell(eui.NewPointInt(x, y))
			cell.State.Attach(mf)
			mf.field = append(mf.field, cell)
		}
	}
}

func (mf *MinedField) Reset() {
	mf.State.SetValue(GameStart)
	mf.resetMarkedMines()
	for _, cell := range mf.field {
		cell.reset()
	}
	mf.State.SetValue(GamePlay)
	mf.Open(mf.firstMove.X, mf.firstMove.Y)
}

func (mf *MinedField) GetField() []*cell         { return mf.field }
func (mf *MinedField) GetRow() int               { return mf.row }
func (mf *MinedField) GetColumn() int            { return mf.column }
func (mf *MinedField) GetTotalMines() int        { return mf.totalMines }
func (mf *MinedField) GetMarkedMines() int       { return mf.markedMines }
func (mf *MinedField) GetPos(idx int) (int, int) { return idx % mf.row, idx / mf.row }
func (mf *MinedField) GetIdx(x, y int) int       { return y*mf.row + x }
func (mx *MinedField) GetCell(x, y int) *cell    { return mx.field[mx.GetIdx(x, y)] }
func (mf *MinedField) isFieldEdge(x, y int) bool {
	return x < 0 || x > mf.row-1 || y < 0 || y > mf.column-1
}
func (mf *MinedField) IsCellOpen(idx int) bool { return mf.field[idx].State.Value() == opened }

func (mf *MinedField) getNeighbours(x, y int) (cells []*cell) {
	for dy := -1; dy < 2; dy++ {
		for dx := -1; dx < 2; dx++ {
			nx := x + dx
			ny := y + dy
			if !mf.isFieldEdge(nx, ny) {
				newCell := mf.GetCell(nx, ny)
				cells = append(cells, newCell)
			}
		}
	}
	return cells
}

func (mf *MinedField) Shuffle(fX, fY int) {
	if mf.State.Value() != GameStart {
		return
	}
	mf.firstMove = eui.NewPointInt(fX, fY)
	rand.Seed(time.Now().UTC().UnixNano())
	var mines int
	for mines < mf.totalMines {
		x, y := rand.Intn(mf.row), rand.Intn(mf.column)
		if x != fX && y != fY {
			cell := mf.field[mf.GetIdx(x, y)]
			if !cell.mined {
				cell.mined = true
				mines++
			}
		}
	}
	for idx, cell := range mf.field {
		var count byte
		if !cell.mined {
			x, y := mf.GetPos(idx)
			neighbours := mf.getNeighbours(x, y)
			for _, newCell := range neighbours {
				if newCell.mined {
					count++
				}
			}
			mf.field[idx].count = count
		}
	}
	mf.State.SetValue(GamePlay)
}

func (mf *MinedField) Open(x, y int) {
	if mf.isFieldEdge(x, y) {
		return
	}
	cell := mf.GetCell(x, y)
	if cell.State.Value() == flagged || cell.State.Value() == opened {
		return
	}
	cell.open()
	if cell.mined {
		cell.State.SetValue(newCellData(firstMined, cellFirstMined, cell.pos))
		mf.GameOver()
		return
	}
	if mf.Winned() {
		return
	}
	if cell.count > 0 {
		return
	}
	for _, newCell := range mf.getNeighbours(x, y) {
		mf.Open(newCell.pos.X, newCell.pos.Y)
	}
}

func (mf *MinedField) Winned() bool {
	if mf.totalOpenCells+mf.totalMines == mf.row*mf.column {
		mf.State.SetValue(GameWin)
		mf.resetMarkedMines()
		for _, cell := range mf.field {
			if cell.mined {
				cell.State.SetValue(newCellData(saved, cellSaved, cell.pos))
			}
		}
		return true
	}
	return false
}

func (mf *MinedField) GameOver() {
	mf.State.SetValue(GameOver)
	mf.resetMarkedMines()
	for _, cell := range mf.field {
		if cell.mined {
			switch cell.State.Value() {
			case closed:
				cell.State.SetValue(newCellData(blown, cellBlown, cell.pos))
			case questioned:
				cell.State.SetValue(newCellData(blown, cellBlown, cell.pos))
			case flagged:
				cell.State.SetValue(newCellData(saved, cellBlown, cell.pos))
			}
		} else if cell.State.Value() == flagged {
			if !cell.mined {
				cell.State.SetValue(newCellData(wrongFlagged, cellWrongFlagged, cell.pos))
			}
		}
	}
}

func (mf *MinedField) AutoMarkAllFlags() {
	var count int
	for {
		for idx := range mf.field {
			if mf.AutoMarkFlags(mf.GetPos(idx)) {
				count++
			}
		}
		if count > 0 {
			mf.AutoMarkAllFlags()
			count = 0
		} else {
			break
		}
	}
}

func (mf *MinedField) AutoMarkFlags(x, y int) (changed bool) {
	if mf.GetCell(x, y).State.Value() == opened {
		var countClosed, countFlags byte
		cellValue := mf.GetCell(x, y).count
		neighbours := mf.getNeighbours(x, y)
		for _, cell := range neighbours {
			if cell.State.Value() == flagged {
				countFlags++
			} else if cell.State.Value() == closed || cell.State.Value() == questioned {
				countClosed++
			}
		}
		if countClosed+countFlags == cellValue {
			for _, cell := range neighbours {
				if cell.State.Value() == closed || cell.State.Value() == questioned {
					cell.State.SetValue(newCellData(flagged, cellFlagged, cell.pos))
					changed = true
				}
			}
		} else if countFlags == cellValue {
			for _, cell := range neighbours {
				if cell.State.Value() == closed || cell.State.Value() == questioned {
					mf.Open(cell.pos.X, cell.pos.Y)
					changed = true
				}
			}
		}
	}
	return changed
}

func (mf *MinedField) MarkFlag(x, y int) {
	mf.field[mf.GetIdx(x, y)].mark()
}

func (mf *MinedField) SaveGame() {
	if mf.State.Value() != GamePlay {
		return
	}
	mf.saved = mf.saved[:0]
	for _, cell := range mf.field {
		mf.saved = append(mf.saved, cell.State.Value())
	}
}

func (mf *MinedField) RestoreGame() {
	for idx, cell := range mf.field {
		cell.State.SetValue(newCellData(mf.saved[idx], cell.String(), cell.pos))
	}
	if mf.State.Value() == GameWin || mf.State.Value() == GameOver {
		mf.State.SetValue(GamePlay)
	}
}

func (mf *MinedField) UpdateData(value interface{}) {
	switch v := value.(type) {
	case *cellData:
		mf.setMarkedMines(v)
	default:
	}
}

func (mf *MinedField) resetMarkedMines() {
	mf.markedMines = 0
	mf.totalOpenCells = 0
}

func (mf *MinedField) setMarkedMines(v *cellData) {
	switch mf.State.Value() {
	case GamePlay:
		switch v.state {
		case opened:
			mf.totalOpenCells++
		case flagged:
			mf.markedMines++
		case questioned:
			mf.markedMines--
		}
	case GameWin, GameOver:
		if v.state == saved {
			mf.markedMines++
		}
	}
}

func (mf *MinedField) String() string {
	str := strconv.Itoa(mf.row) + ","
	str += strconv.Itoa(mf.column) + ",("
	str += strconv.Itoa(mf.markedMines) + "/"
	str += strconv.Itoa(mf.totalMines) + ")"
	str += mf.State.Value().(string) + "\n"
	for y := 0; y < mf.column; y++ {
		for x := 0; x < mf.row; x++ {
			str += mf.field[mf.GetIdx(x, y)].String()
		}
		str += "\n"
	}
	return str
}
