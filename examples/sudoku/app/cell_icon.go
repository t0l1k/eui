package app

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/sudoku/game"
)

type CellIcon struct {
	eui.DrawableBase
	cell      *game.Cell
	btn       *eui.Button
	layout    *eui.GridLayoutRightDown
	show      bool
	f         func(b *eui.Button)
	highlight int
}

func NewCellIcon(cell *game.Cell, f func(b *eui.Button), bg, fg color.RGBA) *CellIcon {
	c := &CellIcon{}
	c.cell = cell
	c.layout = eui.NewGridLayoutRightDown(float64(cell.GetDim()), float64(cell.GetDim()))
	c.f = f
	c.btn = eui.NewButton("-99", f)
	c.Add(c.btn)
	c.Bg(bg)
	c.Fg(fg)
	c.setup()
	c.Visible(true)
	return c
}

func (d *CellIcon) setup() {}

func (c *CellIcon) Highlight(value int) {
	c.highlight = value
	c.Dirty = true
}

func (c *CellIcon) Layout() {
	c.SpriteBase.Layout()
	c.Image().Fill(c.GetBg())
	c.layout.ResetContainerBase()
	if c.cell.GetValue() > 0 {
		lbl := eui.NewText(strconv.Itoa(c.cell.GetValue()))
		c.layout.Add(lbl)
		c.layout.SetDim(1, 1)
		defer lbl.Close()
		if c.cell.GetValue() == c.highlight {
			lbl.Bg(eui.Yellow)
		} else {
			lbl.Bg(eui.Silver)
		}
		lbl.Fg(c.GetFg())
		fmt.Println("Иконка с цифрой", c.cell.GetValue())
	} else {
		size := c.cell.GetDim()
		arr1, _, _ := c.cell.GetNotes()
		if len(arr1) > 0 {
			for i := 0; i < size*size; i++ {
				lbl := eui.NewText("")
				lbl.Bg(eui.Silver)
				lbl.Fg(c.GetFg())
				c.layout.Add(lbl)
				found := eui.IntSliceContains(arr1, i+1)
				if found {
					idx := eui.GetIdxValueFromIntSlice(arr1, i+1)
					lbl.SetText(strconv.Itoa(arr1[idx]))
					if i+1 == c.highlight {
						lbl.Bg(eui.Yellow)
					}
				} else {
					lbl.SetText("")
				}
			}
			c.layout.SetDim(float64(size), float64(size))
			fmt.Println("Иконка с заметкой", arr1)
		} else {
			lbl := eui.NewText("")
			c.layout.Add(lbl)
			c.layout.SetDim(1, 1)
			defer lbl.Close()
			lbl.Bg(eui.Red)
			lbl.Fg(c.GetFg())
			fmt.Println("Иконка без заметок", c.cell.GetValue())
		}
	}
	c.Dirty = false
}

func (c *CellIcon) UpdateData(value interface{}) {
	switch v := value.(type) {
	case int:
		c.Dirty = true
		fmt.Println("cell icon get", v)
	}
}

func (d *CellIcon) Update(dt int) {
	if !d.IsVisible() {
		return
	}
	d.btn.Update(dt)
}

func (c *CellIcon) Draw(surface *ebiten.Image) {
	if !c.IsVisible() {
		return
	}
	if c.Dirty {
		c.Layout()
	}
	for _, v := range c.layout.GetContainer() {
		v.Draw(surface)
	}
}

func (c *CellIcon) IsVisible() bool    { return c.show }
func (c *CellIcon) Visible(value bool) { c.show = value }

func (c *CellIcon) Resize(rect []int) {
	c.Rect(eui.NewRect(rect))
	c.SpriteBase.Resize(rect)
	c.btn.Resize(rect)
	c.layout.Resize(rect)
	c.ImageReset()
}
