package sols

const (
	Col1 Column = iota
	Col2
	Col3
	Col4
	Col5
	Col6
	Col7
	Col8
	Col9
	Col10
	Col11
	Col12
	Col13
	Col14
	Col15
	ColFC1
	ColFc2
	ColFC3
	ColFC4
	ColH1
	ColH2
	ColH3
	ColH4
)

type Column int

func (from Column) IsEq(to Column) bool        { return from == to }
func (from Column) IsValidMove(to Column) bool { return !from.IsEq(to) }
func (f Column) String() (result string) {
	switch f {
	case ColFC1:
		result = "[FCfree cell 1]"
	case ColFc2:
		result = "[FCfree cell 2]"
	case ColFC3:
		result = "[FCfree cell 3]"
	case ColFC4:
		result = "[FCfree cell 4]"
	case ColH1:
		result = "[FChouse cell 1]"
	case ColH2:
		result = "[FChouse cell 2]"
	case ColH3:
		result = "[FChouse cell 3]"
	case ColH4:
		result = "[FChouse cell 4]"
	case Col1:
		result = "[FCcolumn 1 cell]"
	case Col2:
		result = "[FCcolumn 2 cell]"
	case Col3:
		result = "[FCcolumn 3 cell]"
	case Col4:
		result = "[FCcolumn 4 cell]"
	case Col5:
		result = "[FCcolumn 5 cell]"
	case Col6:
		result = "[FCcolumn 6 cell]"
	case Col7:
		result = "[FCcolumn 7 cell]"
	case Col8:
		result = "[FCcolumn 8 cell]"
	case Col9:
		result = "[FCcolumn 9 cell]"
	case Col10:
		result = "[FCcolumn 10 cell]"
	case Col11:
		result = "[FCcolumn 11 cell]"
	case Col12:
		result = "[FCcolumn 12 cell]"
	case Col13:
		result = "[FCcolumn 13 cell]"
	case Col14:
		result = "[FCcolumn 14 cell]"
	case Col15:
		result = "[FCcolumn 15 cell]"
	}
	return result
}
