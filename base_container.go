package eui

// Умею хранить в контейнере тип интерфейса Drawable.
type ContainerBase struct{ container []Drawabler }

func (c *ContainerBase) GetContainer() []Drawabler { return c.container }
func (c *ContainerBase) ResetContainerBase()       { c.container = nil }
func (c *ContainerBase) Add(d Drawabler)           { c.container = append(c.container, d) }
