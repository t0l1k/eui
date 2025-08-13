package eui

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

type ListView struct {
	*Drawable
	children           []Drawabler
	list               []string
	itemSize, rows     int
	contentRect        Rect[int]
	contentImage       *ebiten.Image
	offset, lastOffset int
	cameraRect         image.Rectangle
	isDragging         bool
	dragStartY         int
	dragStartOffset    int
}

func NewListView() *ListView { return &ListView{Drawable: NewDrawable(), rows: 1, itemSize: 30} }

func (c *ListView) Items() []Drawabler { return c.children }

func (l *ListView) SetupListViewButtons(list []string, itemSize, rows int, bg, fg color.Color, f func(b *Button)) {
	l.itemSize = itemSize
	l.rows = rows
	l.list = list
	for _, v := range l.list {
		btn := NewButton(v, f)
		l.AddItem(btn)
		btn.Bg(bg)
		btn.Fg(fg)
	}
	l.resizeChilds()
}

func (l *ListView) SetupListViewCheckBoxs(list []string, itemSize, rows int, bg, fg color.Color, f func(b *Checkbox)) {
	l.itemSize = itemSize
	l.rows = rows
	l.list = list
	for _, v := range l.list {
		chkBox := NewCheckbox(v, f)
		l.AddItem(chkBox)
		chkBox.Bg(bg)
		chkBox.Fg(fg)
	}
	l.resizeChilds()
}

func (l *ListView) GetCheckBoxes() (values []*Checkbox) {
	for _, v := range l.children {
		switch value := v.(type) {
		case *Checkbox:
			values = append(values, value)
		}
	}
	return values
}

func (l *ListView) SetupListViewText(list []string, itemSize, rows int, bg, fg color.Color) {
	l.itemSize = itemSize
	l.rows = rows
	l.list = list
	for _, v := range l.list {
		lbl := NewText(v)
		l.AddItem(lbl)
		lbl.Bg(bg)
		lbl.Fg(fg)
	}
	l.resizeChilds()
}

func (l *ListView) SetListViewTextWithBgFgColors(list []string, bg, fg []color.Color) {
	l.list = list
	for i, str := range l.list {
		lbl := NewText(str)
		l.AddItem(lbl)
		lbl.Bg(bg[i])
		lbl.Fg(fg[i])
	}
	l.resizeChilds()
}

func (l *ListView) AddBgFg(d Drawabler, bg, fg color.Color) {
	l.children = append(l.children, d)
	switch value := d.(type) {
	case *Text:
		value.Bg(bg)
		value.Fg(fg)
		l.list = append(l.list, value.GetText())
	case *Button:
		value.Bg(bg)
		value.Fg(fg)
		l.list = append(l.list, value.Text())
	case *Checkbox:
		value.Bg(bg)
		value.Fg(fg)
		l.list = append(l.list, value.Text())
	}
	l.resizeChilds()
}

func (l *ListView) AddItem(d Drawabler) {
	theme := GetUi().theme
	bg := theme.Get(ListViewItemBg)
	fg := theme.Get(ListViewItemFg)
	l.AddBgFg(d, bg, fg)
}

func (l *ListView) Itemsize(itemSize int) {
	if l.itemSize == itemSize {
		return
	}
	l.itemSize = itemSize
	l.MarkDirty()
}

func (l *ListView) Rows(rows int) {
	if l.rows == rows {
		return
	}
	l.rows = rows
	l.MarkDirty()
}

func (l *ListView) Reset() {
	l.list = nil
	for _, v := range l.children {
		v.Close()
	}
	l.children = nil
	l.contentImage = nil
	l.offset = 0
	l.lastOffset = 0
	l.cameraRect = image.Rect(0, 0, l.rect.W, l.rect.H)
	l.MarkDirty()
}

func (l *ListView) Layout() {
	if l.contentRect.IsEmpty() {
		str := fmt.Sprintf("ListView:Layout:contentRect empty %v %v %v", l.Rect(), l.cameraRect.String(), l.contentRect.String())
		panic(str)
	}
	l.Drawable.Layout()
	w0, h0 := l.contentRect.Size()
	if l.contentImage == nil || w0 != l.contentImage.Bounds().Dx() || h0 != l.contentImage.Bounds().Dy() {
		l.contentImage = nil
		l.contentImage = ebiten.NewImage(w0, h0)
	} else {
		l.contentImage.Clear()
	}
	l.contentImage.Fill(colornames.Fuchsia)
	for _, v := range l.children {
		switch value := v.(type) {
		case *Text, *Button, *Checkbox:
			value.Draw(l.contentImage)
		}
	}
	l.ClearDirty()
}

func (l *ListView) Hit(pt Point[int]) Drawabler {
	if !pt.In(l.rect) || l.IsHidden() {
		return nil
	}
	log.Println("ListView:Hit:", l.Rect(), pt)
	return l
}

func (l *ListView) MouseDown(md MouseData) {
	if l.contentRect.H <= l.rect.H {
		return
	}
	l.isDragging = true
	l.dragStartY = md.Pos().Y
	l.dragStartOffset = l.cameraRect.Min.Y
	log.Println("ListView:MouseDown:", l.Rect(), l.cameraRect.String(), l.isDragging, l.dragStartY, l.dragStartOffset, md)
}

func (l *ListView) MouseDrag(md MouseData) {
	if !l.isDragging {
		return
	}
	delta := l.dragStartY - md.Pos().Y
	newY := l.dragStartOffset + delta
	maxY := l.contentRect.H - l.rect.H
	if newY < 0 {
		newY = 0
	}
	if newY > maxY {
		newY = maxY
	}
	l.cameraRect = image.Rect(0, newY, l.rect.W, newY+l.rect.H)
	l.MarkDirty()
	log.Println("ListView:MouseDrag:", delta, newY, maxY, l.cameraRect)
}

func (l *ListView) MouseUp(md MouseData) {
	l.isDragging = false
	if md.Pos().Y != l.dragStartY {
		return
	}

	x0 := md.Pos().X - l.rect.X + l.cameraRect.Min.X
	y0 := md.Pos().Y - l.rect.Y + l.cameraRect.Min.Y

	pos := NewPoint(x0, y0)

	for _, d := range l.children {
		if mh, ok := d.(interface{ Hit(Point[int]) Drawabler }); ok {
			if mh.Hit(pos) != nil && !d.State().IsFocused() {
				d.SetState(StateFocused)

				if mp, ok := d.(interface{ MouseUp(MouseData) }); ok {
					mp.MouseUp(md)
				}
				if mp, ok := d.(interface{ WantBlur() bool }); ok {
					if mp.WantBlur() {
						d.SetState(StateNormal)
					}
				}

			} else if !d.State().IsBlurred() {
				d.SetState(StateNormal)
			}
		}
	}
	log.Println("ListView:MouseUp:", l.isDragging, l.cameraRect, x0, y0)
}

func (l *ListView) MouseWheel(md MouseData) {
	if l.contentRect.H <= l.rect.H {
		return
	}
	scrollStep := float64(l.itemSize / 2)
	delta := md.WPos().Y * scrollStep
	newY := l.cameraRect.Min.Y - int(delta)
	maxY := l.contentRect.H - l.rect.H
	if newY < 0 {
		newY = 0
	}
	if newY > maxY {
		newY = maxY
	}
	l.cameraRect = image.Rect(0, newY, l.rect.W, newY+l.rect.H)
	l.MarkDirty()
	log.Println("ListView:MouseWheel:", scrollStep, delta, newY, maxY, l.cameraRect)
}
func (l *ListView) WantBlur() bool { return true }

func (l *ListView) Draw(surface *ebiten.Image) {
	if l.IsHidden() {
		return
	}
	if l.IsDirty() {
		l.Layout()
	}
	op := &ebiten.DrawImageOptions{}
	x0, y0 := l.rect.Pos()
	op.GeoM.Translate(float64(x0), float64(y0))
	surface.DrawImage(l.contentImage.SubImage(l.cameraRect).(*ebiten.Image), op)
}

func (l *ListView) SetRect(r Rect[int]) {
	l.Drawable.SetRect(r)
	l.resizeChilds()
	l.cameraRect = image.Rect(0, 0, l.rect.W, l.rect.H)
	l.MarkDirty()
}

func (l *ListView) resizeChilds() {
	x, y := 0, 0
	w, h := l.rect.W/l.rows, l.itemSize
	row := 0
	col := 1
	for _, v := range l.children {
		switch value := v.(type) {
		case *Text, *Button, *Checkbox:
			value.SetRect(NewRect([]int{x, y, w - 1, h - 1}))
		}
		x += w
		row++
		if row > l.rows-1 {
			row = 0
			x = 0
			y += h
			col++
		}
	}
	if y == 0 {
		y = l.rect.H
	} else {
		if len(l.children)%l.rows == 0 {
			col--
		}
		y = col * l.itemSize
	}
	l.contentRect = NewRect([]int{0, 0, l.rect.W, y})
	l.MarkDirty()
}
