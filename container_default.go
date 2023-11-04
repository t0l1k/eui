package eui

// Умею хранить в контейнере тип интерфейса Drawable. При добавлении устанавливаю себя как родитель
type ContainerBase struct {
	DrawableBase
	parent    Layout
	Container []Drawable
}

func (c *ContainerBase) Add(d Drawable) {
	d.Parent(c)
	c.Container = append(c.Container, d)
}

func (c *ContainerBase) Parent(l Layout) { c.parent = l }
