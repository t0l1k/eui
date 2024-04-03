package recur

import (
	"fmt"
)

type Cell struct {
	dim   int
	value int
	notes map[int]bool
}

func NewCell(dim int) *Cell {
	c := &Cell{dim: dim}
	c.Reset()
	return c
}

func (c *Cell) Reset() {
	c.value = 0
	c.notes = make(map[int]bool)
	for i := 0; i < c.Size(); i++ {
		c.notes[i+1] = true
	}
}

func (c *Cell) Add(value int) {
	c.value = value
	c.RemoveNote(value)
}

func (c *Cell) RemoveNote(value int) {
	c.notes[value] = false
}

func (c Cell) GetNotes() (result []int) {
	for k, v := range c.notes {
		if v {
			result = append(result, k)
		}
	}
	return result
}

func (c Cell) Dim() int  { return c.dim }
func (c Cell) Size() int { return c.dim * c.dim }

func (c Cell) String() (s string) {
	if c.value > 0 {
		s += fmt.Sprintf("%v", c.value)
	} else {
		for _, v := range c.GetNotes() {
			s += fmt.Sprintf("(%v)", v)
		}
	}
	return s
}
