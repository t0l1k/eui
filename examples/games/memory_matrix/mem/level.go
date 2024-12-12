package mem

import "strconv"

type Level int

func (l Level) String() string { return strconv.Itoa(int(l)) }

func GetDimForLevel(level Level) Dim {
	w0, h0 := 2, 2
	for l := 1; l < int(level); l++ {
		if w0 == h0 {
			w0 += 1
		} else {
			h0 += 1
		}
	}
	return Dim{w0, h0}
}
