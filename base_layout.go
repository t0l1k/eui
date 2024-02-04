package eui

type LayoutBase struct {
	Dirty bool
	rect  *Rect
}

func (*LayoutBase) Layout()             {}
func (d *LayoutBase) GetRect() *Rect    { return d.rect }
func (d *LayoutBase) Rect(value *Rect)  { d.rect = value; d.Dirty = true }
func (d *LayoutBase) Resize(rect []int) { d.Rect(NewRect(rect)) }
