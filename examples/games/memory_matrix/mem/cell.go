package mem

import "github.com/t0l1k/eui"

const (
	CellEmpty  = " "
	CellFilled = "*"
)

type Cell struct {
	eui.SubjectBase
	readonly, marked bool
}

func NewCell() *Cell {
	c := &Cell{}
	c.SetValue(CellEmpty)
	return c
}

func (c *Cell) IsEmpty() bool    { return c.Value().(string) == CellEmpty }
func (c *Cell) IsReadOnly() bool { return c.readonly }
func (c *Cell) SetReadOnly()     { c.readonly = true }
func (c *Cell) IsMarked() bool   { return c.marked }
func (c *Cell) SetMarked()       { c.marked = true }
func (c *Cell) Move()            { c.SetValue(CellFilled) }
func (c *Cell) String() string {
	if c.IsEmpty() {
		return CellEmpty
	}
	return CellFilled
}
