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
func (a *AbsoluteLayout) Apply(children []Drawabler, r Rect[int]) {
	for _, c := range children {
		if c.ViewType().IsBackground() {
			c.SetRect(r)
		}
	}
	log.Println("AbsoluteLayout:Apply:", children, r)
}

type LayoutBox struct {
	dir     Orientation
	Spacing int // промежуток между элементами (px)
}

func NewVBoxLayout(spacing int) *LayoutBox { return &LayoutBox{dir: Vertical, Spacing: spacing} }
func NewHBoxLayout(spacing int) *LayoutBox { return &LayoutBox{dir: Horizontal, Spacing: spacing} }

func (b *LayoutBox) Apply(children []Drawabler, parentRect Rect[int]) {
	// 1. Фоновые
	for _, v := range children {
		if v.ViewType().IsBackground() {
			v.SetRect(parentRect)
		}
	}
	// 2. Обычные
	var normal []Drawabler
	for _, v := range children {
		if !v.ViewType().IsBackground() {
			normal = append(normal, v)
		}
	}
	n := len(normal)
	if n == 0 {
		return
	}
	w, h := parentRect.Width()-b.Spacing*2, parentRect.Height()-b.Spacing*2
	if b.dir == Vertical {
		totalSpacing := b.Spacing * (n - 1)
		cellH := (h - totalSpacing) / n
		x := parentRect.X + b.Spacing
		y := parentRect.Y + b.Spacing
		for _, c := range normal {
			r := NewRect([]int{x, y, w, cellH})
			c.SetRect(r)
			y += cellH + b.Spacing
		}
	} else {
		totalSpacing := b.Spacing * (n - 1)
		cellW := (w - totalSpacing) / n
		x := parentRect.X + b.Spacing
		y := parentRect.Y + b.Spacing
		for _, c := range normal {
			r := NewRect([]int{x, y, cellW, h})
			c.SetRect(r)
			x += cellW + b.Spacing
		}
	}
	log.Println("LayoutBox:Apply:", parentRect, children)
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

	// 1. Фоновые
	for _, v := range components {
		if v.ViewType().IsBackground() {
			v.SetRect(rect)
		}
	}
	// 2. Обычные
	var normal []Drawabler
	for _, v := range components {
		if !v.ViewType().IsBackground() {
			normal = append(normal, v)
		}
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

	for i, comp := range normal {
		row := i % l.Rows
		col := i / l.Rows
		x := offsetX + float64(row)*(cellW+float64(l.Spacing))
		y := offsetY + float64(col)*(cellH+float64(l.Spacing))
		r := NewRect([]int{int(x), int(y), int(cellW), int(cellH)})
		comp.SetRect(r)
	}
	log.Printf("GridLayout(%v,%v):Apply(%v)%v", l.Rows, l.Columns, len(normal), rect)
}

type StackLayout struct{ spacing int }

func NewStackLayout(spacing int) *StackLayout { return &StackLayout{spacing: spacing} }
func (b *StackLayout) Apply(children []Drawabler, rect Rect[int]) {
	if len(children) == 0 {
		return
	}
	for _, c := range children {
		if c.ViewType().IsBackground() {
			c.SetRect(rect)
			continue
		}
		x := rect.X + b.spacing
		y := rect.Y + b.spacing
		w := rect.W - b.spacing*2
		h := rect.H - b.spacing*2
		c.SetRect(NewRect([]int{x, y, w, h}))
	}
	found := false
	for _, c := range children {
		if c.ViewType().IsBackground() {
			continue
		}
		if !found {
			c.Show()
			found = true
			continue
		}
		c.Hide()
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
	// 1. Сначала фоновые виджеты
	for _, v := range children {
		if v.ViewType().IsBackground() {
			v.SetRect(parentRect)
		}
	}
	// 2. Обычные виджеты
	var normal []Drawabler
	for _, v := range children {
		if !v.ViewType().IsBackground() {
			normal = append(normal, v)
		}
	}
	n := len(normal)
	if n == 0 || len(l.Parts) != n {
		return
	}
	w, h := parentRect.Width()-l.Spacing*2, parentRect.Height()-l.Spacing*2
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
		for i, c := range normal {
			pw := avail * l.Parts[i] / total
			c.SetRect(NewRect([]int{x, y, pw, h}))
			x += pw + l.Spacing
		}
	} else {
		avail := h - l.Spacing*(n-1)
		for i, c := range normal {
			ph := avail * l.Parts[i] / total
			c.SetRect(NewRect([]int{x, y, w, ph}))
			y += ph + l.Spacing
		}
	}
	log.Println("LayoutPercent:Apply:", parentRect, children)
}
