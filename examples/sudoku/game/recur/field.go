package recur

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/sudoku/game"
)

func LoadSudokuField(dim int, diff game.Difficult) (result []int) {
	f := NewField(dim)
	f.Shuffle(f.Size() * f.Size())
	moves := diff.Percent(f.Dim())
	fmt.Printf("Для сложности %v ходов %v\n", diff, moves)
	for moves > 0 {
		x, y := rand.Intn(f.Size()), rand.Intn(f.Size())
		cell := f[f.Idx(x, y)]
		if cell.value > 0 {
			f.ResetCell(x, y)
			moves--
			fmt.Printf("Ход:%v xy[%v,%v]%v\n", moves, x, y, cell)
		}
	}
	for _, v := range f {
		result = append(result, v.value)
	}
	fmt.Println("Сгенерировано поле:", f)
	return result
}

type SudokuField []*Cell

func NewField(dim int) SudokuField {
	size := dim * dim
	f := make(SudokuField, 0)
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			f = append(f, NewCell(dim))
		}
	}
	return f
}

func (f SudokuField) ResetCell(x, y int) {
	idx := f.Idx(x, y)
	f[idx].Reset()
	fmt.Println("Обнулить ход:", idx, x, y, f[f.Idx(x, y)].notes)
	for x0 := 0; x0 < f.Size(); x0++ {
		n0 := f[f.Idx(x0, y)].value
		f[f.Idx(x, y)].RemoveNote(n0)
		// fmt.Println("Обнулить ход gorz:", idx, x0, y, n0, f[f.Idx(x0, y)].notes)
	}
	for y0 := 0; y0 < f.Size(); y0++ {
		n0 := f[f.Idx(x, y0)].value
		f[f.Idx(x, y)].RemoveNote(n0)
		// fmt.Println("Обнулить ход vert:", idx, x, y0, n0, f[f.Idx(x, y0)].notes)
	}
	rX0, rY0 := f.getRectIdx(x, y)
	for i := range f {
		x0, y0 := f.Pos(i)
		rX, rY := f.getRectIdx(x0, y0)
		if rX0 != rX || rY0 != rY {
			continue
		}
		n0 := f[f.Idx(x0, y0)].value
		f[f.Idx(x, y)].RemoveNote(n0)
		// fmt.Println("Обнулить ход rect:", idx, x, y, n0, f[f.Idx(x, y)].notes)
	}
	cell := f[f.Idx(x, y)]
	fmt.Println("Обнуление хода:", idx, x, y, cell.value, cell.notes, f.String())
}

func (d SudokuField) nextNote(x, y int) (note int) {
	rand.Seed(time.Now().UnixNano())
	idx := d.Idx(x, y)
	cell := d[idx]
	notes := cell.GetNotes()
	if len(notes) > 0 {
		note = notes[rand.Intn(len(notes))]
		best, result := d.isValidNote(x, y, note)
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

func (d SudokuField) isValidNote(x, y, note int) ([]int, bool) {
	var (
		notesCount map[int]int
		bestValues []int
		count, min int
		result     bool
	)
	notesCount = make(map[int]int)
	for x0 := d.Size() - 1; x0 >= 0; x0-- {
		cell := d[d.Idx(x0, y)]
		if cell.value > 0 {
			continue
		}
		notes := cell.GetNotes()
		for _, v := range notes {
			notesCount[v]++
		}
	}
	min = d.Size()
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

func (d SudokuField) Move(x, y int) {
	cell := d[d.Idx(x, y)]
	if cell.value > 0 {
		return
	}
	note := d.nextNote(x, y)
	cell.Add(note)
	d.updateNotesAfterMove(x, y, note)
	fmt.Printf("Move:%v [%v,%v]%v %v\n", d.Idx(x, y), x, y, note, d.checkNotesLen())
	d.checkDual(y)
}

func (d SudokuField) checkDual(y0 int) {
	if d.Dim() <= 2 {
		return
	}
	for x0 := 0; x0 < d.Size(); x0++ {
		cellA := d[d.Idx(x0, y0)]
		for x1 := 1; x1 < d.Size(); x1++ {
			cellB := d[d.Idx(x1, y0)]
			if cellA.value > 0 || cellB.value > 0 || len(cellA.notes) != 2 {
				continue
			}
			a, b := cellA.GetNotes(), cellB.GetNotes()
			if eui.IntSlicesIsEqual(a, b) {
				notes := cellB.GetNotes()
				note := notes[rand.Intn(len(notes))]
				cellB.Add(note)
				d.updateNotesAfterMove(x1, y0, note)
				fmt.Printf("___found dual [%v,%v]%v %v{%v,%v} %v", x1, y0, note, notes, cellA, cellB, d)
				break
			}
		}
	}
}

func (d SudokuField) Shuffle(idx int) {
	if idx == 0 {
		return
	}
	x, y := d.Pos(idx - 1)
	d.Move(x, y)
	d.Shuffle(idx - 1)
}

// Если на поле есть заметка с одним числом её открыть сразу после хода
func (d SudokuField) move(idx int) {
	if idx == 0 {
		return
	}
	x, y := d.Pos(idx - 1)
	cell := d[d.Idx(x, y)]
	if !(cell.value > 0) {
		if len(cell.GetNotes()) == 1 {
			note := cell.GetNotes()[0]
			cell.Add(note)
			d.updateNotesAfterMove(x, y, note)
			fmt.Printf("%v:Recur move:[%v,%v]%v\n%v\n", idx, x, y, note, d)
			d.move(d.Size() * d.Size())
		}
	}
	d.move(idx - 1)
}

func (d SudokuField) checkNotesLen() map[int]int {
	lens := make(map[int]int)
	for y := 0; y < d.Size(); y++ {
		for x := 0; x < d.Size(); x++ {
			cell := d[d.Idx(x, y)]
			if !(cell.value > 0) {
				l := len(cell.GetNotes())
				lens[l]++
			}
		}
	}
	return lens
}

func (d SudokuField) updateNotesAfterMove(x0, y0, value int) {
	for x := 0; x < d.Size(); x++ {
		idx := d.Idx(x, y0)
		cell := d[idx]
		cell.RemoveNote(value)
		d.IsValidNotesLen(x, y0, value, cell)
	}
	for y := 0; y < d.Size(); y++ {
		idx := d.Idx(x0, y)
		cell := d[idx]
		cell.RemoveNote(value)
		d.IsValidNotesLen(x0, y, value, cell)
	}
	rX0, rY0 := d.getRectIdx(x0, y0)
	for i, cell := range d {
		x, y := d.Pos(i)
		rX, rY := d.getRectIdx(x, y)
		if rX0 != rX || rY0 != rY {
			continue
		} else {
			cell.RemoveNote(value)
			d.IsValidNotesLen(x, y, value, cell)
		}
	}
	d.move(d.Size() * d.Size())
}

func (f SudokuField) IsValidNotesLen(x, y, value int, cell *Cell) {
	if len(cell.GetNotes()) == 0 && !(cell.value > 0) {
		fmt.Printf("panic notes zero len:[%v,%v]%v c:{%v} %v", x, y, value, cell, f)
		f.resetRow(y)
	}
}

// Обнулить эту строку
func (f SudokuField) resetRow(y0 int) {
	for x := 0; x < f.Size(); x++ {
		f[f.Idx(x, y0)].Reset()
	}
}

func (f SudokuField) getRectIdx(x int, y int) (rX int, rY int) {
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

func (d SudokuField) GetField() []*Cell      { return d }
func (d SudokuField) Dim() int               { return d[0].Dim() }
func (d SudokuField) Size() int              { return d[0].Size() }
func (d SudokuField) Idx(x, y int) int       { return y*d.Size() + x }
func (d SudokuField) Pos(idx int) (int, int) { return idx % d.Size(), idx / d.Size() }

func (d SudokuField) String() string {
	s := fmt.Sprintf("Sudoku %vX%v", d.Size(), d.Size())
	for y := 0; y < d.Size(); y++ {
		s += fmt.Sprintf("\n%v", y)
		for x := 0; x < d.Size(); x++ {
			s += fmt.Sprintf("[ %v:%v ]", x, d[d.Idx(x, y)].String())
		}
	}
	s += "\n"
	return s
}
