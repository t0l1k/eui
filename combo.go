package eui

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
)

type ComboBox struct {
	DrawableBase
	btnPlus, btnMinus *Button
	lblValue, lblText *Text
	valueVar          *SubjectBase
	data              []interface{}
	index             int
	onChange          func(c *ComboBox)
}

func NewComboBox(text string, data []interface{}, index int, f func(*ComboBox)) *ComboBox {
	c := &ComboBox{}
	c.SetupCombo(text, data, index, f)
	c.Visible(true)
	return c
}

func (c *ComboBox) SetupCombo(text string, data []interface{}, index int, f func(*ComboBox)) {
	theme := GetUi().theme
	c.Bg(theme.Get(ComboBoxBg))
	c.Fg(theme.Get(ComboBoxFg))
	c.data = data
	c.index = index
	c.onChange = f
	c.lblValue = NewText(c.GetValueString())
	c.lblValue.Bg(c.bg)
	c.lblValue.Fg(c.fg)
	c.valueVar = NewSubject()
	c.valueVar.Attach(c.lblValue)
	c.btnPlus = NewButton("+", func(b *Button) {
		if c.index < len(c.data)-1 {
			c.index++
			c.valueVar.SetValue(c.GetValueString())
			log.Println("combo: +", c.index, c.Value())
		}
		if c.onChange != nil {
			c.onChange(c)
			log.Println("combo: exec +")
		}
	})
	c.btnMinus = NewButton("-", func(b *Button) {
		if c.index > 0 {
			c.index--
			c.valueVar.SetValue(c.GetValueString())
			log.Println("combo: -", c.index, c.Value())
		}
		if c.onChange != nil {
			c.onChange(c)
			log.Println("combo: exec -")
		}
	})
	c.lblText = NewText(text)
	c.lblText.Bg(c.bg)
	c.lblText.Fg(c.fg)
}

func (c *ComboBox) SetText(text string) {
	c.lblText.SetText(text)
}

func (c *ComboBox) Value() interface{} {
	return c.data[c.index]
}

func (c *ComboBox) GetValueString() (result string) {
	switch c.data[c.index].(type) {
	case int:
		result = strconv.Itoa(c.Value().(int))
	case float64:
		result = fmt.Sprintf("%v", c.Value().(float64))
	case string:
		result = fmt.Sprintf("%v", c.Value().(string))
	}
	return result
}

func (c *ComboBox) SetValue(value interface{}) {
	for i, v := range c.data {
		if v == value {
			c.index = i
			break
		}
	}
	c.lblValue.SetText(c.GetValueString())
}

func (c *ComboBox) Layout() {
	c.SpriteBase.Layout()
	c.Dirty = false
}

func (c *ComboBox) Update(dt int) {
	if c.disabled {
		return
	}
	c.btnPlus.Update(dt)
	c.btnMinus.Update(dt)
}

func (b *ComboBox) Draw(surface *ebiten.Image) {
	if !b.IsVisible() {
		return
	}
	if b.Dirty {
		b.Layout()
		b.btnPlus.Layout()
		b.btnMinus.Layout()
		b.lblValue.Layout()
		b.lblText.Layout()
	}
	b.lblValue.Draw(surface)
	b.lblText.Draw(surface)
	b.btnPlus.Draw(surface)
	b.btnMinus.Draw(surface)
}

func (c *ComboBox) Resize(rect []int) {
	c.Rect(NewRect(rect))
	c.SpriteBase.Rect(NewRect(rect))
	x, y, w0, h0 := c.GetRect().GetRect()
	w, h := h0, h0
	c.lblValue.Resize([]int{x, y, w, h})
	x += h0
	h = h0 / 2
	c.btnPlus.Resize([]int{x, y, w, h})
	y += h
	c.btnMinus.Resize([]int{x, y, w, h})
	x += h0
	y -= h
	w, h = w0-h0*2, h0
	c.lblText.Resize([]int{x, y, w, h})
}
