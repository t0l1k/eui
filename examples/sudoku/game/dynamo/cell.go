package dynamo

import (
	"fmt"

	"github.com/t0l1k/eui"
)

type cell struct {
	value, dim, size int
	notes            []int
}

func newCell(dim, size int) *cell {
	c := &cell{dim: dim, size: size}
	c.reset()
	return c
}

func (c *cell) reset() {
	c.notes = nil
	c.value = -1
	for i := 0; i < c.size; i++ {
		c.notes = append(c.notes, i+1)
	}
}

func (c *cell) add(value int) {
	c.value = value
	c.notes = eui.RemoveFromIntSliceValue(c.notes, value)
}

func (c *cell) removeNote(value int) {
	c.notes = eui.RemoveFromIntSliceValue(c.notes, value)
}

func (c *cell) getNotes() []int { return c.notes }

func (c cell) String() (result string) {
	if c.value > 0 {
		result = fmt.Sprintf("[%3v]", c.value)
	} else {
		result += "["
		for _, v := range c.notes {
			result += fmt.Sprintf("(%v)", v)
		}
		result += "]"
	}
	return result
}
