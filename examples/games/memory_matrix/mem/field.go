package mem

import "math/rand"

type field []*Cell

func newField() *field {
	f := make(field, 0)
	return &f
}

func (f field) Shuffle(level int, dim Dim) *field {
	for y := 0; y < dim.h; y++ {
		for x := 0; x < dim.w; x++ {
			f = append(f, NewCell())
		}
	}
	mined := 0
	for mined < level {
		x, y := rand.Intn(dim.w), rand.Intn(dim.h)
		if f.cell(dim.Idx(x, y)).IsEmpty() {
			f.cell(dim.Idx(x, y)).Move()
			mined++
		}
	}
	for _, v := range f {
		if v.IsEmpty() {
			continue
		}
		v.SetReadOnly()
	}
	return &f
}

func (f field) cell(idx int) *Cell { return f[idx] }
