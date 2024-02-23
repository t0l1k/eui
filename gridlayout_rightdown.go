package eui

// Умею размеры виджетов во мне разделить одинаково на указаное число строк и рядов, начинаю из угла вверх-слева вправо указаное число строк, потом вниз на один ряд и опять с первой строчки вправо потом вниз и т.д
type GridLayoutRightDown struct {
	LayoutBase
	row, column int
	ContainerBase
	cellMargin int
}

func NewGridLayoutRightDown(r, c int) *GridLayoutRightDown {
	return &GridLayoutRightDown{row: r, column: c, cellMargin: 0}
}

func (d *GridLayoutRightDown) SetDim(r, c int) {
	d.row = r
	d.column = c
	d.resize()
}

func (d *GridLayoutRightDown) SetRows(r int)       { d.row = r; d.resize() }
func (d *GridLayoutRightDown) SetColumns(c int)    { d.column = c; d.resize() }
func (d *GridLayoutRightDown) SetCellMargin(m int) { d.cellMargin = m; d.resize() }

func (c *GridLayoutRightDown) resize() { c.Resize(c.GetRect().GetArr()) }

func (c *GridLayoutRightDown) Resize(rect []int) {
	c.Rect(NewRect(rect))
	w0, h0 := c.GetRect().Size()
	x0, y0 := c.GetRect().Pos()
	cellSize := c.getCellSize(c.GetRect())
	marginX := (w0 - cellSize*c.row) / 2
	marginY := (h0 - cellSize*c.column) / 2
	x, y := x0+marginX, y0+marginY
	i := 0
	for _, icon := range c.GetContainer() {
		icon.Resize([]int{x, y, cellSize - c.cellMargin, cellSize - c.cellMargin})
		x += cellSize
		i++
		if i > 0 && i%c.row == 0 {
			x = x0 + marginX
			y += cellSize
		}
	}
}

func (g *GridLayoutRightDown) getCellSize(rect Rect) (size int) {
	r := g.row
	c := g.column
	for r*size < rect.W && c*size < rect.H {
		size += 1
	}
	return size
}
