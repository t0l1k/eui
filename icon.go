package eui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Icon struct {
	Image, icon    *ebiten.Image
	rect           *Rect
	Dirty, Visible bool
}

func NewIcon(icon *ebiten.Image, rect []int) *Icon {
	return &Icon{
		icon:    icon,
		rect:    NewRect(rect),
		Dirty:   true,
		Visible: true,
	}
}

func (i *Icon) Layout() {
	w, h := i.rect.Size()
	if i.Image == nil {
		i.Image = ebiten.NewImage(w, h)
	} else {
		i.Image.Clear()
	}
	op := &ebiten.DrawImageOptions{}
	w1, h1 := i.icon.Size()
	var x1, y1 float64
	x1 = float64(w) / float64(w1)
	y1 = float64(h) / float64(h1)
	op.GeoM.Scale(x1, y1)
	i.Image.DrawImage(i.icon, op)
	i.Dirty = false
}

func (i *Icon) GetIcon() *ebiten.Image {
	return i.Image
}

func (i *Icon) SetIcon(icon *ebiten.Image) {
	if i.icon == icon {
		return
	}
	i.icon = icon
	i.Dirty = true
}

func (i *Icon) Size() (int, int) {
	return i.rect.Size()
}

func (i *Icon) Update(dt int) {}

func (i *Icon) Draw(surface *ebiten.Image) {
	if i.Dirty {
		i.Layout()
	}
	if i.Visible {
		op := &ebiten.DrawImageOptions{}
		x, y := i.rect.Pos()
		op.GeoM.Translate(float64(x), float64(y))
		surface.DrawImage(i.Image, op)
	}
}

func (i *Icon) Resize(rect []int) {
	i.rect = NewRect(rect)
	i.Dirty = true
	i.Image = nil
}

func (i *Icon) Close() {
	i.Image.Dispose()
}
