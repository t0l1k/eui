package eui

type LayoutBase struct {
	Dirty bool
	rect  Rect
}

func (l *LayoutBase) Layout()           {}
func (l *LayoutBase) GetRect() Rect     { return l.rect }
func (l *LayoutBase) Rect(value Rect)   { l.rect = value; l.Dirty = true }
func (l *LayoutBase) Resize(rect []int) { l.Rect(NewRect(rect)) }
