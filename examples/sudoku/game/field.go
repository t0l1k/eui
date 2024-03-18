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

func (f *Field) New() {
	rand.Seed(time.Now().UnixNano())
	countA := 0
	countB := 0
	count := 0
	for y := 0; y < f.size; y++ {
		for x := 0; x < f.size; x++ {
			for {
				n := f.nextRand()
				if countA > f.size*100 {
					idx := f.Idx(x, y)
					f.cells[f.Idx(x, y)].SetValue(0)
					if idx > 0 {
						idx--
					}
					f.cells[idx].SetValue(0)
					x, y = f.Pos(idx)
					fmt.Println("Шаг назад:", x, y, idx, n, countA, countB, f)
					countA = 0
					countB++
				}
				if countB > f.dim {
					for i := f.Idx(x, y); i > 0; i-- {
						f.cells[i].SetValue(0)
					}
					x, y = 0, 0
					countA = 0
					countB = 0
				}
				if n == 0 || n == f.size+1 {
					panic(fmt.Sprintf("got %v %v\n%v", n, countB, f))
				}
				if f.IsValid(x, y, n) {
					f.cells[f.Idx(x, y)].SetValue(n)
					fmt.Println("Генерация успешна:", x, y, n, count)
					fmt.Println(f)
					break
				}
				countA++
				// fmt.Println("Генерация:", x, y, n, count)
				count++
			}
		}
	}
}

func (f *Field) GetCells() []*Cell     { return f.cells }
func (f *Field) nextRand() int         { return rand.Intn(f.size) + 1 }
func (f Field) Idx(x, y int) int       { return y*f.size + x }
func (f Field) Pos(idx int) (int, int) { return idx % f.size, idx / f.size } //x,y
func (f Field) IsValid(x, y, value int) bool {
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
			fmt.Println("rect result 1:", x0, y0, value, rX0, rY0, len(arr), arr, i, arr[i])
			return false
		}
	}
	fmt.Println("rect result 0:", x0, y0, value, rX0, rY0, len(arr), arr)
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
