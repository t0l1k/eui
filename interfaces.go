package eui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Drawabler interface {
	Spriter
	State() ViewState
	SetState(ViewState)
	IsHidden() bool
	Hide()
	Show()
	IsDisabled() bool
	Enable()
	Disable()
	Layout()
	Rect() Rect[int]
	SetRect(Rect[int])
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
	Children() []Drawabler
}

type Sceneer interface {
	Drawabler
	Entered()
	Quit()
}

type KeybordHandlerer interface {
	KeyPressed(KeyboardData)
	KeyReleased(KeyboardData)
}

type MouseHandlerer interface {
	Hit(Point[int]) Drawable
	MouseDown(MouseData)
	MouseUp(MouseData)
	MouseMotion(MouseData)
	MouseDrag(MouseData)
	MouseWheel(MouseData)
	MouseEnter()
	MouseLeave()
	WantBlur() bool
}
