package eui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Умею показать изображение пропорционально маштабированное под размеры
type Icon struct {
	View
	icon *ebiten.Image
}

func NewIcon(icon *ebiten.Image) *Icon {
	i := &Icon{}
	i.SetupIcon(icon)
	return i
}

func (i *Icon) SetupIcon(icon *ebiten.Image) {
	i.icon = icon
	i.SetupView(Black)
	i.Name("icon")
}

func (i *Icon) Layout() {
	i.View.Layout()
	w, h := i.rect.Size()
	op := &ebiten.DrawImageOptions{}
	w1, h1 := i.icon.Size()
	var x1, y1 float64
	x1 = float64(w) / float64(w1)
	y1 = float64(h) / float64(h1)
	op.GeoM.Scale(x1, y1)
	i.image.DrawImage(i.icon, op)
	i.dirty = false
}

func (i *Icon) GetIcon() *ebiten.Image {
	return i.image
}

func (i *Icon) SetIcon(icon *ebiten.Image) {
	if i.icon == icon {
		return
	}
	i.icon = icon
	i.dirty = true
}

func (i *Icon) Draw(surface *ebiten.Image) {
	if i.dirty {
		i.Layout()
	}
	op := &ebiten.DrawImageOptions{}
	x, y := i.rect.Pos()
	op.GeoM.Translate(float64(x), float64(y))
	surface.DrawImage(i.image, op)
}
