package eui

// Умею хранить в контейнере тип интерфейса Drawable. При добавлении устанавливаю себя как родитель
type ContainerBase struct {
	DrawableBase
	Container []Drawable
}

func (c *ContainerBase) Add(d Drawable) {
	c.Container = append(c.Container, d)
}
