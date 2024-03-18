package eui

// Умею хранить в контейнере тип интерфейса Drawable.
type ContainerBase struct{ container []Drawabler }

func NewContainerBase() *ContainerBase             { return &ContainerBase{} }
func (c *ContainerBase) GetContainer() []Drawabler { return c.container }
func (c *ContainerBase) Add(d Drawabler)           { c.container = append(c.container, d) }

func (c *ContainerBase) ResetContainerBase() {
	for _, v := range c.container {
		v.Close()
	}
	c.container = nil
}
