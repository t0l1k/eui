package game

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

type FieldConf struct {
	Dim         int // какой размер поля
	Colors      int // сколько цветов в игре
	Balls       int // сколько на следующем ходе появиться шаров
	WinLineSize int // сколько собрать в линию для удаления
}

func NewFieldConf(dim int) *FieldConf {
	f := &FieldConf{Dim: dim}
	f.WinLineSize = (dim + 1) / 2
	f.Colors = dim - 2
	f.Balls = f.WinLineSize - 2
	return f
}

type Field struct {
	field  []*Cell
	Conf   *FieldConf
	score  int
	InGame bool
}

func NewField(dim int) *Field { return &Field{Conf: NewFieldConf(dim)} }

func (f *Field) NewGame(dim int) {
	f.Conf = NewFieldConf(dim)
	f.field = nil
	f.score = 0
	for y := 0; y < f.Conf.Dim; y++ {
		for x := 0; x < f.Conf.Dim; x++ {
			f.field = append(f.field, NewCell(x, y))
		}
	}
	f.InGame = true
	log.Println("поле игры создано", f.Conf)
}

func (f *Field) MakeMove(x, y int) (col BallColor, way []int) {
	moveCell := f.field[f.Idx(x, y)]
	col = moveCell.color
	if moveCell.IsFilled() || moveCell.IsFilledAfterMove() {
		for _, cell := range f.field {
			if cell.IsMarkedForMove() {
				cell.SetFilled()
				log.Println("удалить ранее отмеченный для хода", moveCell.color, moveCell.pos)
			}
		}
		moveCell.SetMarkedForMove()
		log.Println("отмечаем для хода", moveCell.color, moveCell.pos)
	} else if moveCell.IsEmpty() || moveCell.IsFilledNext() {
		for i, cell := range f.field {
			if cell.IsMarkedForMove() {
				way = f.getWaveWay(i, f.Idx(x, y))
				if len(way) == 0 {
					log.Println("шарик недоступен!!!", len(way), way)
					return col, way
				}
				log.Println("путь к шарику:", way, col)
				if moveCell.IsFilledNext() {
					moveCell.SetFilledAfterMove(cell.color)
					f.NextMoveBall(col)
					log.Printf("ход на отмеченную для следующего хода [%v,%v] из [%v,%v]", moveCell.color, moveCell.pos, cell.color, cell.pos)
				} else {
					moveCell.SetFilledAfterMove(cell.color)
					log.Printf("ход в пустую клетку [%v,%v] из [%v,%v]", moveCell.color, moveCell.pos, cell.color, cell.pos)
				}
				f.field[i].Reset()
				log.Printf("сделан ход [%v,%v] из [%v,%v] путь:%v куда:%v из:%v", moveCell.color, moveCell.pos, cell.color, cell.pos, way, f.Idx(x, y), i)
				f.CheckNextMove(x, y)
				break
			}
		}
	}
	return col, way
}

func (f *Field) CheckNextMove(x, y int) {
	f.ShowFilledNext()
	if f.checkWinLines(x, y) {
		return
	}
	f.NextMoveBalls()
	if f.EmptyCellsCount() == 0 {
		f.ShowFilledNext()
		f.InGame = false
	}
}

func (f *Field) ShowFilledNext() {
	for i, cell := range f.field {
		if cell.IsFilledNext() {
			cell.SetFilled()
			x, y := f.Pos(i)
			f.checkWinLines(x, y)
			log.Printf("показываем новые шары цвет:%v поз:%v", cell.Color(), cell.Pos())
		}
	}
}

func (f *Field) NextMoveBalls() {
	rand.Seed(time.Now().UnixNano())
	count := f.Conf.Balls
	for count > 0 && f.EmptyCellsCount() > 0 {
		col := BallColor(rand.Intn(f.Conf.Colors) + 1)
		if f.NextMoveBall(col) {
			count--
		}
	}
	log.Println("новые шары для следующего хода")
}

func (f *Field) NextMoveBall(col BallColor) bool {
	for {
		x := rand.Intn(f.Conf.Dim)
		y := rand.Intn(f.Conf.Dim)
		cell := f.field[f.Idx(x, y)]
		if cell.IsEmpty() {
			cell.SetFilledNext(col)
			break
		}
	}
	return true
}

func (f *Field) checkWinLines(x, y int) bool {
	count := 0
	if f.checkHorizontal(x, y) {
		count++
	}
	if f.checkVertical(x, y) {
		count++
	}
	if f.checkDiagDown(x, y) {
		count++
	}
	if f.checkDiagUp(x, y) {
		count++
	}
	return count > 0
}

func (f *Field) checkHorizontal(x, y int) bool {
	arr := f.field[f.Idx(0, y):f.Idx(f.Conf.Dim, y)]
	log.Println("проверка собрано ли что-то для удаления по горизонтали", arr)
	return f.checkWinLine(arr)
}

func (f *Field) checkVertical(x, y int) bool {
	var arr []*Cell
	for i := 0; i < f.Conf.Dim; i++ {
		arr = append(arr, f.field[f.Idx(x, i)])
	}
	log.Println("проверка собрано ли что-то для удаления по вертикали", arr)
	return f.checkWinLine(arr)
}

func (f *Field) checkDiagDown(x, y int) bool {
	var arr []*Cell
	idx := f.Idx(x, y)
	for i := idx; i >= 0; i -= f.Conf.Dim + 1 {
		idx = i
	}
	for i := idx; i < f.Conf.Dim*f.Conf.Dim; i += f.Conf.Dim + 1 {
		arr = append(arr, f.field[i])
	}
	log.Println("проверка собрано ли что-то для удаления по диагонали вниз", arr)
	return f.checkWinLine(arr)
}

func (f *Field) checkDiagUp(x0, y0 int) bool {
	var arr []*Cell
	idx := f.Idx(x0, y0)
	for y, x := y0, x0; y < f.Conf.Dim && x >= 0; y++ {
		idx = f.Idx(x, y)
		x--
	}
	x1, y1 := f.Pos(idx)
	for y, x := y1, x1; y >= 0 && x < f.Conf.Dim; y-- {
		arr = append(arr, f.field[f.Idx(x, y)])
		x++
	}
	log.Println("проверка собрано ли что-то для удаления по диагонали вверх", arr)
	return f.checkWinLine(arr)
}

func (f *Field) checkWinLine(arr []*Cell) bool {
	if len(arr) < f.Conf.WinLineSize {
		return false
	}
	var line []*Cell
	count := 1
	col := arr[0].Color()
	line = append(line, arr[0])
	for i := 1; i < len(arr); i++ {
		newCol := arr[i].Color()
		if col == newCol && col != BallNoColor {
			count++
			col = newCol
			line = append(line, arr[i])
		} else if col != newCol {
			if count >= f.Conf.WinLineSize {
				break
			}
			count = 1
			col = newCol
			line = nil
			line = append(line, arr[i])
		}
		log.Printf("линия:%v(цвет:%v,новый цвет:%v,сколько:%v[%v])\n", i, col, newCol, count, line)
	}

	if len(line) >= f.Conf.WinLineSize {
		for i, v := range line {
			if i > f.Conf.WinLineSize {
				f.score += 2
			} else {
				f.score++
			}
			log.Println("отмечаем для удаления:", v.Color(), v.Pos(), f.score)
			v.SetMarkedForDelete()
		}
	}
	return count >= f.Conf.WinLineSize
}

func (f *Field) GetScore() int          { return f.score }
func (f *Field) GetField() []*Cell      { return f.field }
func (f *Field) Dim() (int, int)        { return f.Conf.Dim, f.Conf.Dim }
func (f *Field) GetCell(x, y int) *Cell { return f.field[f.Idx(x, y)] }
func (f *Field) Idx(x, y int) int       { return y*f.Conf.Dim + x }
func (f *Field) Pos(idx int) (int, int) { return idx % f.Conf.Dim, idx / f.Conf.Dim }

func (f *Field) EmptyCellsCount() (count int) {
	for _, v := range f.field {
		if v.color == BallNoColor {
			count++
		}
	}
	return count
}

func (f *Field) GetFilledNext() (cells []*Cell) {
	for _, v := range f.field {
		if v.IsFilledNext() {
			cells = append(cells, v)
		}
	}
	return cells
}

func (f *Field) isFieldEdge(x, y int) bool {
	return x < 0 || x > f.Conf.Dim-1 || y < 0 || y > f.Conf.Dim-1
}

func (f Field) getNeighthors(x, y int) (arr []int) {
	for dir := 0; dir < 4; dir++ {
		var dx, dy int
		switch dir {
		case 0:
			dx--
		case 1:
			dx++
		case 2:
			dy--
		case 3:
			dy++
		}
		nx := x + dx
		ny := y + dy
		idx := f.Idx(nx, ny)
		if !f.isFieldEdge(nx, ny) {
			arr = append(arr, idx)
		}
	}
	return arr
}

func (f Field) getWaveWay(a, b int) (way []int) {
	wave := make([]int, len(f.field))
	d := 0
	var dArr []int
	dArr = append(dArr, a)
	for wave[b] == 0 { // Волна
		for _, index := range dArr {
			x, y := f.Pos(index)
			for _, idx := range f.getNeighthors(x, y) {
				if (f.field[idx].IsEmpty() || f.field[idx].IsFilledNext()) && wave[idx] == 0 {
					wave[idx] = d + 1
				}
			}
		}
		d++
		dArr = dArr[0:0]
		for i, v := range wave {
			if v == d {
				dArr = append(dArr, i)
			}
		}
		if len(dArr) == 0 {
			return
		}
	}

	d = wave[b] // Путь
	way = append(way, b)
	index := b
	for d > 0 {
		x, y := f.Pos(index)
		for _, idx := range f.getNeighthors(x, y) {
			if wave[idx] == d-1 {
				way = append(way, idx)
				index = idx
				break
			}
		}
		d--
	}
	return way
}

func (f Field) String() (result string) {
	result = fmt.Sprintf("Поле:%vx%v\n", f.Conf.Dim, f.Conf.Dim)
	for y := 0; y < f.Conf.Dim; y++ {
		for x := 0; x < f.Conf.Dim; x++ {
			result += f.field[f.Idx(x, y)].String()
		}
		result += "\n"
	}
	return result
}
