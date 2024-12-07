package mem

import (
	"testing"
)

func TestLevelDim(t *testing.T) {
	var tt = map[int]Dim{
		1:  NewDim(2, 2),
		2:  NewDim(3, 2),
		3:  NewDim(3, 3),
		4:  NewDim(4, 3),
		13: NewDim(8, 8),
		15: NewDim(9, 9),
		16: NewDim(10, 9),
	}
	t.Run("Test level dim", func(t *testing.T) {
		for k, v := range tt {
			got := GetDimForLevel(Level(k))
			want := v
			if got != want {
				t.Errorf("got:%v want:%v", got, want)
			}

		}
	})
}
