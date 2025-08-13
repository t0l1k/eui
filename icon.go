package eui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Умею показать изображение пропорционально маштабированное под размеры
type Icon struct {
	*Drawable
	icon *ebiten.Image
}

func NewIcon(icon *ebiten.Image) *Icon { return &Icon{Drawable: NewDrawable(), icon: icon} }
func (i *Icon) Icon() *ebiten.Image    { return i.image }
func (i *Icon) SetIcon(icon *ebiten.Image) {
	if i.icon == icon {
		return
	}
	i.icon = icon
	i.MarkDirty()
}
func (i *Icon) Layout() {
	i.Drawable.Layout()
	w, h := i.rect.Size()
	iconSize := i.icon.Bounds().Size()
	var x1, y1 float64
	x1 = float64(w) / float64(iconSize.X) //w
	y1 = float64(h) / float64(iconSize.Y) //h
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(x1, y1)
	i.Image().DrawImage(i.icon, op)
	i.ClearDirty()
}
