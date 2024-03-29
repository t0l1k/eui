package eui

// Умею размеры виджетов во мне разделить одинаково на указаное число строк и рядов, начинаю из угла вверх-слева вниз указаное число рядов, потом влево на одну строку и со следующей строчки сверху вниз и т.д
type GridLayoutDownRight struct {
	row, column int
	ContainerBase
	cellMargin int
}

func NewGridLayoutDownRight(r, c int) *GridLayoutDownRight {
	return &GridLayoutDownRight{row: r, column: c, cellMargin: 0}
}

func (d *GridLayoutDownRight) SetRows(r int)       { d.row = r }
func (d *GridLayoutDownRight) SetColumns(c int)    { d.column = c }
func (d *GridLayoutDownRight) SetCellMargin(m int) { d.cellMargin = m }

func (c *GridLayoutDownRight) Resize(rect []int) {
	r := NewRect(rect)
	w0, h0 := r.Size()
	x0, y0 := r.Pos()
	sz := c.column
	size := w0
	if w0 > h0 {
		size = h0
	} else {
		size = w0
	}
	cellSize := size / sz
	marginX := (w0 - cellSize*c.row) / 2
	marginY := (h0 - cellSize*c.column) / 2
	x, y := x0+marginX, y0+marginY
	i := 0
	for _, icon := range c.GetContainer() {
		icon.Resize([]int{x, y, cellSize - c.cellMargin, cellSize - -c.cellMargin})
		y += cellSize
		i++
		if i > 0 && i%c.column == 0 {
			y = y0 + marginY
			x += cellSize
		}
	}
}
