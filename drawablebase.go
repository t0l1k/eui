package eui

import "github.com/hajimehoshi/ebiten/v2"

type DrawableBase struct {
	name  string
	Dirty bool
	Rect  *Rect
}

func (d *DrawableBase) Update(dt int)              {}
func (d *DrawableBase) Draw(surface *ebiten.Image) {}
func (d *DrawableBase) GetName() string            { return d.name }
func (d *DrawableBase) Parent(Layout)              {}
func (d *DrawableBase) Layout()                    {}
func (d *DrawableBase) Resize(rect []int)          {}
func (d *DrawableBase) Close()                     {}
