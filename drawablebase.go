package eui

import "github.com/hajimehoshi/ebiten/v2"

type DrawableBase struct {
	Dirty, Visible bool
	Rect           *Rect
	Image          *ebiten.Image
}

func (d *DrawableBase) Update(dt int)              {}
func (d *DrawableBase) Draw(surface *ebiten.Image) {}
func (d *DrawableBase) Layout()                    {}
func (d *DrawableBase) Resize(rect []int)          {}
func (d *DrawableBase) Close()                     {}
