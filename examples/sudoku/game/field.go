package game

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/t0l1k/eui"
)

type Field []*Cell

func newField(dim int) *Field {
	size := dim * dim
	f := make(Field, 0)
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			f = append(f, NewCell(dim))
		}
	}
	return &f
}

func (f Field) Reset() {
	for i := range f {
		f[i].Reset()
	}
}

func (f Field) Shuffle(idx int) {
	if idx == 0 {
		return
	}
	x, y := f.Pos(idx - 1)
	f.Move(x, y)
	f.Shuffle(idx - 1)
}

func (f Field) prepareFor(diff Difficult) {
	moves := diff.Percent(f.Dim())
	fmt.Printf("Для сложности %v ходов %v\n", diff, moves)
	for moves > 0 {
		x, y := rand.Intn(f.Size()), rand.Intn(f.Size())
		cell := f[f.Idx(x, y)]
		if cell.GetValue() > 0 {
			f.ResetCell(x, y)
			moves--
			fmt.Printf("Ход:%v xy[%v,%v]%v\n", moves, x, y, cell)
		}
	}
	for _, v := range f {
		if v.GetValue() > 0 {
			v.SetReadOnly()
		}
	}
}

func (f Field) ResetCell(x, y int) {
	idx := f.Idx(x, y)
	f[idx].Reset()
	fmt.Println("Обнулить ход:", idx, x, y, f[f.Idx(x, y)].notes)
	for x0 := 0; x0 < f.Size(); x0++ {
		n0 := f[f.Idx(x0, y)].GetValue()
		f[f.Idx(x, y)].UpdateNote(n0)
	}
	for y0 := 0; y0 < f.Size(); y0++ {
		n0 := f[f.Idx(x, y0)].GetValue()
		f[f.Idx(x, y)].UpdateNote(n0)
	}
	rX0, rY0 := f.getRectIdx(x, y)
	for i := range f {
		x0, y0 := f.Pos(i)
		rX, rY := f.getRectIdx(x0, y0)
		if rX0 != rX || rY0 != rY {
			continue
		}
		n0 := f[f.Idx(x0, y0)].GetValue()
		f[f.Idx(x, y)].UpdateNote(n0)
	}
	cell := f[f.Idx(x, y)]
	fmt.Println("Обнуление хода:", idx, x, y, cell.GetValue(), cell.notes, f.String())
}

func (f Field) Move(x, y int) {
	cell := f[f.Idx(x, y)]
	if cell.GetValue() > 0 {
		return
	}
	note := f.nextNote(x, y)
	cell.Add(note)
	f.updateNotesAfterMove(x, y, note)
	fmt.Printf("Move:%v [%v,%v]%v\n", f.Idx(x, y), x, y, note)
	f.checkDual(y)
}

func (f Field) nextNote(x, y int) (note int) {
	rand.Seed(time.Now().UnixNano())
	idx := f.Idx(x, y)
	cell := f[idx]
	notes := cell.GetNotes()
	if len(notes) > 0 {
		note = notes[rand.Intn(len(notes))]
		best, result := f.isValidNote(x, y, note)
		if !result {
			for _, v := range notes {
				if eui.IntSliceContains(best, v) {
					note = v
					break
				}
			}
			fmt.Println("new nextNote", idx, note, notes)
		}
	}
	return note
}

func (f Field) isValidNote(x, y, note int) ([]int, bool) {
	var (
		notesCount map[int]int
		bestValues []int
		count, min int
		result     bool
	)
	notesCount = make(map[int]int)
	for x0 := f.Size() - 1; x0 >= 0; x0-- {
		cell := f[f.Idx(x0, y)]
		if cell.GetValue() > 0 {
			continue
		}
		notes := cell.GetNotes()
		for _, v := range notes {
			notesCount[v]++
		}
	}
	min = f.Size()
	for _, v := range notesCount {
		if v < min {
			min = v
		}
	}
	for k, v := range notesCount {
		if v > min {
			continue
		}
		count++
		bestValues = append(bestValues, k)
	}
	result = eui.IntSliceContains(bestValues, note)
	fmt.Printf("Результат проверки заметок:[%v,%v]%v nc:{%v} best:%v count:%v {%v}", x, y, note, notesCount, bestValues, count, result)
	return bestValues, result
}

func (f Field) checkDual(y0 int) {
	if f.Dim() < 3 {
		return
	}
	for x0 := 0; x0 < f.Size(); x0++ {
		cellA := f[f.Idx(x0, y0)]
		for x1 := 1; x1 < f.Size(); x1++ {
			cellB := f[f.Idx(x1, y0)]
			if cellA.GetValue() > 0 || cellB.GetValue() > 0 || len(cellA.notes) != 2 {
				continue
			}
			a, b := cellA.GetNotes(), cellB.GetNotes()
			if eui.IntSlicesIsEqual(a, b) {
				notes := cellB.GetNotes()
				note := notes[rand.Intn(len(notes))]
				cellB.Add(note)
				f.updateNotesAfterMove(x1, y0, note)
				fmt.Printf("___found dual [%v,%v]%v %v{%v,%v} %v", x1, y0, note, notes, cellA, cellB, f)
				break
			}
		}
	}
}

func (f Field) openSingleNote(idx int) {
	if idx == 0 {
		return
	}
	x, y := f.Pos(idx - 1)
	cell := f[f.Idx(x, y)]
	if !(cell.GetValue() > 0) {
		if len(cell.GetNotes()) == 1 {
			note := cell.GetNotes()[0]
			cell.Add(note)
			f.updateNotesAfterMove(x, y, note)
			fmt.Printf("%v:Recur move:[%v,%v]%v\n%v\n", idx, x, y, note, f)
			f.openSingleNote(f.Size() * f.Size())
		}
	}
	f.openSingleNote(idx - 1)
}

func (f Field) updateNotesAfterMove(x0, y0, value int) {
	for x := 0; x < f.Size(); x++ {
		idx := f.Idx(x, y0)
		cell := f[idx]
		cell.UpdateNote(value)
		f.IsValidNotesLen(x, y0, value, cell)
	}
	for y := 0; y < f.Size(); y++ {
		idx := f.Idx(x0, y)
		cell := f[idx]
		cell.UpdateNote(value)
		f.IsValidNotesLen(x0, y, value, cell)
	}
	rX0, rY0 := f.getRectIdx(x0, y0)
	for i, cell := range f {
		x, y := f.Pos(i)
		rX, rY := f.getRectIdx(x, y)
		if rX0 != rX || rY0 != rY {
			continue
		} else {
			cell.UpdateNote(value)
			f.IsValidNotesLen(x, y, value, cell)
		}
	}
	f.openSingleNote(f.Size() * f.Size())
}

func (f Field) IsValidNotesLen(x, y, value int, cell *Cell) {
	if len(cell.GetNotes()) == 0 && !(cell.GetValue() > 0) {
		fmt.Printf("panic notes zero len:[%v,%v]%v c:{%v} %v", x, y, value, cell, f)
		f.resetRow(y)
	}
}

// Обнулить эту строку
func (f Field) resetRow(y0 int) {
	for x := 0; x < f.Size(); x++ {
		f[f.Idx(x, y0)].Reset()
	}
	f.Shuffle(f.Size() * f.Size())
}

func (f Field) getRectIdx(x int, y int) (rX int, rY int) {
	szX := f.Size()
	szY := f.Size()
	rX = f.Dim()
	rY = f.Dim()
	for szX > x {
		szX -= f.Dim()
		rX--
	}
	for szY > y {
		szY -= f.Dim()
		rY--
	}
	return rX, rY
}

func (f Field) GetField() []*Cell      { return f }
func (f Field) Dim() int               { return f[0].Dim() }
func (f Field) Size() int              { return f[0].Size() }
func (f Field) Idx(x, y int) int       { return y*f.Size() + x }
func (f Field) Pos(idx int) (int, int) { return idx % f.Size(), idx / f.Size() }

func (f Field) String() string {
	s := fmt.Sprintf("Sudoku %vX%v", f.Size(), f.Size())
	for y := 0; y < f.Size(); y++ {
		s += fmt.Sprintf("\n%v", y)
		for x := 0; x < f.Size(); x++ {
			s += fmt.Sprintf("[ %v:%v ]", x, f[f.Idx(x, y)].String())
		}
	}
	s += "\n"
	return s
}
