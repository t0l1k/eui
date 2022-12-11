package eui

type ContainerDefault struct {
	Container []Sprite
}

func (c *ContainerDefault) Add(d Sprite) {
	c.Container = append(c.Container, d)
}
