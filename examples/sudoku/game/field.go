package game

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/t0l1k/eui"
)

type Field struct {
	dim, size int
	cells     []*Cell
}

func NewField(dim int) *Field {
	f := &Field{dim: dim, size: dim * dim}
	for y := 0; y < f.size; y++ {
		for x := 0; x < f.size; x++ {
			f.cells = append(f.cells, NewCell(f.dim, eui.NewPointInt(x, y)))
		}
	}
	return f
}

func (f *Field) Reset() {
	for y := 0; y < f.size; y++ {
		for x := 0; x < f.size; x++ {
			f.cells[f.Idx(x, y)].Reset()
		}
	}
}

func (f *Field) New() {
	rand.Seed(time.Now().UnixNano())
	for y := 0; y < f.size; y++ {
		for x := 0; x < f.size; x++ {
			f.nextRand(x, y)
		}
	}
}

func (f *Field) nextRand(x, y int) {
	var (
		rValues, wValues []int
	)
	idx := f.Idx(x, y)
	rValues = f.cells[idx].GetNotes()
	fmt.Println(idx, "00 rand", rValues, len(rValues), x, y)
	if len(rValues) > 0 {
		rnd := rand.Intn(len(rValues))
		n := rValues[rnd]
		f.Add(x, y, rValues[rnd])
		if f.IsValidPath(x, y) {
			fmt.Println(idx, "01 rand is valid path:", x, y, rnd, n, rValues[rnd], rValues, f.StringFull())
			return
		} else {
			wValues = append(wValues, n)
			fmt.Println(idx, "02 !rand remove n:", x, y, n, wValues, f.StringFull())
		}
	}
	if len(rValues) == 0 {
		fmt.Println(idx, "03 rand len 0 update values:", x, y, rValues, f.StringFull())
		rValues = f.cells[idx].GetNotes()
	}
	fmt.Println(idx, "04 rand found zero len notes", rValues, wValues, idx, x, y)
	for i := 0; i < len(rValues); i++ {
		if ArrIsContain(wValues, rValues[i]) {
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

func (f *Field) Add(x, y, n int) {
	idx := f.Idx(x, y)
	f.cells[idx].Add(n)
	fmt.Println("Сделан ход:", n, idx, x, y, f.cells[f.Idx(x, y)].notes.values)
	for x0 := 0; x0 < f.size; x0++ {
		f.cells[f.Idx(x0, y)].notes.RemoveNote(n)
		fmt.Println("Сделан ход gorz:", n, idx, x, y, x0, f.cells[f.Idx(x0, y)].notes.values)
	}
	for y0 := 0; y0 < f.size; y0++ {
		f.cells[f.Idx(x, y0)].notes.RemoveNote(n)
		fmt.Println("Сделан ход vert:", n, idx, x, y, y0, f.cells[f.Idx(x, y0)].notes.values)
	}

	rX0, rY0 := f.getRectIdx(x, y)
	for i, v := range f.cells {
		x1, y1 := f.Pos(i)
		rX, rY := f.getRectIdx(x1, y1)
		if rX0 != rX || rY0 != rY {
			continue
		}
		v.notes.RemoveNote(n)
		fmt.Println("Сделан ход rect:", n, idx, x, y, x1, y1, f.cells[f.Idx(x1, y1)].notes.values)
	}
	fmt.Println("Результат хода:", n, idx, x, y, f.cells[f.Idx(x, y)].notes.values, f.StringFull())
}

func (f *Field) ResetCell(x, y int) {
	idx := f.Idx(x, y)
	f.cells[idx].Reset()
	fmt.Println("Обнулить ход:", idx, x, y, f.cells[f.Idx(x, y)].notes.values)
	for x0 := 0; x0 < f.size; x0++ {
		n0 := f.cells[f.Idx(x0, y)].Value().(int)
		f.cells[f.Idx(x, y)].notes.RemoveNote(n0)
		fmt.Println("Обнулить ход gorz:", idx, x0, y, n0, f.cells[f.Idx(x0, y)].notes.values)
	}
	for y0 := 0; y0 < f.size; y0++ {
		n0 := f.cells[f.Idx(x, y0)].Value().(int)
		f.cells[f.Idx(x, y)].notes.RemoveNote(n0)
		fmt.Println("Обнулить ход vert:", idx, x, y0, n0, f.cells[f.Idx(x, y0)].notes.values)
	}
	rX0, rY0 := f.getRectIdx(x, y)
	for i := range f.cells {
		x0, y0 := f.Pos(i)
		rX, rY := f.getRectIdx(x0, y0)
		if rX0 != rX || rY0 != rY {
			continue
		}
		n0 := f.cells[f.Idx(x0, y0)].Value().(int)
		f.cells[f.Idx(x, y)].notes.RemoveNote(n0)
		fmt.Println("Обнулить ход rect:", idx, x, y, n0, f.cells[f.Idx(x, y)].notes.values)
	}
	cell := f.cells[f.Idx(x, y)]
	fmt.Println("Обнуление хода:", idx, x, y, cell.Value(), cell.notes.values, f.StringFull())
}

func (f *Field) GetCells() []*Cell     { return f.cells }
func (f Field) Idx(x, y int) int       { return y*f.size + x }
func (f Field) Pos(idx int) (int, int) { return idx % f.size, idx / f.size } //x,y

func (f *Field) IsValidPath(x0, y0 int) (result bool) {
	result = true
	idx := f.Idx(x0, y0)
	dim := f.size * f.size
	if idx <= dim-1 {
		idx++
	}
	for i := idx; i < dim; i++ {
		x, y := f.Pos(i)
		if f.isEmptyNotes(x, y) {
			if f.cells[f.Idx(x0, y0)].Value().(int) > 0 {
				f.ResetCell(x0, y0)
				fmt.Println(idx, "Обнулить 1 вызвавшей появление пустой:", i, len(f.cells[i].notes.values), f.cells[i].notes.String(), f)
			}
			fmt.Println(idx, "Есть пустые заметки:", i, len(f.cells[i].notes.values), f.cells[i].notes.String())
			result = false
		}
		if !result {
			f.ResetCell(f.Pos(i))
			fmt.Println(idx, "Заметка обновлена:", i, len(f.cells[i].notes.values), f.cells[i].notes.String(), f.StringFull())
		}
	}
	return result
}

func (f Field) isEmptyNotes(x, y int) bool {
	idx := f.Idx(x, y)
	return len(f.cells[idx].notes.values) == 0
}

func (f Field) IsValidMove(x, y, value int) bool {
	return f.isValidHor(y, value) && f.isValidVer(x, value) && f.isValidRect(x, y, value)
}

func (f Field) isValidRect(x0, y0, value int) bool {
	rX0, rY0 := f.getRectIdx(x0, y0)
	var arr []int
	for i, v := range f.cells {
		x, y := f.Pos(i)
		rX, rY := f.getRectIdx(x, y)
		if rX0 != rX || rY0 != rY {
			continue
		}
		arr = append(arr, v.Value().(int))
	}
	for i := 0; i < len(arr); i++ {
		if arr[i] == 0 {
			continue
		}
		if arr[i] == value {
			return false
		}
	}
	return true
}

func (f Field) getRectIdx(x int, y int) (rX int, rY int) {
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

func (f Field) isValidHor(y0, value int) bool {
	for x0 := 0; x0 < f.size; x0++ {
		idx := f.Idx(x0, y0)
		if f.cells[idx].Value() == value {
			return false
		}
	}
	return true
}
func (f Field) isValidVer(x, value int) bool {
	for y0 := 0; y0 < f.size; y0++ {
		idx := f.Idx(x, y0)
		if f.cells[idx].Value() == value {
			return false
		}
	}
	return true
}

func (f Field) String() (result string) {
	result = fmt.Sprintf("sudoku %vX%v\n", f.size, f.size)
	for y := 0; y < f.size; y++ {
		for x := 0; x < f.size; x++ {
			result += f.cells[f.Idx(x, y)].StringValueShort()
		}
		result += "\n"
	}
	return result
}

func (f Field) StringFull() (result string) {
	result = fmt.Sprintf("sudoku %vX%v\n", f.size, f.size)
	for y := 0; y < f.size; y++ {
		for x := 0; x < f.size; x++ {
			cell := f.cells[f.Idx(x, y)]
			result += cell.StringValueShort()
		}
		result += "\n"
	}
	return result
}
