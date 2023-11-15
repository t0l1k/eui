package eui

// Умею размеры виджетов во мне разделить одинаково на указаное число строк и рядов, начинаю из угла вверх-слева вправо указаное число строк, потом вниз на один ряд и опять с первой строчки вправо потом вниз и т.д
type GridLayoutRightDown struct {
	row, column int
	ContainerBase
	cellMargin int
}

func NewGridLayoutRightDown(r, c int) *GridLayoutRightDown {
	return &GridLayoutRightDown{row: r, column: c, cellMargin: 0}
}

func (d *GridLayoutRightDown) SetRows(r int)       { d.row = r }
func (d *GridLayoutRightDown) SetColumns(c int)    { d.column = c }
func (d *GridLayoutRightDown) SetCellMargin(m int) { d.cellMargin = m }

func (c *GridLayoutRightDown) Resize(rect []int) {
	r := NewRect(rect)
	w0, h0 := r.Size()
	x0, y0 := r.Pos()
	cellSize := c.getCellSize(r)
	marginX := (w0 - cellSize*c.row) / 2
	marginY := (h0 - cellSize*c.column) / 2
	x, y := x0+marginX, y0+marginY
	i := 0
	for _, icon := range c.Container {
		icon.Resize([]int{x, y, cellSize - c.cellMargin, cellSize - c.cellMargin})
		x += cellSize
		i++
		if i > 0 && i%c.row == 0 {
			x = x0 + marginX
			y += cellSize
		}
	}
}

func (g *GridLayoutRightDown) getCellSize(rect *Rect) (size int) {
	r := g.row
	c := g.column
	for r*size < rect.W && c*size < rect.H {
		size += 1
	}
	return size
}
