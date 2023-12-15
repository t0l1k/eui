package game

import (
	"strconv"

	"github.com/t0l1k/eui"
)

// Умею обновить клетку игры по подписке на состояние клетки поля
type CellIcon struct {
	eui.View
	Btn   *eui.Button
	field *MinedField
}

func NewCellIcon(field *MinedField, f func(b *eui.Button)) *CellIcon {
	c := &CellIcon{}
	c.SetupView()
	c.Btn = eui.NewButton(cellClosed, f)
	c.Add(c.Btn)
	c.Setup(field, f)
	return c
}

func (c *CellIcon) Setup(field *MinedField, f func(b *eui.Button)) {
	c.field = field
	c.Btn.SetupButton(cellClosed, f)
	c.Btn.Bg(eui.Gray)
	c.Btn.Fg(eui.Red)
}

func (c *CellIcon) UpdateData(value interface{}) {
	switch v := value.(type) {
	case *cellData:
		switch v.state {
		case closed:
			c.Btn.SetText(cellClosed)
			c.Btn.Bg(eui.Gray)
			c.Btn.Fg(eui.Red)
		case flagged:
			c.Btn.SetText(cellFlagged)
		case questioned:
			c.Btn.SetText(cellQuestioned)
		case firstMined:
			c.Btn.SetText(cellFirstMined)
		case saved:
			c.Btn.Fg(eui.Yellow)
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
				c.Btn.Bg(eui.Silver)
			case 1:
				c.Btn.SetText(strconv.Itoa(int(cell.count)))
				c.Btn.Bg(eui.Silver)
				c.Btn.Fg(eui.Blue)
			case 2:
				c.Btn.SetText(strconv.Itoa(int(cell.count)))
				c.Btn.Bg(eui.Silver)
				c.Btn.Fg(eui.Orange)
			case 3:
				c.Btn.SetText(strconv.Itoa(int(cell.count)))
				c.Btn.Bg(eui.Silver)
				c.Btn.Fg(eui.Green)
			case 4:
				c.Btn.SetText(strconv.Itoa(int(cell.count)))
				c.Btn.Bg(eui.Silver)
				c.Btn.Fg(eui.Aqua)
			case 5:
				c.Btn.SetText(strconv.Itoa(int(cell.count)))
				c.Btn.Bg(eui.Silver)
				c.Btn.Fg(eui.Navy)
			case 6:
				c.Btn.SetText(strconv.Itoa(int(cell.count)))
				c.Btn.Bg(eui.Silver)
				c.Btn.Fg(eui.Fuchsia)
			case 7:
				c.Btn.SetText(strconv.Itoa(int(cell.count)))
				c.Btn.Bg(eui.Silver)
				c.Btn.Fg(eui.Purple)
			case 8:
				c.Btn.SetText(strconv.Itoa(int(cell.count)))
				c.Btn.Bg(eui.Silver)
				c.Btn.Fg(eui.Black)
			}
		}
	}
}
