package eui

import (
	"fmt"
	"image/color"
	"log"
	"strconv"
)

type ComboBox struct {
	View
	btnPlus, btnMinus *Button
	lblValue, lblText *Text
	valueVar          *StringVar
	data              []interface{}
	index             int
	onChange          func(c *ComboBox)
}

func NewComboBox(text string, data []interface{}, index int, bg, fg color.Color, f func(*ComboBox)) *ComboBox {
	c := &ComboBox{}
	c.SetupView(bg)
	c.SetupCombo(text, data, index, bg, fg, f)
	return c
}

func (c *ComboBox) SetupCombo(text string, data []interface{}, index int, bg, fg color.Color, f func(*ComboBox)) {
	c.data = data
	c.index = index
	c.onChange = f
	c.lblValue = NewText(c.GetValueString(), bg, fg)
	c.valueVar = NewStringVar(c.GetValueString())
	c.valueVar.Attach(c.lblValue)
	c.btnPlus = NewButton("+", bg, fg, func(b *Button) {
		if c.index < len(c.data)-1 {
			c.index++
			c.valueVar.Set(c.GetValueString())
			log.Println("combo: +", c.index, c.Value())
		}
		if c.onChange != nil {
			c.onChange(c)
			log.Println("combo: exec +")
		}
	})
	c.btnMinus = NewButton("-", bg, fg, func(b *Button) {
		if c.index > 0 {
			c.index--
			c.valueVar.Set(c.GetValueString())
			log.Println("combo: -", c.index, c.Value())
		}
		if c.onChange != nil {
			c.onChange(c)
			log.Println("combo: exec -")
		}
	})
	c.lblText = NewText(text, bg, fg)
	c.Add(c.lblValue)
	c.Add(c.btnPlus)
	c.Add(c.btnMinus)
	c.Add(c.lblText)
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
			c.Dirty(true)
			break
		}
	}
}

func (c *ComboBox) Resize(r []int) {
	c.Rect(r)
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
	w, h = w0-x, h0
	c.lblText.Resize([]int{x, y, w, h})
}
