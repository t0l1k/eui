package sheet

import (
	"testing"
)

func TestInitAddUpdateCells(t *testing.T) {
	sheet := NewSheet()
	indexs := []string{"a1", "a2", "a3", "a4"}
	values := []string{"1", "2", "3", "=a1+a2+a3"}
	for _, v := range indexs { // init cells
		grid := GridParse(v)
		sheet.InitCell(grid, NewCell(grid))
	}

	for i, v := range indexs { // add values
		grid := GridParse(v)
		sheet.AddValue(grid, values[i])
	}

	t.Run("Test Formula Add", func(t *testing.T) {
		got := sheet.Cell(GridParse(indexs[3])).GetValue()
		want := "6"
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}

		sheet.Cell(GridParse(indexs[3])).SetValue(got)
		got = sheet.Cell(GridParse(indexs[3])).Value().(string)
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("Test New Value Updated And Update Formula", func(t *testing.T) {
		sheet.Cell(GridParse(indexs[0])).SetValue("10")
		got := sheet.Cell(GridParse(indexs[3])).GetValue()
		want := "15"
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}
