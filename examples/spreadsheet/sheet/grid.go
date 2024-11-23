package sheet

import (
	"fmt"
	"strconv"
	"unicode"
)

type Grid struct{ row, column int }

func NewGrid(row, column int) Grid  { return Grid{row, column} }
func (g Grid) Get() (int, int)      { return g.row, g.column }
func (g Grid) GetRow() string       { return string('A' + rune(g.row)) }
func (g Grid) GetColumn() string    { return strconv.Itoa(g.column) }
func (g *Grid) Set(row, column int) { g.row = row; g.column = column }
func (g Grid) String() string       { return fmt.Sprintf("[%v%v]", g.GetRow(), g.GetColumn()) }

func GridParse(value string) Grid {
	var (
		a, b   []rune
		aa, bb int
		res    string
	)
	for _, v := range value {
		if unicode.IsNumber(v) {
			b = append(b, v)
		} else {
			if unicode.IsLower(v) {
				v = unicode.ToUpper(v)
			}
			a = append(a, v)
		}
	}
	for _, v := range a {
		aa += int(v - 'A')
	}
	for _, v := range b {
		res += string(v)
	}
	bb, err := strconv.Atoi(res)
	if err != nil {
		panic(err)
	}
	return NewGrid(aa, bb)
}
