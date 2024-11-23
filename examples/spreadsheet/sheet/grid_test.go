package sheet

import (
	"testing"
)

func TestGridParse(t *testing.T) {
	strParse := map[string][]string{
		"a1":  {"a1", "[A1]"},
		"A1":  {"A1", "[A1]"},
		"b1":  {"b1", "[B1]"},
		"B1":  {"B1", "[B1]"},
		"b10": {"b10", "[B10]"},
		"B10": {"B10", "[B10]"},
		"Z10": {"Z10", "[Z10]"},
	}

	for k, v := range strParse {
		got := GridParse(v[0]).String()
		want := v[1]
		if got != want {
			t.Errorf("t:%v got %v want %v", k, got, want)
		}
	}
}
