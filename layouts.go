package eui

import (
	"log"
)

type AbsoluteLayout struct{}

func NewAbsoluteLayout() *AbsoluteLayout              { return &AbsoluteLayout{} }
func (a *AbsoluteLayout) Apply(d []Drawabler, r Rect) { log.Println("AbsoluteLayout:Apply:", d, r) }

// Умею размеры виджетов во мне разделить одинаково по горизонтали или по вертикали
type BoxLayout struct {
	horizontal bool
	spacing    float64
}

func NewHBoxLayout(spacing float64) *BoxLayout { return &BoxLayout{horizontal: true, spacing: spacing} }
func NewVBoxLayout(spacing float64) *BoxLayout {
	return &BoxLayout{horizontal: false, spacing: spacing}
}
func (l *BoxLayout) Apply(widgets []Drawabler, rect Rect) {
	count := len(widgets)
	if count == 0 {
		return
	}
	x0, y0 := rect.Pos()
	if l.horizontal {
		width := float64(rect.Width()-(count-1)*int(l.spacing)) / float64(count)
		w, h := width, rect.Height()
		for i, v := range widgets {
			x, y := w*float64(i)+l.spacing, l.spacing
			v.Resize(NewRect([]int{x0 + int(x), y0 + int(y), int(w), h}))
		}
	} else {
		height := (float64(rect.Height()) - float64(count-1)*l.spacing) / float64(count)
		w, h := rect.Width(), height
		for i, v := range widgets {
			x, y := l.spacing, h*float64(i)+l.spacing
			v.Resize(NewRect([]int{x0 + int(x), y0 + int(y), w, int(h)}))
		}
	}
	log.Println("BoxLayout:Apply", count, widgets, rect)
}

// Умею размеры виджетов во мне разделить одинаково на указаное число строк и рядов, начинаю из угла вверх-слева вправо указаное число строк, потом вниз на один ряд и опять с первой строчки вправо потом вниз и т.д
type GridLayout struct {
	rows, columns, spacing float64
	square                 bool
}

func NewGridLayout(r, c, spacing float64) *GridLayout {
	return &GridLayout{rows: r, columns: c, spacing: spacing, square: false}
}

func NewSquareGridLayout(r, c, spacing float64) *GridLayout {
	return &GridLayout{rows: r, columns: c, spacing: spacing, square: true}
}

func (c *GridLayout) GetCellSize(r Rect) (size float64) {
	row := c.rows
	col := c.columns
	for row*size < float64(r.W) && col*size < float64(r.H) {
		size += 0.01
	}
	return size
}

func (c *GridLayout) GetRowSize(r Rect) (size float64) {
	row := c.rows
	for row*size < float64(r.W) {
		size += 0.01
	}
	return size
}
func (c *GridLayout) GetColumnSize(r Rect) (size float64) {
	col := c.columns
	for col*size < float64(r.H) {
		size += 0.01
	}
	return size
}

func (c *GridLayout) Apply(components []Drawabler, rect Rect) {
	w0, h0 := rect.Size()
	x0, y0 := rect.Pos()
	if !c.square {
		cellSize := c.GetCellSize(rect)
		marginX := (float64(w0) - cellSize*c.rows) / 2
		marginY := (float64(h0) - cellSize*c.columns) / 2
		x, y := float64(x0)+marginX, float64(y0)+marginY
		i := 0
		for _, icon := range components {
			icon.Resize(NewRect([]int{int(x), int(y), int(cellSize - c.spacing), int(cellSize - c.spacing)}))
			x += cellSize
			i++
			if i > 0 && i%int(c.rows) == 0 {
				x = float64(x0) + marginX
				y += cellSize
			}
		}
	} else {
		cellSizeRow := c.GetRowSize(rect)
		cellSizeColumn := c.GetColumnSize(rect)
		marginX := (float64(w0) - cellSizeRow*c.rows) / 2
		marginY := (float64(h0) - cellSizeColumn*c.columns) / 2
		x, y := float64(x0)+marginX, float64(y0)+marginY
		i := 0
		for _, icon := range components {
			icon.Resize(NewRect([]int{int(x), int(y), int(cellSizeRow - c.spacing), int(cellSizeColumn - c.spacing)}))
			x += cellSizeRow
			i++
			if i > 0 && i%int(c.rows) == 0 {
				x = float64(x0) + marginX
				y += cellSizeColumn
			}
		}
	}
	log.Printf("GridLayout(%v,%v):Apply(%v)%v", c.rows, c.columns, len(components), rect)
}
