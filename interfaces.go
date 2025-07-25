package eui

import "github.com/hajimehoshi/ebiten/v2"

type Drawabler interface {
	Spriter
	IsVisible() bool
	Visible(bool)
	Enable()
	Disable()
	Layout()
	Rect() Rect[int]
	SetRect(Rect[int])
	Resize(Rect[int])
	IsDirty() bool
	MarkDirty()
	ClearDirty()
	Traverse(func(d Drawabler), bool)
	Close()
}

type Spriter interface {
	Update(dt int)
	Draw(*ebiten.Image)
}

type Layouter interface {
	Apply([]Drawabler, Rect[int])
}

type Containerer interface {
	Drawabler
	Add(Drawabler)
	Childrens() []Drawabler
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
