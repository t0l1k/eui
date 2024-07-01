package game

import (
	"fmt"
	"math/rand"
)

type field []*Cell

func newField(size int) *field {
	f := make(field, 0)
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			f = append(f, newCell())
		}
	}
	return &f
}

func (f field) reset(size int) {
	for _, v := range f {
		v.reset(size)
	}
}

func (f field) cell(idx int) *Cell { return f[idx] }

func (f field) isFoundEmptyCells() bool {
	for _, v := range f {
		if v.GetValue() == 0 {
			return true
		}
	}
	return false
}

func (f field) isFoundEmptyNote() bool {
	for _, v := range f {
		if v.GetValue() == 0 && len(v.GetNotes()) == 0 {
			return true
		}
	}
	return false
}

func (f field) prepareFor(diff Difficult, size int) (percent int) {
	moves, percent := diff.Percent(size)
	fmt.Printf("Для сложности %v ходов %v\n", diff, moves)
	for moves > 0 {
		x, y := rand.Intn(size), rand.Intn(size)
		idx := y*size + x
		cell := f[idx]
		if cell.GetValue() > 0 {
			f.cell(idx).reset(size)
			moves--
			fmt.Printf("Ход:%v xy[%v,%v]%v\n", moves, x, y, cell)
		}
	}
	for _, v := range f {
		if v.GetValue() > 0 {
			v.setReadOnly()
		}
	}
	return percent
}
