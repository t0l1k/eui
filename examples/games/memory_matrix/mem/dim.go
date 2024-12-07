package mem

import "fmt"

type Dim struct {
	w, h int
}

func NewDim(w, h int) Dim { return Dim{w: w, h: h} }

func (d Dim) Size() int              { return d.w * d.h }
func (d Dim) Width() int             { return d.w }
func (d Dim) Height() int            { return d.h }
func (d Dim) Idx(x, y int) int       { return y*d.w + x }
func (d Dim) Pos(idx int) (int, int) { return idx % d.w, idx / d.h }
func (d Dim) String() string         { return fmt.Sprintf("dim:%vx%v", d.w, d.h) }
