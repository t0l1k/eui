package game

import (
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/t0l1k/eui"
	"github.com/t0l1k/eui/colors"
)

// Умею обновить клетку игры по подписке на состояние клетки поля
type CellIcon struct {
	eui.DrawableBase
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
	c.Btn.SetupButton(cellClosed, f)
	c.Btn.Bg(colors.Gray)
	c.Btn.Fg(colors.Red)
}

func (c *CellIcon) UpdateData(value interface{}) {
	switch v := value.(type) {
	case *cellData:
		switch v.state {
		case closed:
			c.Btn.SetText(cellClosed)
			c.Btn.Bg(colors.Gray)
			c.Btn.Fg(colors.Red)
		case flagged:
			c.Btn.SetText(cellFlagged)
		case questioned:
			c.Btn.SetText(cellQuestioned)
		case firstMined:
			c.Btn.SetText(cellFirstMined)
		case saved:
			c.Btn.Fg(colors.Yellow)
			c.Btn.SetText(cellSaved)
		case blown:
			c.Btn.SetText(cellBlown)
		case wrongFlagged:
			c.Btn.SetText(cellWrongFlagged)
		case opened:
			cell := c.field.GetCell(v.pos.X, v.pos.Y)
			switch cell.count {
			case 0:
				c.Btn.SetText(cellClosed)
				c.Btn.Bg(colors.Silver)
			case 1:
				c.Btn.SetText(strconv.Itoa(int(cell.count)))
				c.Btn.Bg(colors.Silver)
				c.Btn.Fg(colors.Blue)
			case 2:
				c.Btn.SetText(strconv.Itoa(int(cell.count)))
				c.Btn.Bg(colors.Silver)
				c.Btn.Fg(colors.Orange)
			case 3:
				c.Btn.SetText(strconv.Itoa(int(cell.count)))
				c.Btn.Bg(colors.Silver)
				c.Btn.Fg(colors.Green)
			case 4:
				c.Btn.SetText(strconv.Itoa(int(cell.count)))
				c.Btn.Bg(colors.Silver)
				c.Btn.Fg(colors.Aqua)
			case 5:
				c.Btn.SetText(strconv.Itoa(int(cell.count)))
				c.Btn.Bg(colors.Silver)
				c.Btn.Fg(colors.Navy)
			case 6:
				c.Btn.SetText(strconv.Itoa(int(cell.count)))
				c.Btn.Bg(colors.Silver)
				c.Btn.Fg(colors.Fuchsia)
			case 7:
				c.Btn.SetText(strconv.Itoa(int(cell.count)))
				c.Btn.Bg(colors.Silver)
				c.Btn.Fg(colors.Purple)
			case 8:
				c.Btn.SetText(strconv.Itoa(int(cell.count)))
				c.Btn.Bg(colors.Silver)
				c.Btn.Fg(colors.Black)
			}
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

func (c *CellIcon) Resize(rect []int) {
	c.Rect(eui.NewRect(rect))
	c.Btn.Resize(rect)
}
