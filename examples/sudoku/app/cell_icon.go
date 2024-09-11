package app

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/colors"
	"github.com/t0l1k/eui/examples/sudoku/game"
)

type CellIcon struct {
	eui.DrawableBase
	cell            *game.Cell
	dim             game.Dim
	btn             *eui.Button
	layout          *eui.GridLayoutRightDown
	show, showNotes bool
	f               func(b *eui.Button)
	highlight       int
}

func NewCellIcon(dim game.Dim, cell *game.Cell, f func(b *eui.Button), bg, fg color.RGBA) *CellIcon {
	c := &CellIcon{}
	c.cell = cell
	c.dim = dim
	c.layout = eui.NewGridLayoutRightDown(1, 1)
	c.f = f
	c.btn = eui.NewButton("0", f)
	c.Add(c.btn)
	c.Bg(bg)
	if c.cell.IsReadOnly() {
		c.Fg(colors.Blue)
	} else {
		c.Fg(fg)
	}
	c.Visible(true)
	return c
}

func (c *CellIcon) ShowNotes(value bool)         { c.showNotes = value; c.Dirty = true }
func (c *CellIcon) Highlight(value int)          { c.highlight = value; c.Dirty = true }
func (c *CellIcon) UpdateData(value interface{}) { c.Dirty = true }

func (c *CellIcon) Layout() {
	c.SpriteBase.Layout()
	c.Image().Fill(c.GetBg())
	c.layout.ResetContainerBase()
	value := c.cell.GetValue()
	if value > 0 {
		lbl := eui.NewText(strconv.Itoa(value))
		c.layout.Add(lbl)
		c.layout.SetDim(1, 1)
		defer lbl.Close()
		if value == c.highlight {
			lbl.Bg(colors.Yellow)
		} else {
			lbl.Bg(colors.Silver)
		}
		if c.cell.IsReadOnly() {
			lbl.Fg(colors.Blue)
		} else {
			lbl.Fg(c.GetFg())
		}
		// log.Println("Иконка с цифрой", c.cell.GetValue())
	} else {
		notes := c.cell.GetNotes()
		if c.showNotes && len(notes) > 0 {
			for i := 0; i < c.dim.Size(); i++ {
				lbl := eui.NewText("")
				lbl.Bg(colors.Silver)
				lbl.Fg(c.GetFg())
				c.layout.Add(lbl)
				found := notes.IsContain(i + 1)
				if found {
					idx, _ := notes.Index(i + 1)
					lbl.SetText(strconv.Itoa(notes[idx]))
					if i+1 == c.highlight {
						lbl.Bg(colors.Yellow)
					}
				} else {
					lbl.SetText("")
				}
			}
			c.layout.SetDim(float64(c.dim.W), float64(c.dim.H))
			// log.Println("Иконка с заметкой", arr1)
		} else {
			lbl := eui.NewText("")
			c.layout.Add(lbl)
			c.layout.SetDim(1, 1)
			defer lbl.Close()
			if len(notes) == 0 {
				lbl.Bg(colors.Orange)
			} else {
				lbl.Bg(colors.Silver)
			}
			lbl.Fg(c.GetFg())
			// log.Println("Иконка без заметок", c.cell.GetValue())
		}
	}
	c.Dirty = false
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
