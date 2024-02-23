package eui

// Умею размеры виджетов во мне разделить одинаково по горизонтали или по вертикали
type BoxLayout struct {
	LayoutBase
	horizontal bool
	ContainerBase
}

func NewHLayout() *BoxLayout { return &BoxLayout{horizontal: true} }
func NewVLayout() *BoxLayout { return &BoxLayout{horizontal: false} }

func (d *BoxLayout) SetHorizontal() { d.horizontal = true; d.resize() }
func (d *BoxLayout) SetVertical()   { d.horizontal = false; d.resize() }

func (c *BoxLayout) resize() { c.Resize(c.GetRect().GetArr()) }

func (c *BoxLayout) Resize(rect []int) {
	c.Rect(NewRect(rect))
	w0, h0 := c.GetRect().Size()
	x0, y0 := c.GetRect().Pos()
	if count := len(c.GetContainer()); count > 0 {
		if c.horizontal {
			w, h := w0/count, h0
			for i, v := range c.GetContainer() {
				x, y := w*i, 0
				v.Resize([]int{x0 + x, y0 + y, w, h})
			}
		} else {
			w, h := w0, h0/count
			for i, v := range c.GetContainer() {
				x, y := 0, h*i
				v.Resize([]int{x0 + x, y0 + y, w, h})
			}
		}
	}
}
