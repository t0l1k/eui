package eui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type ListView struct {
	View
	list               []string
	itemSize           int
	contentRect        *Rect
	contentImage       *ebiten.Image
	offset, lastOffset int
	cameraRect         image.Rectangle
}

func NewListView(list []string, itemSize int) *ListView {
	l := &ListView{}
	l.SetupListView(list, itemSize)
	return l
}

func (l *ListView) SetupListView(list []string, itemSize int) {
	theme := GetUi().theme
	l.SetupView()
	l.itemSize = itemSize
	l.list = list
	for _, v := range list {
		lbl := NewText(v)
		lbl.Bg(theme.Get(ListViewItemBg))
		lbl.Fg(theme.Get(ListViewItemFg))
		l.Add(lbl)
	}
}

func (l *ListView) Reset() {
	l.list = nil
	l.Container = l.Container[:0]
	l.contentImage = nil
	l.offset = 0
	l.lastOffset = 0
	l.cameraRect = image.Rect(0, 0, l.rect.W, l.rect.H)
	l.dirty = true
}

func (l *ListView) Add(d Drawable) {
	l.ContainerBase.Add(d)
	if l.rect != nil {
		l.resizeChilds()
	}
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
	for _, v := range l.Container {
		v.Draw(l.contentImage)
	}
	l.dirty = false
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
}

func (l *ListView) Draw(surface *ebiten.Image) {
	if !l.visible || l.disabled {
		return
	}
	if l.dirty {
		l.Layout()
		for _, c := range l.Container {
			c.Layout()
		}
	}
	op := &ebiten.DrawImageOptions{}
	x0, y0 := l.rect.Pos()
	op.GeoM.Translate(float64(x0), float64(y0))
	surface.DrawImage(l.contentImage.SubImage(l.cameraRect).(*ebiten.Image), op)
}

func (l *ListView) Resize(r []int) {
	l.Rect(r)
	l.resizeChilds()
	l.cameraRect = image.Rect(0, 0, l.rect.W, l.rect.H)
}

func (l *ListView) resizeChilds() {
	x, y := 0, 0
	w, h := l.rect.W, l.itemSize
	for _, v := range l.Container {
		v.(*Text).Resize([]int{x, y, w, h})
		y += h
	}
	if y == 0 {
		y = l.rect.H
	}
	l.contentRect = NewRect([]int{0, 0, w, y})
	l.dirty = true
}
