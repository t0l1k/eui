package game

import (
	"errors"
	"fmt"
	"log"
	"math/rand"

	"github.com/t0l1k/eui/utils"
)

const (
	err00 = "найдена пустая клетка"
	err01 = "вариант 01 появления пустой клетки"
	err02 = "вариант 02 появления пустой клетки"
	err03 = "вариант 03 исчерпали заметки"
)

type Game struct {
	dim     Dim
	diff    Difficult
	Percent int
	field   *field
	path    map[int]utils.IntList
	history utils.IntList
	inGame  bool
}

func NewGame(dim Dim) *Game {
	g := &Game{dim: dim, field: newField(dim.Size())}
	return g
}

func (g Game) Dim() Dim               { return g.dim }
func (g Game) Size() int              { return g.dim.Size() }
func (g Game) Cell(x, y int) *Cell    { return g.field.cell(g.Idx(x, y)) }
func (g Game) Idx(x, y int) int       { return y*g.Size() + x }
func (g Game) Pos(idx int) (int, int) { return idx % g.Size(), idx / g.Size() }

func (g *Game) Load(diff Difficult) {
	switch diff {
	case Easy, Normal, Hard, Extreme:
		g.diff = diff
		g.shuffle()
		g.Percent = g.field.prepareFor(g.diff, g.dim.Size())
		g.UpdateAllFieldNotes()
		g.inGame = true
	case Manual:
		g.field.reset(g.Size())
		g.path = make(map[int]utils.IntList)
	}
}

func (g *Game) GetPercent() int {
	moves := g.field.GetMovesCount()
	sz := g.dim.Size()
	g.Percent = moves * 100 / (sz * sz)
	return g.Percent
}

func (g *Game) MarkReadOnly() {
	g.field.MarkReadOnly()
	g.inGame = true
}

func (g *Game) shuffle() {
	g.field.reset(g.Size())
	idx := 0
	count := 0
	g.path = make(map[int]utils.IntList)
	for g.field.isFoundEmptyCells() {
		if err := g.guess(idx); err == nil {
			idx++
			if idx >= g.Size()*g.Size() {
				idx = 0
			}
			fmt.Printf("%v,%v сделан ход на поле\n%v", count, idx, g.String())
		} else {
			switch err.Error() {
			case err03:
				for k := range g.path {
					if k >= idx {
						delete(g.path, idx)
					}
				}
				idx--
			}
			x0, y0 := g.Pos(idx)
			for y := y0; y < g.Size(); y++ {
				for x := x0; x < g.Size(); x++ {
					g.field.cell(g.Idx(x, y)).reset(g.Size())
				}
			}
			g.UpdateAllFieldNotes()
			fmt.Printf("%v,%v на поле есть пустые заметки\n%v\n%v\n", count, idx, g.String(), err)
		}
		count++
	}
}

func (g *Game) guess(idx int) error {
	if g.Cell(g.Pos(idx)).GetValue() > 0 {
		return nil
	}
	cell := g.field.cell(idx)
	notes := cell.GetNotes()
	for _, v := range cell.GetNotes() {
		if g.path[idx].IsContain(v) {
			notes, _ = notes.RemoveValue(v)
			fmt.Printf("Удаляем метку:%v из заметок:[%v] рузультат пути:[%v]", v, notes, g.path)
		}
	}
	if len(notes) == 0 {
		return errors.New(err03)
	}
	note := g.getRndNote(notes)
	x, y := g.Pos(idx)
	g.MakeMove(x, y, note)
	fmt.Printf("%v ход %v ячейка[%v] успешен на поле\n%v\n", idx, note, cell, g.String())
	if err := g.isFoundEmptyNotes(idx); err != nil {
		fmt.Printf("%v ход %v ячейка[%v]есть пустая клетка на поле\n%v\n%v\n", idx, note, cell, g.String(), err)
		return err
	}
	return nil
}

func (*Game) getRndNote(notes []int) int {
	note := notes[rand.Intn(len(notes))]
	return note
}

func (g *Game) isFoundEmptyNotes(idx int) error {
	notes := make(map[int]bool)
	for i := 1; i <= g.Size(); i++ {
		notes[i] = false
	}
	if g.field.isFoundEmptyNote() {
		fmt.Println(notes)
		return errors.New(err00)
	}
	var allNotes []utils.IntList
	_, y0 := g.Pos(idx)
	for x := 0; x < g.Size(); x++ {
		value := g.Cell(x, y0).GetValue()
		if value > 0 {
			notes[value] = true
		} else {
			cell := g.Cell(x, y0)
			allNotes = append(allNotes, cell.GetNotes())
			for _, v := range cell.GetNotes() {
				notes[v] = true
			}
		}
	}
	for _, v := range notes {
		if !v {
			fmt.Println(notes, allNotes)
			return errors.New(err01)
		}
	}
	for i, v1 := range allNotes {
		if len(v1) > 1 {
			continue
		}
		for j, v2 := range allNotes {
			if len(v2) > 1 || i == j {
				continue
			}
			if v1.Equals(v2) {
				fmt.Println(notes, allNotes)
				return errors.New(err02)
			}
		}
	}
	return nil
}

func (g *Game) MakeMove(x, y, note int) bool {
	idx := g.Idx(x, y)
	if !g.inGame {
		g.path[idx] = append(g.path[idx], note)
		fmt.Printf("Ход %v метка:%v путь %v\n", idx, note, g.path)
	} else {
		g.history = g.history.Add(idx)
	}
	cell := g.field.cell(idx)
	if !cell.add(note) {
		fmt.Println("move on read-only cell")
		return true
	}
	count := g.UpdateAllFieldNotes()
	fmt.Printf("Ход %v метка:%v обновленно меток:%v в истории ходов:%v\n", idx, note, count, len(g.history))
	return count > 0
}

func (g *Game) Undo() {
	if g.history.IsEmpty() {
		return
	}
	x, y := g.LastMovePos()
	g.ResetCell(x, y)
	g.history = g.history.Pop()
	log.Printf("Undo move[%v,%v]%v в истории ходов:%v\n", x, y, g.Cell(x, y), len(g.history))
}

func (g *Game) LastMovePos() (int, int) { return g.Pos(g.history[g.history.Size()-1]) }
func (g *Game) MovesCount() int         { return g.history.Size() }

func (g *Game) IsWin() bool {
	var checkHor, checkVert, checkRects bool
	var arr utils.IntList
	for x0 := 0; x0 < g.Size(); x0++ {
		arr = g.fillArr()
		for y := 0; y < g.Size(); y++ {
			value := g.Cell(x0, y).GetValue()
			if value == 0 {
				return false
			} else {
				arr, _ = arr.RemoveValue(value)
			}
		}
		if len(arr) == 0 {
			checkHor = true
		} else {
			return false
		}
	}

	for y0 := 0; y0 < g.Size(); y0++ {
		arr = g.fillArr()
		for x := 0; x < g.Size(); x++ {
			value := g.Cell(x, y0).GetValue()
			if value == 0 {
				return false
			} else {
				arr, _ = arr.RemoveValue(value)
			}
		}
		if len(arr) == 0 {
			checkVert = true
		} else {
			return false
		}
	}

	for rY0 := 0; rY0 < g.dim.W; rY0++ {
		for rX0 := 0; rX0 < g.dim.H; rX0++ {
			arr = g.fillArr()
			for i, v := range *g.field {
				x, y := g.Pos(i)
				rX, rY := g.getRectIdx(x, y)
				if rX0 != rX || rY0 != rY {
					continue
				}
				value := v.GetValue()
				if value == 0 {
					return false
				} else {
					arr, _ = arr.RemoveValue(value)
				}
			}
			if len(arr) == 0 {
				checkRects = true
			} else {
				return false
			}
		}
	}
	return checkHor && checkVert && checkRects
}

func (g *Game) fillArr() (arr []int) {
	for i := 1; i <= g.Size(); i++ {
		arr = append(arr, i)
	}
	return arr
}

func (g *Game) UpdateAllFieldNotes() (count int) {
	for i := range *g.field {
		x, y := g.Pos(i)
		cell := g.Cell(x, y)
		if cell.GetValue() > 0 {
			continue
		}
		g.ResetCell(x, y)
		count++
	}
	return count
}

func (g *Game) ResetCell(x0, y0 int) {
	cell := g.Cell(x0, y0)
	if cell.IsReadOnly() {
		return
	}
	cell.reset(g.Size())
	for y := 0; y < g.Size(); y++ {
		value := g.Cell(x0, y).GetValue()
		if value > 0 {
			g.Cell(x0, y0).setNote(value)
		}
	}
	for x := 0; x < g.Size(); x++ {
		value := g.Cell(x, y0).GetValue()
		if value > 0 {
			g.Cell(x0, y0).setNote(value)
		}
	}
	rX0, rY0 := g.getRectIdx(x0, y0)
	for i := range *g.field {
		x1, y1 := g.Pos(i)
		rX, rY := g.getRectIdx(x1, y1)
		if rX0 != rX || rY0 != rY {
			continue
		}
		value := g.Cell(x1, y1).GetValue()
		if value > 0 {
			g.Cell(x0, y0).setNote(value)
		}
	}
}

func (g *Game) ValuesCount() (counts map[int]int) {
	counts = make(map[int]int)
	for i := 1; i <= g.Size(); i++ {
		counts[i] = g.field.valueCount(i)
	}
	return counts
}

func (g Game) getRectIdx(x int, y int) (rX int, rY int) {
	szX := g.Size()
	szY := g.Size()
	rX = g.dim.W
	rY = g.dim.H
	if g.dim.W != g.dim.H {
		rX = g.dim.H
		rY = g.dim.W
	}
	for szX > x {
		szX -= g.dim.W
		rX--
	}
	for szY > y {
		szY -= g.dim.H
		rY--
	}
	return rX, rY
}

func (g Game) String() string {
	res := "Sudoku: "
	res += g.dim.String() + "\n"
	for y := 0; y < g.Size(); y++ {
		for x := 0; x < g.Size(); x++ {
			res += fmt.Sprintf("[ %2v ]", g.Cell(x, y).String())
		}
		res += "\n"
	}
	return res
}
