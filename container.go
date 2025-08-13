package eui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Container struct {
	*Drawable
	container []Drawabler
	layout    Layouter
}

func NewContainer(layout Layouter) *Container {
	return &Container{Drawable: NewDrawable(), layout: layout}
}
func (c *Container) Children() []Drawabler { return c.container }
func (c *Container) SetLayout(l Layouter)  { c.layout = l; c.MarkDirty() }
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
		c.Drawable.SetRect(c.Rect())
		c.layout.Apply(c.container, c.Rect())
		for _, child := range c.Children() {
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
func (c *Container) Show()         { c.Traverse(func(d Drawabler) { d.Show() }, false) }
func (c *Container) Hide()         { c.Traverse(func(d Drawabler) { d.Hide() }, false) }
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
			// log.Println("traverse:check:", v.Rect())
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
