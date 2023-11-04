package eui

// Умею размеры виджетов во мне разделить одинаково по горизонтали или по вертикали
type BoxLayout struct {
	horizontal bool
	ContainerBase
}

func NewHLayout() *BoxLayout { return &BoxLayout{horizontal: true} }
func NewVLayout() *BoxLayout { return &BoxLayout{horizontal: false} }

func (d *BoxLayout) SetHorizontal()  { d.horizontal = true }
func (d *BoxLayout) SetVertical()    { d.horizontal = false }
func (d *BoxLayout) GetName() string { return "BoxLayout" }

func (c *BoxLayout) Resize(rect []int) {
	crect := NewRect(rect)
	w0, h0 := crect.Size()
	x0, y0 := crect.Pos()
	if count := len(c.Container); count > 0 {
		if c.horizontal {
			w, h := w0/count, h0
			for i, v := range c.Container {
				x, y := w*i, 0
				v.Resize([]int{x0 + x, y0 + y, w, h})
			}
		} else {
			w, h := w0, h0/count
			for i, v := range c.Container {
				x, y := 0, h*i
				v.Resize([]int{x0 + x, y0 + y, w, h})
			}
		}
	}
}
