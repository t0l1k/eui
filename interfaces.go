package eui

import "github.com/hajimehoshi/ebiten/v2"

type Drawabler interface {
	IsVisible() bool
	Visible(bool)
	Enable()
	Disable()
	Layouter
	Spriter
	Close()
}

type Spriter interface {
	Update(dt int)
	Draw(*ebiten.Image)
}

type Layouter interface {
	Layout()
	Resize([]int)
}

type Containerer interface {
	Add(Drawabler)
	GetContainer() []Layouter
}

type Sceneer interface {
	Spriter
	Entered()
	Resize()
	Quit()
}

type Inputer interface {
	UpdateInput(interface{})
}
