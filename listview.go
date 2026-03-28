package eui

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

type ListView struct {
	*Drawable
	children            []Drawabler
	list                []string
	index               int
	itemSize, rows      int
	contentRect         Rect[int]
	contentImage        *ebiten.Image
	offset, lastOffset  int
	cameraRect          image.Rectangle
	isDragging          bool
	dragStartY          int
	dragStartOffset     int
	longText            bool
	isHoveringScrollbar bool
	isDraggingScrollbar bool
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
		btn.SetBg(bg)
		btn.SetFg(fg)
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
		chkBox.SetBg(bg)
		chkBox.SetFg(fg)
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
		lbl := NewLabel(v)
		l.AddItem(lbl)
		lbl.SetBg(bg)
		lbl.SetFg(fg)
	}
	l.resizeChilds()
}

func (l *ListView) SetListViewTextWithBgFgColors(list []string, bg, fg []color.Color) {
	l.list = list
	for i, str := range l.list {
		lbl := NewLabel(str)
		l.AddItem(lbl)
		lbl.SetBg(bg[i])
		lbl.SetFg(fg[i])
	}
	l.resizeChilds()
}

func (l *ListView) AddLongText(txt *Label) {
	l.longText = true
	l.rows = 1
	l.itemSize = 30
	l.AddItem(txt)
	l.resizeChilds()
}

func (l *ListView) AddBgFg(d Drawabler, bg, fg color.Color) {
	l.children = append(l.children, d)
	switch value := d.(type) {
	case *Label:
		value.SetBg(bg)
		value.SetFg(fg)
		l.list = append(l.list, value.Text())
	case *Button:
		value.SetBg(bg)
		value.SetFg(fg)
		l.list = append(l.list, value.Text())
	case *Checkbox:
		value.SetBg(bg)
		value.SetFg(fg)
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

func (l *ListView) Close() {
	for i := range l.children {
		if l.children[i] != nil {
			l.children[i].Close()
			l.children[i] = nil
		}
	}
	l.children = nil
	l.list = nil
	l.contentImage = nil
	l.Drawable.Close()
}

func (l *ListView) Reset() {
	l.Close()
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
		case *Label, *Button, *Checkbox:
			value.Draw(l.contentImage)
		}
	}
	l.ClearDirty()
}

func (l *ListView) Hit(pt Point[int]) Drawabler {
	if !pt.In(l.rect) || l.IsHidden() {
		return nil
	}
	return l
}

func (l *ListView) MouseDown(md MouseData) {
	l.dragStartY = md.Pos().Y // Сохраняем всегда для корректной обработки кликов в MouseUp

	// Проверяем, попал ли клик в зону полосы прокрутки (последние 12 пикселей справа)
	sbClickZone := 12
	if md.Pos().X >= l.rect.X+l.rect.W-sbClickZone {
		if l.contentRect.H > l.rect.H {
			l.isDraggingScrollbar = true
			l.scrollTo(md.Pos().Y)
			return
		}
	}

	if l.contentRect.H <= l.rect.H {
		return
	}
	l.isDragging = true
	l.dragStartOffset = l.cameraRect.Min.Y
}

func (l *ListView) MouseDrag(md MouseData) {
	if l.isDraggingScrollbar {
		l.scrollTo(md.Pos().Y)
		return
	}
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
}

func (l *ListView) MouseMotion(md MouseData) {
	sbClickZone := 12
	isOver := md.Pos().X >= l.rect.X+l.rect.W-sbClickZone && l.contentRect.H > l.rect.H
	if isOver != l.isHoveringScrollbar {
		l.isHoveringScrollbar = isOver
		l.MarkDirty()
	}
}

func (l *ListView) MouseLeave() {
	if l.isHoveringScrollbar {
		l.isHoveringScrollbar = false
		l.MarkDirty()
	}
}

func (l *ListView) MouseUp(md MouseData) {
	wasDraggingSB := l.isDraggingScrollbar
	l.isDragging = false
	l.isDraggingScrollbar = false

	// Если мы прокручивали список (ползунком или контентом), не обрабатываем клик по элементам
	if wasDraggingSB || md.Pos().Y != l.dragStartY {
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

	// Рисуем полосу прокрутки
	l.drawScrollbar(surface)
}

func (l *ListView) scrollTo(mouseY int) {
	contentH := float64(l.contentRect.H)
	visibleH := float64(l.rect.H)
	if contentH <= visibleH {
		return
	}

	// Вычисляем процент прокрутки на основе положения мыши относительно высоты виджета
	percent := float64(mouseY-l.rect.Y) / visibleH
	newY := int(percent * (contentH - visibleH))

	maxY := int(contentH - visibleH)
	if newY < 0 {
		newY = 0
	} else if newY > maxY {
		newY = maxY
	}
	l.cameraRect = image.Rect(0, newY, l.rect.W, newY+l.rect.H)
	l.MarkDirty()
}

func (l *ListView) drawScrollbar(surface *ebiten.Image) {
	if l.contentRect.H <= l.rect.H {
		return
	}

	// Настройки внешнего вида
	sbWidth := 6.0
	padding := 2.0

	visibleH := float64(l.rect.H)
	contentH := float64(l.contentRect.H)

	// Расчет высоты ползунка (пропорционально видимой области)
	thumbH := (visibleH / contentH) * visibleH
	if thumbH < 20 {
		thumbH = 20 // Минимальная высота ползунка
	}

	// Расчет позиции ползунка
	// Процент прокрутки от 0.0 до 1.0
	scrollPercent := float64(l.cameraRect.Min.Y) / (contentH - visibleH)

	// Доступный путь для перемещения ползунка внутри видимой области
	trackH := visibleH - thumbH
	thumbY := float64(l.rect.Y) + (scrollPercent * trackH)
	thumbX := float64(l.rect.X+l.rect.W) - sbWidth - padding

	// Отрисовка ползунка (используем Fg цвет темы)
	thumbColor := l.Fg()
	if l.isHoveringScrollbar || l.isDraggingScrollbar {
		thumbColor = colornames.White
	}
	vector.FillRect(surface, float32(thumbX), float32(thumbY), float32(sbWidth), float32(thumbH), thumbColor, true)
}

func (l *ListView) SetRect(r Rect[int]) {
	l.Drawable.SetRect(r)
	l.resizeChilds()
	l.cameraRect = image.Rect(0, 0, l.rect.W, l.rect.H)
	l.MarkDirty()
}

func (l *ListView) resizeChilds() {
	if l.rect.IsEmpty() {
		return
	}
	x, y := 0, 0
	w, h := l.rect.W/l.rows, l.itemSize
	if l.longText && len(l.children) > 0 {
		if lbl, ok := l.children[0].(*Label); ok {
			// Вычисляем высоту текста заранее, зная доступную ширину
			fnt := GetUi().FontDefault()
			_, size := fnt.WordWrapText(lbl.Text(), lbl.FontSize(), l.rect.W)
			h = size.Y
			lbl.SetRect(NewRect([]int{0, 0, l.rect.W, h}))
			y = h
		}
	} else {
		row := 0
		col := 1
		for _, v := range l.children {
			v.SetRect(NewRect([]int{x, y, w - 1, h - 1}))
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
	}
	l.contentRect = NewRect([]int{0, 0, l.rect.W, y})
	l.MarkDirty()
}
