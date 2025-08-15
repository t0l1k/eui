package eui

import (
	"log"
)

type Orientation int

const (
	Horizontal Orientation = iota
	Vertical
)

type AbsoluteLayout struct{}

func NewAbsoluteLayout() *AbsoluteLayout { return &AbsoluteLayout{} }
func (a *AbsoluteLayout) Apply(d []Drawabler, r Rect[int]) {
	log.Println("AbsoluteLayout:Apply:", d, r)
}

// Умею размеры виджетов во мне разделить одинаково по горизонтали или по вертикали
type BoxLayout struct {
	dir Orientation
	// horizontal bool
	spacing float64
}

func NewHBoxLayout(spacing float64) *BoxLayout { return &BoxLayout{dir: Horizontal, spacing: spacing} }
func NewVBoxLayout(spacing float64) *BoxLayout {
	return &BoxLayout{dir: Vertical, spacing: spacing}
}
func (l *BoxLayout) Apply(widgets []Drawabler, rect Rect[int]) {
	count := len(widgets)
	if count == 0 {
		return
	}
	x0, y0 := rect.Pos()
	if l.dir == Horizontal {
		width := float64(rect.Width()-(count-1)*int(l.spacing)) / float64(count)
		w, h := width, rect.Height()
		for i, v := range widgets {
			x, y := w*float64(i)+l.spacing, l.spacing
			v.SetRect(NewRect([]int{x0 + int(x), y0 + int(y), int(w), h}))
		}
	} else {
		height := (float64(rect.Height()) - float64(count-1)*l.spacing) / float64(count)
		w, h := rect.Width(), height
		for i, v := range widgets {
			x, y := l.spacing, h*float64(i)+l.spacing
			v.SetRect(NewRect([]int{x0 + int(x), y0 + int(y), w, int(h)}))
		}
	}
	log.Println("BoxLayout:Apply", count, widgets, rect)
}

type GridLayout struct {
	Rows, Columns int
	Spacing       int
	square        bool
}

func NewGridLayout(rows, cols, spacing int) *GridLayout {
	return &GridLayout{Rows: rows, Columns: cols, Spacing: spacing, square: false}
}

func NewSquareGridLayout(rows, cols, spacing int) *GridLayout {
	return &GridLayout{Rows: rows, Columns: cols, Spacing: spacing, square: true}
}

func (l *GridLayout) Apply(components []Drawabler, rect Rect[int]) {
	count := len(components)
	if count == 0 || l.Rows == 0 || l.Columns == 0 {
		return
	}
	cellW := float64(rect.Width()-(l.Rows-1)*l.Spacing) / float64(l.Rows)
	cellH := float64(rect.Height()-(l.Columns-1)*l.Spacing) / float64(l.Columns)
	if l.square {
		cellSize := min(cellH, cellW)
		cellW, cellH = cellSize, cellSize
	}

	// Центрирование
	gridW := cellW*float64(l.Rows) + float64(l.Rows-1)*float64(l.Spacing)
	gridH := cellH*float64(l.Columns) + float64(l.Columns-1)*float64(l.Spacing)
	offsetX := float64(rect.X) + (float64(rect.Width())-gridW)/2
	offsetY := float64(rect.Y) + (float64(rect.Height())-gridH)/2

	for i, comp := range components {
		row := i % l.Rows
		col := i / l.Rows
		x := offsetX + float64(row)*(cellW+float64(l.Spacing))
		y := offsetY + float64(col)*(cellH+float64(l.Spacing))
		r := NewRect([]int{int(x), int(y), int(cellW), int(cellH)})
		comp.SetRect(r)
	}
	log.Printf("GridLayout(%v,%v):Apply(%v)%v", l.Rows, l.Columns, len(components), rect)
}

type StackLayout struct{ spacing int }

func NewStackLayout(spacing int) *StackLayout { return &StackLayout{spacing: spacing} }
func (b *StackLayout) Apply(children []Drawabler, rect Rect[int]) {
	if len(children) == 0 {
		return
	}
	for _, c := range children {
		x := rect.X + b.spacing
		y := rect.Y + b.spacing
		w := rect.W - b.spacing*2
		h := rect.H - b.spacing*2
		c.SetRect(NewRect([]int{x, y, w, h}))
	}
	log.Println("StackLayout:Apply:", children, rect)
}

type LayoutPercent struct {
	dir     Orientation
	Parts   []int // проценты для каждого ребёнка (сумма не обязательно 100)
	Spacing int   // промежуток между элементами (px)
}

func NewLayoutVerticalPercent(parts []int, spacing int) *LayoutPercent {
	return &LayoutPercent{dir: Vertical, Parts: parts, Spacing: spacing}
}
func NewLayoutHorizontalPercent(parts []int, spacing int) *LayoutPercent {
	return &LayoutPercent{dir: Horizontal, Parts: parts, Spacing: spacing}
}
func (l *LayoutPercent) Apply(children []Drawabler, parentRect Rect[int]) {
	n := len(children)
	if n == 0 || len(l.Parts) != n {
		return
	}
	w, h := parentRect.W-l.Spacing*2, parentRect.H-l.Spacing*2
	total := 0
	for _, p := range l.Parts {
		total += p
	}
	if total == 0 {
		return
	}
	x, y := parentRect.X+l.Spacing, parentRect.Y+l.Spacing
	if l.dir == Horizontal {
		avail := w - l.Spacing*(n-1)
		for i, c := range children {
			pw := avail * l.Parts[i] / total
			c.SetRect(NewRect([]int{x, y, pw, h}))
			x += pw + l.Spacing
		}
	} else {
		avail := h - l.Spacing*(n-1)
		for i, c := range children {
			ph := avail * l.Parts[i] / total
			c.SetRect(NewRect([]int{x, y, w, ph}))
			y += ph + l.Spacing
		}
	}
	log.Println("LayoutPercent:Apply:", parentRect, children)
}
