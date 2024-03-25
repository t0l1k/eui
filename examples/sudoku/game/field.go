package game

import (
	"fmt"

	"github.com/t0l1k/eui"
)

type Field struct {
	dim, size int
	cells     []*Cell
	generated []int
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

func (f *Field) Load(field []int) {
	f.generated = field
	f.Reset()
	for i := range f.cells {
		x, y := f.Pos(i)
		f.Add(x, y, field[i])
	}
}

func (f *Field) Reset() {
	for i := range f.cells {
		f.cells[i].Reset()
	}
}

func (f Field) Idx(x, y int) int       { return y*f.size + x }
func (f Field) Pos(idx int) (int, int) { return idx % f.size, idx / f.size } //x,y
func (f *Field) GetCells() []*Cell     { return f.cells }

func (f *Field) Add(x, y, n int) {
	idx := f.Idx(x, y)
	f.cells[idx].Add(n)
	fmt.Println("Сделан ход:", n, idx, x, y, f.cells[idx].notes)
	for x0 := 0; x0 < f.size; x0++ {
		f.cells[f.Idx(x0, y)].notes.RemoveNote(n)
		fmt.Println("Сделан ход gorz:", n, idx, x, y, x0, f.cells[f.Idx(x0, y)].notes)
	}
	for y0 := 0; y0 < f.size; y0++ {
		f.cells[f.Idx(x, y0)].notes.RemoveNote(n)
		fmt.Println("Сделан ход vert:", n, idx, x, y, y0, f.cells[f.Idx(x, y0)].notes)
	}

	rX0, rY0 := f.getRectIdx(x, y)
	for i, v := range f.cells {
		x1, y1 := f.Pos(i)
		rX, rY := f.getRectIdx(x1, y1)
		if rX0 != rX || rY0 != rY {
			continue
		}
		v.notes.RemoveNote(n)
		fmt.Println("Сделан ход rect:", n, idx, x, y, x1, y1, f.cells[f.Idx(x1, y1)].notes)
	}
	fmt.Println("Результат хода:", n, idx, x, y, f.cells[idx].notes, f.String())
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

func (f Field) String() (result string) {
	result = fmt.Sprintf("sudoku %vX%v\n", f.size, f.size)
	for y := 0; y < f.size; y++ {
		for x := 0; x < f.size; x++ {
			result += f.cells[y*f.size+x].StringValueShort()
		}
		result += "\n"
	}
	return result
}
