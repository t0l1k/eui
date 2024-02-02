package eui

import "github.com/hajimehoshi/ebiten/v2"

type Drawable interface {
	Layout
	Sprite
	Close()
}

type Sprite interface {
	Update(dt int)
	Draw(*ebiten.Image)
}

type Layout interface {
	Layout()
	Resize([]int)
}

type Container interface {
	Add(Drawable)
	GetContainer() []Drawable
}

type Scene interface {
	Sprite
	Entered()
	Resize()
	Quit()
}

type Observer interface {
	UpdateData(interface{})
}

type Input interface {
	UpdateInput(interface{})
}
