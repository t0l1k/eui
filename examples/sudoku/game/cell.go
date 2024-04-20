package game

import (
	"fmt"

	"github.com/t0l1k/eui"
)

type Cell struct {
	eui.SubjectBase
	dim      int
	readOnly bool
	notes    map[int]bool
}

func NewCell(dim int) *Cell {
	c := &Cell{dim: dim}
	c.Reset()
	return c
}

func (c *Cell) Reset() {
	c.readOnly = false
	c.SetValue(0)
	c.notes = make(map[int]bool)
	for i := 1; i <= c.Size(); i++ {
		c.notes[i] = true
	}
}

func (c *Cell) Add(value int) bool {
	if c.readOnly || !c.notes[value] {
		return false
	}
	c.SetValue(value)
	c.UpdateNote(value)
	return true
}

func (c *Cell) Dim() int                  { return c.dim }
func (c *Cell) Size() int                 { return c.dim * c.dim }
func (c *Cell) IsReadOnly() bool          { return c.readOnly }
func (c *Cell) SetReadOnly()              { c.readOnly = true }
func (c *Cell) GetValue() int             { return c.Value().(int) }
func (c Cell) IsFoundNote(value int) bool { return c.notes[value] }
func (c *Cell) UpdateNote(value int)      { c.notes[value] = false }

func (c *Cell) GetNotes() (notes []int) {
	for k, v := range c.notes {
		if v {
			notes = append(notes, k)
		}
	}
	return notes
}

func (c Cell) String() (result string) {
	result = "["
	if c.Value().(int) > 0 {
		result = fmt.Sprintf("%3v", c.Value())
	} else {
		for k, v := range c.notes {
			if v {
				result += fmt.Sprintf("(%v)", k)
			}
		}
	}
	return result + "]"
}
