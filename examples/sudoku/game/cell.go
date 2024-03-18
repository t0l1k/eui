package game

import (
	"fmt"

	"github.com/t0l1k/eui"
)

type Cell struct {
	eui.SubjectBase
	notes *Note
	pos   eui.PointInt
}

func NewCell(dim int, pos eui.PointInt) *Cell {
	c := &Cell{pos: pos}
	c.SetValue(0)
	c.notes = NewNote(dim)
	return c
}

func (c *Cell) Add(value int) {
	c.SetValue(value)
}

func (c *Cell) AddNote(value int) {
	c.notes.AddNote(value)
}

func (c Cell) StringValueShort() (result string) {
	if c.Value().(int) > 0 {
		result = fmt.Sprintf("%3v", c.Value())
	} else {
		result = "..."
	}
	return result
}

func (c Cell) String() (result string) {
	if c.Value().(int) > 0 {
		result = fmt.Sprintf("\n%3v\n", c.Value().(int))
	} else {
		result = c.notes.String()
	}
	return result
}
