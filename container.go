package eui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Container struct {
	*Drawable
	children []Drawabler
	layout   Layouter
}

func NewContainer(layout Layouter) *Container {
	return &Container{Drawable: NewDrawable(), layout: layout}
}
func (c *Container) Children() []Drawabler { return c.children }
func (c *Container) SetLayout(l Layouter)  { c.layout = l; c.MarkDirty() }
func (c *Container) Add(d Drawabler) *Container {
	c.children = append(c.children, d)
	c.MarkDirty()
	return c
}

// ResetContainer после обнуление контейнера, обязательно вызвать func (*Container)Layout()
func (c *Container) ResetContainer() *Container {
	c.Traverse(func(d Drawabler) { d.Close() }, false)
	c.children = nil
	c.MarkDirty()
	return c
}
func (c *Container) Layout() {
	c.Drawable.SetRect(c.Rect())
	c.layout.Apply(c.children, c.Rect())
	if c.layout != nil {
		for _, child := range c.Children() {
			if child.IsDirty() {
				child.Layout()
			}
		}
	} else {
		panic("Layout nil")
	}
	c.ClearDirty()
}
func (c *Container) Update() { c.Traverse(func(d Drawabler) { d.Update() }, false) }
func (c *Container) Show() {
	c.Traverse(func(d Drawabler) { d.Show() }, false)
	c.SetState(StateNormal)
}
func (c *Container) Hide() {
	c.Traverse(func(d Drawabler) { d.Hide() }, false)
	c.SetState(StateHidden)
}
func (c *Container) Draw(surface *ebiten.Image) {
	if c.IsHidden() {
		return
	}
	if c.IsDirty() {
		c.Layout()
	}
	c.Traverse(func(d Drawabler) {
		if d.ViewType().IsBackground() {
			d.Draw(surface)
		}
	}, false)
	c.Traverse(func(d Drawabler) {
		if !d.ViewType().IsBackground() {
			d.Draw(surface)
		}
	}, false)
}
func (c *Container) Traverse(action func(d Drawabler), reverse bool) {
	for _, d := range c.Children() {
		traverse(d, action, reverse)
	}
}
func traverse(d Drawabler, action func(d Drawabler), reverse bool) {
	if reverse {
		action(d)
	}
	if container, ok := d.(interface{ Children() []Drawabler }); ok {
		for _, v := range container.Children() {
			traverse(v, action, reverse)
		}
	}
	if !reverse {
		action(d)
	}
}

// func (c *Container) Hit(value Point[int]) Drawabler {
// 	if !value.In(c.rect) || c.state.IsHidden() {
// 		return nil
// 	}
// 	for i := len(c.Children()) - 1; i >= 0; i-- {
// 		c := c.Children()[i]
// 		if mh, ok := c.(interface{ Hit(Point[int]) Drawabler }); ok {
// 			if hit := mh.Hit(value); hit != nil {
// 				return hit
// 			}
// 		}

// 	}
// 	return c
// }

type ContainerVisibleByFilter struct {
	*Container
	items  *Container
	filter func(Drawabler, ViewState) bool
}

func NewContainerVisibleByFilter(l Layouter, items *Container, fn func(Drawabler, ViewState) bool) *ContainerVisibleByFilter {
	return &ContainerVisibleByFilter{Container: NewContainer(l), items: items, filter: fn}
}
func (c *ContainerVisibleByFilter) UpdateBy(s ViewState) {
	c.ResetContainer()
	for _, item := range c.items.Children() {
		item.SetState(StateHidden)
		match := c.filter(item, s)
		if match {
			item.SetState(StateNormal)
			c.Add(item)
		}
	}
	c.MarkDirty()
	if c.Rect().IsEmpty() {
		return
	}
	c.Layout()
}
