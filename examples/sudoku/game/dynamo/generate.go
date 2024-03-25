package dynamo

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/sudoku/game"
)

type GenSudokuField struct {
	dim, size int
	diff      game.Difficult
	cells     []*cell
}

func NewGenSudokuField(d int, diff game.Difficult) *GenSudokuField {
	f := &GenSudokuField{dim: d, size: d * d, diff: diff}
	f.New()
	return f
}

func (f *GenSudokuField) New() {
	for y := 0; y < f.size; y++ {
		for x := 0; x < f.size; x++ {
			f.cells = append(f.cells, newCell(f.dim, f.size))
		}
	}
}

func (f *GenSudokuField) GetField() (result []int) {
	f.shuffle()
	moves := f.diff.Percent(f.dim)
	fmt.Printf("Для сложности %v ходов %v\n", f.diff, moves)
	for moves > 0 {
		x, y := rand.Intn(f.size), rand.Intn(f.size)
		cell := f.cells[f.Idx(x, y)]
		if cell.value > 0 {
			f.ResetCell(x, y)
			moves--
			fmt.Printf("Ход:%v xy[%v,%v]%v\n", moves, x, y, cell)
		}
	}
	for _, v := range f.cells {
		result = append(result, v.value)
	}
	fmt.Println("Сгенерировано поле:", f)
	return result
}

func (f *GenSudokuField) getCells() []*cell     { return f.cells }
func (f GenSudokuField) Idx(x, y int) int       { return y*f.size + x }
func (f GenSudokuField) Pos(idx int) (int, int) { return idx % f.size, idx / f.size } //x,y

func (f *GenSudokuField) shuffle() {
	rand.Seed(time.Now().UnixNano())
	for y := 0; y < f.size; y++ {
		for x := 0; x < f.size; x++ {
			f.nextRand(x, y)
		}
	}
}

func (f *GenSudokuField) nextRand(x, y int) {
	var (
		rValues, wValues []int
	)
	idx := f.Idx(x, y)
	rValues = f.cells[idx].getNotes()
	fmt.Println(idx, "00 rand", rValues, len(rValues), x, y)
	if len(rValues) > 0 {
		rnd := rand.Intn(len(rValues))
		n := rValues[rnd]
		f.Add(x, y, rValues[rnd])
		if f.IsValidPath(x, y) {
			fmt.Println(idx, "01 rand is valid path:", x, y, rnd, n, rValues[rnd], rValues, f.String())
			return
		} else {
			wValues = append(wValues, n)
			fmt.Println(idx, "02 !rand remove n:", x, y, n, wValues, f.String())
		}
	}
	if len(rValues) == 0 {
		fmt.Println(idx, "03 rand len 0 update values:", x, y, rValues, f.String())
		rValues = f.cells[idx].getNotes()
	}
	fmt.Println(idx, "04 rand found zero len notes", rValues, wValues, idx, x, y)
	for i := 0; i < len(rValues); i++ {
		if eui.IntSliceContains(wValues, rValues[i]) {
			fmt.Println(idx, "05 rand skip used note", rValues, wValues, idx, x, y)
			continue
		}
		f.Add(x, y, rValues[i])
		if f.IsValidPath(x, y) {
			f.Add(x, y, rValues[i])
			fmt.Println(idx, "06 rand len 0 found path", i, rValues[i], rValues, len(rValues), x, y)
			return
		}
		wValues = append(wValues, rValues[i])
		fmt.Println(idx, "07 rand len 0 repeat find move:", i, rValues[i], rValues, wValues, len(rValues), x, y)
	}
}

func (f *GenSudokuField) Add(x, y, n int) {
	idx := f.Idx(x, y)
	f.cells[idx].add(n)
	fmt.Println("Сделан ход:", n, idx, x, y, f.cells[f.Idx(x, y)].notes)
	for x0 := 0; x0 < f.size; x0++ {
		f.cells[f.Idx(x0, y)].removeNote(n)
		fmt.Println("Сделан ход gorz:", n, idx, x, y, x0, f.cells[f.Idx(x0, y)].notes)
	}
	for y0 := 0; y0 < f.size; y0++ {
		f.cells[f.Idx(x, y0)].removeNote(n)
		fmt.Println("Сделан ход vert:", n, idx, x, y, y0, f.cells[f.Idx(x, y0)].notes)
	}

	rX0, rY0 := f.getRectIdx(x, y)
	for i, v := range f.cells {
		x1, y1 := f.Pos(i)
		rX, rY := f.getRectIdx(x1, y1)
		if rX0 != rX || rY0 != rY {
			continue
		}
		v.removeNote(n)
		fmt.Println("Сделан ход rect:", n, idx, x, y, x1, y1, f.cells[f.Idx(x1, y1)].notes)
	}
	fmt.Println("Результат хода:", n, idx, x, y, f.cells[f.Idx(x, y)].notes, f.String())
}

func (f *GenSudokuField) ResetCell(x, y int) {
	idx := f.Idx(x, y)
	f.cells[idx].reset()
	fmt.Println("Обнулить ход:", idx, x, y, f.cells[f.Idx(x, y)].notes)
	for x0 := 0; x0 < f.size; x0++ {
		n0 := f.cells[f.Idx(x0, y)].value
		f.cells[f.Idx(x, y)].removeNote(n0)
		fmt.Println("Обнулить ход gorz:", idx, x0, y, n0, f.cells[f.Idx(x0, y)].notes)
	}
	for y0 := 0; y0 < f.size; y0++ {
		n0 := f.cells[f.Idx(x, y0)].value
		f.cells[f.Idx(x, y)].removeNote(n0)
		fmt.Println("Обнулить ход vert:", idx, x, y0, n0, f.cells[f.Idx(x, y0)].notes)
	}
	rX0, rY0 := f.getRectIdx(x, y)
	for i := range f.cells {
		x0, y0 := f.Pos(i)
		rX, rY := f.getRectIdx(x0, y0)
		if rX0 != rX || rY0 != rY {
			continue
		}
		n0 := f.cells[f.Idx(x0, y0)].value
		f.cells[f.Idx(x, y)].removeNote(n0)
		fmt.Println("Обнулить ход rect:", idx, x, y, n0, f.cells[f.Idx(x, y)].notes)
	}
	cell := f.cells[f.Idx(x, y)]
	fmt.Println("Обнуление хода:", idx, x, y, cell.value, cell.notes, f.String())
}

func (f *GenSudokuField) IsValidPath(x0, y0 int) (result bool) {
	result = true
	idx := f.Idx(x0, y0)
	dim := f.size * f.size
	if idx <= dim-1 {
		idx++
	}
	for i := idx; i < dim; i++ {
		x, y := f.Pos(i)
		if f.isEmptyNotes(x, y) {
			if f.cells[f.Idx(x0, y0)].value > 0 {
				f.ResetCell(x0, y0)
				fmt.Println(idx, "Обнулить 1 вызвавшей появление пустой:", i, len(f.cells[i].notes), f.cells[i].String(), f)
			}
			fmt.Println(idx, "Есть пустые заметки:", i, len(f.cells[i].notes), f.cells[i].String())
			result = false
		}
		if !result {
			f.ResetCell(f.Pos(i))
			fmt.Println(idx, "Заметка обновлена:", i, len(f.cells[i].notes), f.cells[i].String(), f.String())
		}
	}
	return result
}

func (f GenSudokuField) getRectIdx(x int, y int) (rX int, rY int) {
	szX := f.size
	szY := f.size
	rX = f.dim
	rY = f.dim
	for szX > x {
		szX -= f.dim
		rX--
	}
	for szY > y {
		szY -= f.dim
		rY--
	}
	return rX, rY
}

func (f GenSudokuField) isEmptyNotes(x, y int) bool {
	idx := f.Idx(x, y)
	return len(f.cells[idx].notes) == 0
}

func (f GenSudokuField) String() (result string) {
	result = fmt.Sprintf("sudoku %vX%v\n", f.size, f.size)
	for y := 0; y < f.size; y++ {
		for x := 0; x < f.size; x++ {
			cell := f.cells[f.Idx(x, y)]
			result += cell.String()
		}
		result += "\n"
	}
	return result
}
