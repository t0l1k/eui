package sheet

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/t0l1k/eui"
)

type Formula interface {
	Calc() string
	String() string
}

type FormulaAdd struct {
	eui.SubjectBase
	sheet   *Sheet
	srcCell *Cell
	formula string
}

func NewFormulaAdd(sheet *Sheet, srcCell *Cell, formula string) *FormulaAdd {
	f := &FormulaAdd{}
	f.sheet = sheet
	f.srcCell = srcCell
	f.formula = formula

	f.attach()
	f.srcCell.SetFormula(f)
	return f
}

func (f *FormulaAdd) attach() {
	for _, v := range f.parseFormula(f.formula) {
		grid := GridParse(v)
		cell := f.sheet.Cell(grid)
		cell.Attach(f.sheet.Cell(f.srcCell.grid))
	}
}

func (f *FormulaAdd) getValues() (result []string) {
	for _, v := range f.parseFormula(f.formula) {
		grid := GridParse(v)
		cell := f.sheet.Cell(grid)
		result = append(result, cell.GetValue())
	}
	return result
}

func (f *FormulaAdd) parseFormula(value string) (result []string) {
	result = strings.Split(value, "=")
	for _, v := range result {
		result = strings.Split(v, "+")
	}
	return result
}

func (f *FormulaAdd) parseValues(values []string) string {
	var arr []float64
	for _, v := range values {
		num, err := strconv.ParseFloat(v, 64)
		if err != nil {
			panic(err)
		}
		arr = append(arr, num)
	}

	sum := 0.0
	for _, v := range arr {
		sum += v
	}
	x := fmt.Sprintf("%.0f", sum)
	f.srcCell.SetValue(x)
	return x
}

func (f *FormulaAdd) Calc() string {
	return f.parseValues(f.getValues())
}

func (f *FormulaAdd) String() string {
	return f.formula
}
