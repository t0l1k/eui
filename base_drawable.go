package eui

import "github.com/hajimehoshi/ebiten/v2"

type DrawableBase struct {
	SpriteBase
	ContainerBase
}

func (d *DrawableBase) Update(dt int) {
	for _, v := range d.GetContainer() {
		v.Update(dt)
	}
}

func (d *DrawableBase) Draw(surface *ebiten.Image) {
	for _, v := range d.GetContainer() {
		v.Draw(surface)
	}
}

func (d *DrawableBase) Resize(rect []int) {
	for _, v := range d.GetContainer() {
		v.Resize(rect)
	}
}

func (d *DrawableBase) Close() {
	for _, v := range d.GetContainer() {
		v.Close()
	}
}
