package game

import (
	"fmt"
	"strconv"

	"github.com/t0l1k/eui"
)

type Cell struct {
	eui.SubjectBase
	notes    map[int]bool
	readOnly bool
}

func newCell() *Cell { return &Cell{} }

func (c *Cell) reset(size int) {
	c.SetValue(0)
	c.notes = make(map[int]bool)
	for i := 1; i <= size; i++ {
		c.notes[i] = true
	}
	c.readOnly = false
}

func (c Cell) GetValue() int { return c.Value().(int) }
func (c *Cell) add(value int) bool {
	if c.IsReadOnly() {
		return false
	}
	c.SetValue(value)
	c.setNote(value)
	return true
}
func (c Cell) IsReadOnly() bool   { return c.readOnly }
func (c *Cell) setReadOnly()      { c.readOnly = true }
func (c *Cell) setNote(value int) { c.notes[value] = false }
func (c Cell) GetNotes() (res []int) {
	for i, v := range c.notes {
		if v {
			res = append(res, i)
		}
	}
	return res
}

func (c Cell) String() (res string) {
	if c.GetValue() > 0 {
		res = strconv.Itoa(c.GetValue())
	} else {
		for k, v := range c.notes {
			if v {
				res += fmt.Sprintf("(%v)", k)
			}
		}
	}
	return res
}
