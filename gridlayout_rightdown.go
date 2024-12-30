package eui

// Умею размеры виджетов во мне разделить одинаково на указаное число строк и рядов, начинаю из угла вверх-слева вправо указаное число строк, потом вниз на один ряд и опять с первой строчки вправо потом вниз и т.д
type GridLayoutRightDown struct {
	LayoutBase
	row, column, cellMargin, cellSize, cellSizeRow, cellSizeColumn float64
	ContainerBase
	ItemsRect Rect
	fitToDim  bool
}

func NewGridLayoutRightDown(r, c float64) *GridLayoutRightDown {
	return &GridLayoutRightDown{row: r, column: c, cellMargin: 0, fitToDim: false}
}

func NewGridLayoutRightDownFitToDim(r, c float64) *GridLayoutRightDown {
	return &GridLayoutRightDown{row: r, column: c, cellMargin: 0, fitToDim: true}
}

// Вычислить квадратной ячейку
func (c *GridLayoutRightDown) GetCellSize() (size float64) {
	row := c.row
	col := c.column
	for row*size < float64(c.GetRect().W) && col*size < float64(c.GetRect().H) {
		size += 0.01
	}
	return size
}

func (c *GridLayoutRightDown) GetRowSize() (size float64) {
	row := c.row
	for row*size < float64(c.GetRect().W) {
		size += 0.01
	}
	return size
}
func (c *GridLayoutRightDown) GetColumnSize() (size float64) {
	col := c.column
	for col*size < float64(c.GetRect().H) {
		size += 0.01
	}
	return size
}
func (d *GridLayoutRightDown) FitToDim(value bool)     { d.fitToDim = value; d.resize() }
func (d *GridLayoutRightDown) SetDim(r, c float64)     { d.row = r; d.column = c; d.resize() }
func (d *GridLayoutRightDown) SetRows(r float64)       { d.row = r; d.resize() }
func (d *GridLayoutRightDown) SetColumns(c float64)    { d.column = c; d.resize() }
func (d *GridLayoutRightDown) SetCellMargin(m float64) { d.cellMargin = m; d.resize() }
func (c *GridLayoutRightDown) resize()                 { c.Resize(c.GetRect().GetArr()) }
func (c *GridLayoutRightDown) Resize(rect []int) {
	c.Rect(NewRect(rect))
	w0, h0 := c.GetRect().Size()
	x0, y0 := c.GetRect().Pos()

	if !c.fitToDim {
		c.cellSize = c.GetCellSize()
		marginX := (float64(w0) - c.cellSize*c.row) / 2
		marginY := (float64(h0) - c.cellSize*c.column) / 2
		x, y := float64(x0)+marginX, float64(y0)+marginY
		c.ItemsRect = NewRect([]int{int(x), int(y), int(c.cellSize * c.row), int(c.cellSize * c.column)})
		i := 0
		for _, icon := range c.GetContainer() {
			icon.Resize([]int{int(x), int(y), int(c.cellSize - c.cellMargin), int(c.cellSize - c.cellMargin)})
			x += c.cellSize
			i++
			if i > 0 && i%int(c.row) == 0 {
				x = float64(x0) + marginX
				y += c.cellSize
			}
		}
	} else {
		c.cellSizeRow = c.GetRowSize()
		c.cellSizeColumn = c.GetColumnSize()
		marginX := (float64(w0) - c.cellSizeRow*c.row) / 2
		marginY := (float64(h0) - c.cellSizeColumn*c.column) / 2
		x, y := float64(x0)+marginX, float64(y0)+marginY
		c.ItemsRect = NewRect([]int{int(x), int(y), int(c.cellSizeRow * c.row), int(c.cellSizeColumn * c.column)})
		i := 0
		for _, icon := range c.GetContainer() {
			icon.Resize([]int{int(x), int(y), int(c.cellSizeRow - c.cellMargin), int(c.cellSizeColumn - c.cellMargin)})
			x += c.cellSizeRow
			i++
			if i > 0 && i%int(c.row) == 0 {
				x = float64(x0) + marginX
				y += c.cellSizeColumn
			}
		}

	}
}
