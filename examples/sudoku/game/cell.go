package game

import (
	"fmt"

	"github.com/t0l1k/eui"
)

type Cell struct {
	eui.SubjectBase
	notes           *Note
	readOnly, wrong bool
}

func NewCell(dim int) *Cell {
	c := &Cell{}
	c.SetValue(0)
	c.notes = NewNote(dim)
	return c
}

func (c *Cell) Reset() {
	c.readOnly = false
	c.wrong = false
	c.SetValue(0)
	c.notes.Reset()
}

func (c *Cell) Add(value int) bool {
	if c.readOnly {
		return false
	}
	if !c.notes.UpdateNote(value) {
		return false
	}
	c.SetValue(value)
	return true
}

func (c *Cell) UpdateNote(value int) {
	v := c.Value().(int)
	c.SetValue(v)
	c.notes.UpdateNote(value)
}

func (c *Cell) IsReadOnly() bool { return c.readOnly }
func (c *Cell) SetReadOnly()     { c.readOnly = true }
func (c *Cell) GetDim() int      { return c.notes.dim }
func (c *Cell) GetValue() int    { return c.Value().(int) }
func (c *Cell) GetNotes() []int  { return c.notes.GetNoteValues() }

func (c Cell) String() (result string) {
	if c.Value().(int) > 0 {
		result = fmt.Sprintf("[%3v]", c.Value())
	} else {
		result = fmt.Sprintf("[%3v]", c.notes.String())
	}
	return result
}
