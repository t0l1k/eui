package sheet

import (
	"strings"
)

type Sheet map[Grid]*Cell

func NewSheet() *Sheet { s := make(Sheet, 0); return &s }

func (s Sheet) InitCell(index Grid, value *Cell) { s[index] = value }

func (s Sheet) AddValue(index Grid, value string) {
	cell := s.Cell(index)
	if !cell.IsContainValue() {
		isFormula, newValue := s.IsFormula(cell, value)
		if isFormula {
			s[index].Emit(newValue)
		} else {
			s[index].Emit(value)
		}
	}
}

func (s Sheet) Cell(index Grid) *Cell { return s[index] }

func (s *Sheet) IsFormula(cell *Cell, value string) (bool, string) {
	if strings.HasPrefix(value, "=") {
		formula := NewFormulaAdd(s, cell, value)
		sum := formula.Calc()
		return true, sum
	}
	return false, ""
}
