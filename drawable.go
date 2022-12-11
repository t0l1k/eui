package eui

import "github.com/hajimehoshi/ebiten/v2"

type Drawable interface {
	Layout()
	Sprite
	Resize([]int)
}

type Sprite interface {
	Update(dt int)
	Draw(surface *ebiten.Image)
	Close()
}
