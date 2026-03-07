package circle

import (
	"fmt"
	"log"

	"github.com/t0l1k/eui"
)

type Canvas struct {
	*eui.Container
	clicked *eui.Signal[eui.Point[int]]
}

func NewCanvas() *Canvas {
	c := &Canvas{
		Container: eui.NewContainer(eui.NewAbsoluteLayout()),
		clicked:   eui.NewSignal(func(a, b eui.Point[int]) bool { return a.Eq(b) }),
	}
	return c
}
func (*Canvas) WantBlur() bool { return true }
func (c *Canvas) Hit(pt eui.Point[int]) eui.Drawabler {
	if !pt.In(c.Rect()) || c.IsHidden() || c.IsDisabled() {
		return nil
	}
	return c
}
func (c *Canvas) MouseUp(md eui.MouseData) {
	c.clicked.Emit(md.Pos())
	c.MarkDirty()
	log.Println("Canvas:MouseUp", md.Pos(), c.Rect())
}

func (c *Canvas) String() string {
	return fmt.Sprintf("Canvas:%v[%v]", len(c.Children()), c.Children())
}
