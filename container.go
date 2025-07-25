package eui

import "github.com/hajimehoshi/ebiten/v2"

type Container struct {
	*Drawable
	container []Drawabler
	layout    Layouter
}

func NewContainer(layout Layouter) *Container {
	return &Container{Drawable: NewDrawable(), layout: layout}
}
func (c *Container) Childrens() []Drawabler { return c.container }
func (c *Container) SetLayout(l Layouter)   { c.layout = l; c.MarkDirty() }
func (c *Container) Add(d Drawabler) *Container {
	c.container = append(c.container, d)
	c.MarkDirty()
	return c
}
func (c *Container) ResetContainer() *Container {
	for _, v := range c.container {
		v.Close()
	}
	c.container = nil
	c.MarkDirty()
	return c
}
func (c *Container) Layout() {
	if c.layout != nil {
		c.layout.Apply(c.container, c.Rect())
		for _, child := range c.Childrens() {
			if child.IsDirty() {
				child.Layout()
			}
		}
	} else {
		c.Traverse(func(d Drawabler) {
			if d.IsDirty() {
				d.Layout()
			}
		}, false)
	}
	c.ClearDirty()
}
func (c *Container) Update(dt int) { c.Traverse(func(d Drawabler) { d.Update(dt) }, false) }
func (c *Container) Draw(surface *ebiten.Image) {
	if c.IsHidden() {
		return
	}
	if c.IsDirty() {
		c.Layout()
	}
	c.Traverse(func(d Drawabler) { d.Draw(surface) }, false)
}
func (c *Container) Traverse(action func(d Drawabler), reverse bool) {
	for _, d := range c.Childrens() {
		traverse(d, action, reverse)
	}
}
func traverse(d Drawabler, action func(d Drawabler), reverse bool) {
	if reverse {
		action(d)
	}
	if container, ok := d.(interface{ Container() []Drawabler }); ok {
		for _, v := range container.Container() {
			traverse(v, action, reverse)
		}
	}
	if !reverse {
		action(d)
	}
}
