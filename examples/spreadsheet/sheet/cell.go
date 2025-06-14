package sheet

import (
	"fmt"
	"strconv"

	"github.com/t0l1k/eui"
)

type Cell struct {
	*eui.Signal
	grid    Grid
	active  bool
	formula Formula
}

func NewCell(grid Grid) *Cell {
	c := &Cell{Signal: eui.NewSignal()}
	c.Emit("")
	c.grid = grid
	c.active = false
	c.formula = nil
	return c
}

func (c *Cell) GetValue() string         { return c.Value().(string) }
func (c *Cell) IsContainValue() bool     { return !(c.Value() == "") }
func (c *Cell) GetFormula() Formula      { return c.formula }
func (c *Cell) IsContainFormula() bool   { return c.formula != nil }
func (c *Cell) SetFormula(value Formula) { c.formula = value }
func (c *Cell) RemoveFormula()           { c.formula = nil; c.Emit("") }

func (c *Cell) GetNum() int {
	if c.IsContainValue() {
		v, err := strconv.Atoi(c.Value().(string))
		if err != nil {
			fmt.Println("error cell get num", v, c.Value(), err)
			panic(err)
		}
		return v
	}
	return 0
}

func (c *Cell) IsActive() bool { return c.active }
func (c *Cell) SetActive()     { c.active = true }
func (c *Cell) SetInActive()   { c.active = false }

// func (c *Cell) UpdateData(value any) {
// 	switch value.(type) {
// 	case string:
// 		c.formula.Calc()
// 	}
// }

func (c *Cell) String() string { return fmt.Sprintf("%v", c.Value()) }
