package eui

// Умею размеры виджетов во мне разделить одинаково на указаное число строк и рядов, начинаю из угла вверх-слева вправо указаное число строк, потом вниз на один ряд и опять с первой строчки вправо потом вниз и т.д
type GridLayoutRightDown struct {
	LayoutBase
	row, column, cellMargin, cellSize float64
	ContainerBase
	ItemsRect Rect
}

func NewGridLayoutRightDown(r, c float64) *GridLayoutRightDown {
	return &GridLayoutRightDown{row: r, column: c, cellMargin: 0}
}

func (c *GridLayoutRightDown) GetCellSize() (size float64) {
	row := c.row
	col := c.column
	for row*size < float64(c.GetRect().W) && col*size < float64(c.GetRect().H) {
		size += 0.01
	}
	return size
}
func (d *GridLayoutRightDown) SetDim(r, c float64)     { d.row = r; d.column = c; d.resize() }
func (d *GridLayoutRightDown) SetRows(r float64)       { d.row = r; d.resize() }
func (d *GridLayoutRightDown) SetColumns(c float64)    { d.column = c; d.resize() }
func (d *GridLayoutRightDown) SetCellMargin(m float64) { d.cellMargin = m; d.resize() }
func (c *GridLayoutRightDown) resize()                 { c.Resize(c.GetRect().GetArr()) }
func (c *GridLayoutRightDown) Resize(rect []int) {
	c.Rect(NewRect(rect))
	w0, h0 := c.GetRect().Size()
	x0, y0 := c.GetRect().Pos()

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
}
