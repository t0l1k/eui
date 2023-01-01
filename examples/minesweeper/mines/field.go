package mines

import (
	"math/rand"
	"strconv"
	"time"
)

type GameState string

const (
	GameStart GameState = "start"
	GamePlay  GameState = "play"
	GamePause GameState = "pause"
	GameWin   GameState = "win"
	GameOver  GameState = "game over"
)

type MinedField struct {
	state                   GameState
	field                   []*cell
	saved                   []cellState
	row, column, totalMines int
	firstMove               position
}

func NewMinedField(r, c, m int) *MinedField {
	mf := &MinedField{
		row:        r,
		column:     c,
		totalMines: m,
	}
	mf.New()
	return mf
}

func (mf *MinedField) New() {
	mf.state = GameStart
	if mf.field != nil {
		mf.field = mf.field[:0]
	}
	for y := 0; y < mf.column; y++ {
		for x := 0; x < mf.row; x++ {
			mf.field = append(mf.field, newCell(position{x, y}))
		}
	}
}

func (mf *MinedField) Reset() {
	for _, cell := range mf.field {
		cell.reset()
	}
	mf.Open(mf.firstMove.x, mf.firstMove.y)
	mf.state = GamePlay
}

func (mf *MinedField) GetField() []*cell  { return mf.field }
func (mf *MinedField) GetRow() int        { return mf.row }
func (mf *MinedField) GetColumn() int     { return mf.column }
func (mf *MinedField) GetTotalMines() int { return mf.totalMines }
func (mf *MinedField) GetLeftMines() int {
	leftMines := 0
	for _, cell := range mf.GetField() {
		if cell.state == flagged || cell.state == saved {
			leftMines++
		}
	}
	return leftMines
}
func (mf *MinedField) GetState() GameState       { return mf.state }
func (mf *MinedField) GetPos(idx int) (int, int) { return idx % mf.row, idx / mf.row }
func (mf *MinedField) GetIdx(x, y int) int       { return y*mf.row + x }
func (mf *MinedField) isFieldEdge(x, y int) bool {
	return x < 0 || x > mf.row-1 || y < 0 || y > mf.column-1
}

func (mf *MinedField) IsCellOpen(idx int) bool {
	return mf.field[idx].state == opened
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

func (mx *MinedField) GetCell(x, y int) *cell {
	return mx.field[mx.GetIdx(x, y)]
}

func (mf *MinedField) Shuffle(fX, fY int) {
	if mf.state != GameStart {
		return
	}
	mf.firstMove.x, mf.firstMove.y = fX, fY
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
	mf.state = GamePlay
}

func (mf *MinedField) Open(x, y int) {
	if mf.isFieldEdge(x, y) {
		return
	}
	cell := mf.GetCell(x, y)
	if cell.state == flagged || cell.state == opened {
		return
	}
	cell.open()
	if cell.mined {
		cell.state = firstMined
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
		mf.Open(newCell.position.x, newCell.position.y)
	}
}

func (mf *MinedField) Winned() bool {
	var count int
	for _, cell := range mf.field {
		if cell.state == opened {
			count++
		}
	}
	if count+mf.totalMines == mf.row*mf.column {
		for _, cell := range mf.field {
			if cell.mined {
				cell.state = saved
			}
		}
		mf.state = GameWin
		return true
	}
	return false
}

func (mf *MinedField) GameOver() {
	for _, cell := range mf.field {
		if cell.mined {
			switch cell.state {
			case closed:
				cell.state = blown
			case questioned:
				cell.state = blown
			case flagged:
				cell.state = saved
			}
		} else if cell.state == flagged {
			if !cell.mined {
				cell.state = wrongFlagged
			}
		}
	}
	mf.state = GameOver
}

func (mf *MinedField) AutoMarkFlags(x, y int) {
	if mf.GetCell(x, y).state == opened {
		var countClosed, countFlags byte
		cellValue := mf.GetCell(x, y).count
		neighbours := mf.getNeighbours(x, y)
		for _, value := range neighbours {
			if value.state == flagged {
				countFlags++
			} else if value.state == closed || value.state == questioned {
				countClosed++
			}
		}
		if countClosed+countFlags == cellValue {
			for _, value := range neighbours {
				if value.state == closed || value.state == questioned {
					value.state = flagged
				}
			}
		} else if countFlags == cellValue {
			for _, value := range neighbours {
				if value.state == closed || value.state == questioned {
					mf.Open(value.position.x, value.position.y)
				}
			}
		}
	}
}

func (mf *MinedField) MarkFlag(x, y int) {
	mf.field[mf.GetIdx(x, y)].mark()
}

func (mf *MinedField) SaveGame() {
	if mf.state != GamePlay {
		return
	}
	mf.saved = mf.saved[:0]
	for _, cell := range mf.field {
		mf.saved = append(mf.saved, cell.state)
	}
}

func (mf *MinedField) RestoreGame() {
	for idx, cell := range mf.field {
		cell.state = mf.saved[idx]
	}
	if mf.state == GameWin || mf.state == GameOver {
		mf.state = GamePlay
	}
}

func (mf *MinedField) String() string {
	str := strconv.Itoa(mf.row) + ","
	str += strconv.Itoa(mf.column) + ","
	str += strconv.Itoa(mf.totalMines) + "\n"
	for y := 0; y < mf.column; y++ {
		for x := 0; x < mf.row; x++ {
			str += mf.field[mf.GetIdx(x, y)].String()
		}
		str += "\n"
	}
	return str
}
