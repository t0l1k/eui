package game

import (
	"fmt"

	"github.com/t0l1k/eui"
)

type Cell struct {
	eui.SubjectBase
	notes    *Note
	readOnly bool
}

func NewCell(dim int) *Cell {
	c := &Cell{}
	c.SetValue(0)
	c.notes = NewNote(dim)
	return c
}

func (c *Cell) Reset() {
	c.readOnly = false
	c.SetValue(0)
	c.notes.Reset()
}

func (c *Cell) Add(value int) bool {
	if c.readOnly {
		return false
	}
	c.SetValue(value)
	c.notes.RemoveNote(value)
	return true
}

func (c *Cell) IsReadOnly() bool { return c.readOnly }
func (c *Cell) MarkReadOnly() {
	c.readOnly = true
}

func (c *Cell) RemoveNote(value int) {
	v := c.Value().(int)
	c.SetValue(v)
	c.notes.RemoveNote(value)
}

func (c *Cell) GetDim() int                     { return c.notes.dim }
func (c *Cell) GetValue() int                   { return c.Value().(int) }
func (c *Cell) GetNotes() ([]int, []int, []int) { return c.notes.GetNoteValues() }

func (c *Cell) AddNote(value int) {
	c.notes.AddNote(value)
}

func (c Cell) String() (result string) {
	if c.Value().(int) > 0 {
		result = fmt.Sprintf("[%5v]", c.Value())
	} else {
		result = fmt.Sprintf("[%5v]", c.notes.String())
	}
	return result
}
