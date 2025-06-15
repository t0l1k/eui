package game

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/t0l1k/eui"
)

type MinedField struct {
	State                                              *eui.Signal
	field                                              []*cell
	saved                                              []string
	row, column, totalMines, countMarked, countFlagged int
	firstMove                                          eui.PointInt
}

func NewMinedField(r, c, m int) *MinedField {
	mf := &MinedField{
		row:        r,
		column:     c,
		totalMines: m,
		State:      eui.NewSignal(),
	}
	mf.State.Emit(GameStart)
	mf.State.Connect(func(data any) {
		v := data.(string)
		switch v {
		case GameStart:
		case GamePlay:
		case GameOver:
			mf.resetMarkedMines()
			for _, cell := range mf.field {
				if cell.mined {
					switch cell.State.Value().(*cellData).State() {
					case closed:
						cell.State.Emit(newCellData(blown, cellBlown, cell.pos))
					case questioned:
						cell.State.Emit(newCellData(blown, cellBlown, cell.pos))
					case flagged:
						cell.State.Emit(newCellData(saved, cellBlown, cell.pos))
					}
				} else if cell.State.Value().(*cellData).State() == flagged {
					if !cell.mined {
						cell.State.Emit(newCellData(wrongFlagged, cellWrongFlagged, cell.pos))
					}
				}
			}
		case GameWin:
			for _, cell := range mf.field {
				if cell.mined {
					switch cell.State.Value().(*cellData).State() {
					case closed:
						cell.State.Emit(newCellData(saved, cellSaved, cell.pos))
					case questioned:
						cell.State.Emit(newCellData(saved, cellSaved, cell.pos))
					case flagged:
						cell.State.Emit(newCellData(saved, cellSaved, cell.pos))
					}
				}
			}
		}
		log.Println("MinedField:State:", v)
	})
	return mf
}

func (mf *MinedField) New() {
	mf.State.Emit(GameStart)
	mf.resetMarkedMines()
	if mf.field != nil {
		mf.field = mf.field[:0]
	}
	for y := 0; y < mf.column; y++ {
		for x := 0; x < mf.row; x++ {
			cell := newCell(eui.NewPointInt(x, y))
			cell.State.Connect(mf.UpdateCell)
			mf.field = append(mf.field, cell)
		}
	}
}

func (mf *MinedField) Reset() {
	mf.State.Emit(GameStart)
	mf.resetMarkedMines()
	for _, cell := range mf.field {
		cell.reset()
	}
	mf.State.Emit(GamePlay)
	mf.Open(mf.firstMove.X, mf.firstMove.Y)
}

func (mf *MinedField) GetField() []*cell         { return mf.field }
func (mf *MinedField) GetRow() int               { return mf.row }
func (mf *MinedField) GetColumn() int            { return mf.column }
func (mf *MinedField) GetTotalMines() int        { return mf.totalMines }
func (mf *MinedField) GetMarkedMines() int       { return mf.countFlagged }
func (mf *MinedField) GetPos(idx int) (int, int) { return idx % mf.row, idx / mf.row }
func (mf *MinedField) GetIdx(x, y int) int       { return y*mf.row + x }
func (mx *MinedField) GetCell(x, y int) *cell    { return mx.field[mx.GetIdx(x, y)] }
func (mf *MinedField) isFieldEdge(x, y int) bool {
	return x < 0 || x > mf.row-1 || y < 0 || y > mf.column-1
}
func (mf *MinedField) IsCellOpen(idx int) bool {
	return mf.field[idx].State.Value().(*cellData).State() == opened
}

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
				mf.countMarked++
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
	mf.State.Emit(GamePlay)
}

func (mf *MinedField) Open(x, y int) {
	if mf.isFieldEdge(x, y) {
		return
	}
	cell := mf.GetCell(x, y)
	if cell.State.Value().(*cellData).State() == flagged || cell.State.Value().(*cellData).State() == opened {
		return
	}
	cell.open()
	mf.countMarked++
	if cell.mined {
		cell.State.Emit(newCellData(firstMined, cellFirstMined, cell.pos))
		mf.GameOver()
		return
	}
	if mf.Winned() {
		mf.State.Emit(GameWin)
		return
	}
	if cell.count > 0 {
		return
	}
	for _, newCell := range mf.getNeighbours(x, y) {
		mf.Open(newCell.pos.X, newCell.pos.Y)
	}
}

func (mf *MinedField) Winned() bool { return mf.countMarked == mf.column*mf.row }
func (mf *MinedField) GameOver()    { mf.State.Emit(GameOver) }

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
	if mf.GetCell(x, y).State.Value().(*cellData).State() == opened {
		var countClosed, countFlags byte
		cellValue := mf.GetCell(x, y).count
		neighbours := mf.getNeighbours(x, y)
		for _, cell := range neighbours {
			if cell.State.Value().(*cellData).State() == flagged {
				countFlags++
			} else if cell.State.Value().(*cellData).State() == closed || cell.State.Value().(*cellData).State() == questioned {
				countClosed++
			}
		}
		if countClosed+countFlags == cellValue {
			for _, cell := range neighbours {
				if cell.State.Value().(*cellData).State() == closed || cell.State.Value().(*cellData).State() == questioned {
					cell.State.Emit(newCellData(flagged, cellFlagged, cell.pos))
					changed = true
					mf.countFlagged++
				}
			}
		} else if countFlags == cellValue {
			for _, cell := range neighbours {
				if cell.State.Value().(*cellData).State() == closed || cell.State.Value().(*cellData).State() == questioned {
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
		value := cell.State.Value().(*cellData)
		mf.saved = append(mf.saved, value.State())
	}
}

func (mf *MinedField) RestoreGame() {
	for idx, cell := range mf.field {
		cell.State.Emit(newCellData(mf.saved[idx], cell.String(), cell.pos))
	}
	if mf.State.Value() == GameWin || mf.State.Value() == GameOver {
		mf.State.Emit(GamePlay)
	}
}

func (mf *MinedField) UpdateCell(value interface{}) {
	switch v := value.(type) {
	case *cellData:
		mf.setMarkedMines(v)
		log.Println("MinedField:UpdateCell:", v)
	}
}

func (mf *MinedField) resetMarkedMines() {
	mf.countMarked = 0
	mf.countFlagged = 0
}

func (mf *MinedField) setMarkedMines(v *cellData) {
	switch mf.State.Value() {
	case GamePlay:
		switch v.state {
		case flagged:
			mf.countFlagged++
		case questioned:
			mf.countFlagged--
		}
	}
}

func (mf *MinedField) String() string {
	str := strconv.Itoa(mf.row) + ","
	str += strconv.Itoa(mf.column) + ",("
	str += strconv.Itoa(mf.countFlagged) + "/"
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
