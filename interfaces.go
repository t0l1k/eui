package eui

import "github.com/hajimehoshi/ebiten/v2"

type Drawable interface {
	Layout
	Sprite
}

type Sprite interface {
	Update(dt int)
	Draw(*ebiten.Image)
}

type Layout interface {
	GetName() string
	Parent(Layout)
	Layout()
	Resize([]int)
	Close()
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
