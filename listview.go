package eui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type ListView struct {
	View
	list               []string
	itemSize, rows     int
	contentRect        Rect
	contentImage       *ebiten.Image
	offset, lastOffset int
	cameraRect         image.Rectangle
}

func NewListView() *ListView {
	l := &ListView{rows: 1, itemSize: 30}
	l.SetupView()
	return l
}

func (l *ListView) SetupListViewButtons(list []string, itemSize, rows int, bg, fg color.Color, f func(b *Button)) {
	l.itemSize = itemSize
	l.rows = rows
	l.list = list
	for _, v := range l.list {
		btn := NewButton(v, f)
		l.Add(btn)
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
		l.Add(chkBox)
		chkBox.Bg(bg)
		chkBox.Fg(fg)
	}
	l.resizeChilds()
}

func (l *ListView) GetCheckBoxes() (values []*Checkbox) {
	for _, v := range l.GetContainer() {
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
		l.Add(lbl)
		lbl.Bg(bg)
		lbl.Fg(fg)
	}
	l.resizeChilds()
}

func (l *ListView) SetListViewTextWithBgFgColors(list []string, bg, fg []color.Color) {
	l.list = list
	for i, str := range l.list {
		lbl := NewText(str)
		l.Add(lbl)
		lbl.Bg(bg[i])
		lbl.Fg(fg[i])
	}
	l.resizeChilds()
}

func (l *ListView) Add(d Drawabler) {
	l.ContainerBase.Add(d)
	theme := GetUi().theme
	bg := theme.Get(ListViewItemBg)
	fg := theme.Get(ListViewItemFg)
	switch value := d.(type) {
	case *Text:
		value.Bg(bg)
		value.Fg(fg)
		l.list = append(l.list, value.GetText())
	case *Button:
		value.Bg(bg)
		value.Fg(fg)
		l.list = append(l.list, value.GetText())
	case *Checkbox:
		value.Bg(bg)
		value.Fg(fg)
		l.list = append(l.list, value.GetText())
	}
	l.resizeChilds()
}

func (l *ListView) Itemsize(itemSize int) {
	if l.itemSize == itemSize {
		return
	}
	l.itemSize = itemSize
	l.Dirty = true
}

func (l *ListView) Rows(rows int) {
	if l.rows == rows {
		return
	}
	l.rows = rows
	l.Dirty = true
}

func (l *ListView) Reset() {
	l.list = nil
	l.ResetContainerBase()
	l.contentImage = nil
	l.offset = 0
	l.lastOffset = 0
	l.cameraRect = image.Rect(0, 0, l.rect.W, l.rect.H)
	l.Dirty = true
}

func (l *ListView) Layout() {
	l.View.Layout()
	w0, h0 := l.contentRect.Size()
	if l.contentImage == nil || w0 != l.contentImage.Bounds().Dx() || h0 != l.contentImage.Bounds().Dy() {
		l.contentImage = nil
		l.contentImage = ebiten.NewImage(w0, h0)
	} else {
		l.contentImage.Clear()
	}
	l.contentImage.Fill(l.bg)
	for _, v := range l.GetContainer() {
		switch value := v.(type) {
		case *Text, *Button, *Checkbox:
			value.Draw(l.contentImage)
		}
	}
	l.Dirty = false
}

func (l *ListView) Update(dt int) {
	if l.isDragging {
		l.offset = -(l.dragEndPoint.Offset(l.dragStartPoint).Y)
		y := l.cameraRect.Min.Y + l.offset - l.lastOffset
		h := y + l.rect.H
		height := l.contentRect.H - l.rect.H
		if y <= 0 {
			y = 0
		} else if y >= height {
			y = height
		}
		if h < l.rect.H {
			h = l.rect.H
		}
		l.cameraRect = image.Rect(0, y, l.rect.W, h)
		l.lastOffset = l.offset
	} else if !l.isDragging {
		l.lastOffset = 0
	}
	if l.state != ViewStateNormal {
		x0 := l.dragEndPoint.X - l.rect.X + l.cameraRect.Min.X
		y0 := l.dragEndPoint.Y - l.rect.Y + l.cameraRect.Min.Y
		for _, v := range l.GetContainer() {
			switch value := v.(type) {
			case *Button:
				r := value.rect
				if r.InRect(x0, y0) && value.state != l.state {
					value.SetState(l.state)
					if l.lastOffset == 0 && l.offset == 0 {
						value.Update(dt)
					}
					l.Dirty = true
				} else if value.state != ViewStateNormal {
					value.SetState(ViewStateNormal)
				}
			case *Checkbox:
				r := value.btn.rect
				if r.InRect(x0, y0) && value.btn.state != l.state {
					value.btn.SetState(l.state)
					if l.lastOffset == 0 && l.offset == 0 {
						value.btn.Update(dt)
					}
					l.Dirty = true
				} else if value.btn.state != ViewStateNormal {
					value.btn.SetState(ViewStateNormal)
				}
			}
		}
	}
}

func (l *ListView) Draw(surface *ebiten.Image) {
	if !l.visible || l.disabled {
		return
	}
	if l.Dirty {
		l.Layout()
	}
	op := &ebiten.DrawImageOptions{}
	x0, y0 := l.rect.Pos()
	op.GeoM.Translate(float64(x0), float64(y0))
	surface.DrawImage(l.contentImage.SubImage(l.cameraRect).(*ebiten.Image), op)
}

func (l *ListView) Resize(r []int) {
	l.Rect(NewRect(r))
	l.resizeChilds()
	l.cameraRect = image.Rect(0, 0, l.rect.W, l.rect.H)
}

func (l *ListView) resizeChilds() {
	x, y := 0, 0
	w, h := l.rect.W/l.rows, l.itemSize
	row := 0
	col := 1
	for _, v := range l.GetContainer() {
		switch value := v.(type) {
		case *Text, *Button, *Checkbox:
			value.Resize([]int{x, y, w - 1, h - 1})
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
		if len(l.GetContainer())%l.rows == 0 {
			col--
		}
		y = col * l.itemSize
	}
	l.contentRect = NewRect([]int{0, 0, l.rect.W, y})
	l.Dirty = true
}
