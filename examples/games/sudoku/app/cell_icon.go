package app

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/examples/games/sudoku/game"
	"golang.org/x/image/colornames"
)

type CellIcon struct {
	*eui.Container
	cell      *game.Cell
	dim       game.Dim
	btn       *eui.Button
	layout    *eui.Container
	showNotes bool
	f         func(b *eui.Button)
	highlight int
}

func NewCellIcon(dim game.Dim, cell *game.Cell, f func(b *eui.Button), bg, fg color.RGBA) *CellIcon {
	c := &CellIcon{Container: eui.NewContainer(eui.NewAbsoluteLayout())}
	c.cell = cell
	c.dim = dim
	c.layout = eui.NewContainer(eui.NewGridLayout(1, 1, 1))
	c.f = f
	c.btn = eui.NewButton("0", f)
	c.Add(c.btn)
	c.Bg(bg)
	if c.cell.IsReadOnly() {
		c.Fg(colornames.Blue)
	} else {
		c.Fg(fg)
	}
	c.Add(c.layout)
	return c
}

func (c *CellIcon) ShowNotes(value bool) { c.showNotes = value; c.MarkDirty() }
func (c *CellIcon) Highlight(value int)  { c.highlight = value; c.MarkDirty() }
func (c *CellIcon) UpdateData(value int) { c.MarkDirty() }

func (c *CellIcon) Layout() {
	c.Drawable.Layout()
	c.Image().Fill(c.GetBg())
	c.layout.ResetContainer()
	value := c.cell.GetValue()
	if value > 0 {
		lbl := eui.NewLabel(strconv.Itoa(value))
		c.layout.Add(lbl)
		c.layout.SetLayout(eui.NewGridLayout(1, 1, 1))
		defer lbl.Close()
		if value == c.highlight {
			lbl.Bg(colornames.Yellow)
		} else {
			lbl.Bg(colornames.Silver)
		}
		if c.cell.IsReadOnly() {
			lbl.Fg(colornames.Blue)
		} else {
			lbl.Fg(c.GetFg())
		}
		// log.Println("Иконка с цифрой", c.cell.GetValue())
	} else {
		notes := c.cell.GetNotes()
		if c.showNotes && len(notes) > 0 {
			for i := 0; i < c.dim.Size(); i++ {
				lbl := eui.NewLabel("")
				lbl.Bg(colornames.Silver)
				lbl.Fg(c.GetFg())
				c.layout.Add(lbl)
				found := notes.IsContain(i + 1)
				if found {
					idx, _ := notes.Index(i + 1)
					lbl.SetText(strconv.Itoa(notes[idx]))
					if i+1 == c.highlight {
						lbl.Bg(colornames.Yellow)
					}
				} else {
					lbl.SetText("")
				}
			}
			c.layout.SetLayout(eui.NewGridLayout(c.dim.W, c.dim.H, 1))
			// log.Println("Иконка с заметкой", arr1)
		} else {
			lbl := eui.NewLabel("")
			c.layout.Add(lbl)
			c.layout.SetLayout(eui.NewGridLayout(1, 1, 1))
			defer lbl.Close()
			if len(notes) == 0 {
				lbl.Bg(colornames.Orange)
			} else {
				lbl.Bg(colornames.Silver)
			}
			lbl.Fg(c.GetFg())
			// log.Println("Иконка без заметок", c.cell.GetValue())
		}
	}
	c.ClearDirty()
}

func (d *CellIcon) Update(dt int) {
	if d.IsHidden() {
		return
	}
	d.btn.Update(dt)
}

func (d *CellIcon) Draw(surface *ebiten.Image) {
	if d.IsHidden() {
		return
	}
	if d.IsDirty() {
		d.Layout()
	}
	d.Container.Draw(surface)
}

func (c *CellIcon) SetRect(rect eui.Rect[int]) {
	c.Container.SetRect(rect)
	c.btn.SetRect(rect)
	c.layout.SetRect(rect)
	c.ImageReset()
}
