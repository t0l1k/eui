package game

import (
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"golang.org/x/image/colornames"
)

// Умею обновить клетку игры по подписке на состояние клетки поля
type CellIcon struct {
	eui.Drawable
	Btn   *eui.Button
	field *MinedField
}

func NewCellIcon(field *MinedField, f func(b *eui.Button)) *CellIcon {
	c := &CellIcon{}
	c.Btn = eui.NewButton(cellClosed, f)
	c.Setup(field, f)
	return c
}

func (c *CellIcon) Setup(field *MinedField, f func(b *eui.Button)) {
	c.field = field
	// c.Btn.SetupButton(cellClosed, f)
	c.Btn.SetText(cellClosed)
	c.Btn.Bg(colornames.Gray)
	c.Btn.Fg(colornames.Red)
}

func (c *CellIcon) UpdateData(value *cellData) {
	switch value.state {
	case closed:
		c.Btn.SetText(cellClosed)
		c.Btn.Bg(colornames.Gray)
		c.Btn.Fg(colornames.Red)
	case flagged:
		c.Btn.SetText(cellFlagged)
	case questioned:
		c.Btn.SetText(cellQuestioned)
	case firstMined:
		c.Btn.SetText(cellFirstMined)
	case saved:
		c.Btn.Fg(colornames.Yellow)
		c.Btn.SetText(cellSaved)
	case blown:
		c.Btn.SetText(cellBlown)
	case wrongFlagged:
		c.Btn.SetText(cellWrongFlagged)
	case opened:
		cell := c.field.GetCell(value.pos.X, value.pos.Y)
		switch cell.count {
		case 0:
			c.Btn.SetText(cellClosed)
			c.Btn.Bg(colornames.Silver)
		case 1:
			c.Btn.SetText(strconv.Itoa(int(cell.count)))
			c.Btn.Bg(colornames.Silver)
			c.Btn.Fg(colornames.Blue)
		case 2:
			c.Btn.SetText(strconv.Itoa(int(cell.count)))
			c.Btn.Bg(colornames.Silver)
			c.Btn.Fg(colornames.Orange)
		case 3:
			c.Btn.SetText(strconv.Itoa(int(cell.count)))
			c.Btn.Bg(colornames.Silver)
			c.Btn.Fg(colornames.Green)
		case 4:
			c.Btn.SetText(strconv.Itoa(int(cell.count)))
			c.Btn.Bg(colornames.Silver)
			c.Btn.Fg(colornames.Aqua)
		case 5:
			c.Btn.SetText(strconv.Itoa(int(cell.count)))
			c.Btn.Bg(colornames.Silver)
			c.Btn.Fg(colornames.Navy)
		case 6:
			c.Btn.SetText(strconv.Itoa(int(cell.count)))
			c.Btn.Bg(colornames.Silver)
			c.Btn.Fg(colornames.Fuchsia)
		case 7:
			c.Btn.SetText(strconv.Itoa(int(cell.count)))
			c.Btn.Bg(colornames.Silver)
			c.Btn.Fg(colornames.Purple)
		case 8:
			c.Btn.SetText(strconv.Itoa(int(cell.count)))
			c.Btn.Bg(colornames.Silver)
			c.Btn.Fg(colornames.Black)
		}
	}
}

func (c *CellIcon) Update(dt int) {
	c.Btn.Update(dt)
}

func (c *CellIcon) Draw(surface *ebiten.Image) {
	c.Btn.Draw(surface)
}

func (c *CellIcon) Layout() {
}

func (c *CellIcon) SetRect(rect eui.Rect[int]) {
	c.Drawable.SetRect(rect)
	c.Btn.SetRect(rect)
}
