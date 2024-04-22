package game

import (
	"fmt"
	"log"

	"github.com/t0l1k/eui"
)

type Game struct {
	dim, size int
	diff      Difficult
	field     *Field
	history   []int
}

func NewGame(dim int) *Game {
	g := &Game{dim: dim, size: dim * dim}
	g.field = newField(dim)
	return g
}

func (g *Game) generate() {
	g.field.Shuffle(g.field.Size() * g.field.Size())
	g.field.prepareFor(g.diff)
}

func (g *Game) Load(diff Difficult) {
	g.diff = diff
	g.field.Reset()
	g.generate()
}

func (g *Game) ValuesCount() (counts map[int]int) {
	counts = make(map[int]int)
	for i := 1; i <= g.size; i++ {
		counts[i] = 0
		for _, cell := range g.GetCells() {
			if cell.Value().(int) == i {
				counts[i]++
			}
		}
	}
	return counts
}

func (g Game) Dim() int           { return g.dim }
func (g Game) Size() int          { return g.size }
func (g *Game) GetCells() []*Cell { return *g.field }
func (g *Game) GetField() *Field  { return g.field }

func (g *Game) Add(x, y, n int) {
	idx := g.field.Idx(x, y)
	if !g.GetCells()[idx].Add(n) {
		return
	}
	fmt.Println("Сделан ход:", n, idx, x, y, g.field.GetField()[idx].notes)
	for x0 := 0; x0 < g.size; x0++ {
		g.field.GetField()[g.field.Idx(x0, y)].UpdateNote(n)
	}
	for y0 := 0; y0 < g.size; y0++ {
		g.field.GetField()[g.field.Idx(x, y0)].UpdateNote(n)
	}

	rX0, rY0 := g.getRectIdx(x, y)
	for i, v := range g.field.GetField() {
		x1, y1 := g.field.Pos(i)
		rX, rY := g.getRectIdx(x1, y1)
		if rX0 != rX || rY0 != rY {
			continue
		}
		v.UpdateNote(n)
	}
	g.history = append(g.history, idx)
	fmt.Println("Результат хода:", n, idx, x, y, g.GetCells()[idx].notes, g.String())
}

func (g *Game) ResetCell(x, y int) {
	idx := g.GetField().Idx(x, y)
	if g.GetCells()[idx].IsReadOnly() {
		return
	}
	g.GetCells()[idx].Reset()
	fmt.Println("Обнулить ход:", idx, x, y, g.GetCells()[g.GetField().Idx(x, y)].notes)
	for x0 := 0; x0 < g.size; x0++ {
		value := g.GetCells()[g.GetField().Idx(x0, y)].Value().(int)
		g.GetCells()[g.GetField().Idx(x, y)].UpdateNote(value)
	}
	for y0 := 0; y0 < g.size; y0++ {
		value := g.GetCells()[g.GetField().Idx(x, y0)].Value().(int)
		g.GetCells()[g.GetField().Idx(x, y)].UpdateNote(value)
	}
	rX0, rY0 := g.getRectIdx(x, y)
	for i := range g.GetCells() {
		x0, y0 := g.GetField().Pos(i)
		rX, rY := g.getRectIdx(x0, y0)
		if rX0 != rX || rY0 != rY {
			continue
		}
		value := g.GetCells()[g.GetField().Idx(x0, y0)].Value().(int)
		g.GetCells()[g.GetField().Idx(x, y)].UpdateNote(value)
	}
	cell := g.GetCells()[g.GetField().Idx(x, y)]
	fmt.Println("Обнуление хода:", idx, x, y, cell.Value().(int), cell.notes, g.String())
}

func (g *Game) Undo() {
	if len(g.history) == 0 {
		return
	}
	x, y := g.LastMovePos()
	g.ResetCell(x, y)
	g.history = eui.PopIntSlice(g.history)
	log.Println("undo", x, y, g.history)
}

func (g Game) ReseAllCells(idx int) {
	if idx == 0 {
		return
	}
	x, y := g.GetField().Pos(idx - 1)
	if g.GetCells()[g.GetField().Idx(x, y)].GetValue() == 0 {
		g.ResetCell(x, y)
	}
	g.ReseAllCells(idx - 1)
}

func (g *Game) LastMovePos() (int, int) { return g.GetField().Pos(g.history[len(g.history)-1]) }

func (g Game) getRectIdx(x int, y int) (rX int, rY int) {
	szX := g.size
	szY := g.size
	rX = g.dim
	rY = g.dim
	for szX > x {
		szX -= g.dim
		rX--
	}
	for szY > y {
		szY -= g.dim
		rY--
	}
	return rX, rY
}

func (g Game) String() (result string) {
	result = fmt.Sprintf("sudoku %vX%v\n", g.size, g.size)
	for y := 0; y < g.size; y++ {
		for x := 0; x < g.size; x++ {
			result += g.GetCells()[g.GetField().Idx(x, y)].String()
		}
		result += "\n"
	}
	return result
}
