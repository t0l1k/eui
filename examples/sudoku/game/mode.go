package game

import (
	"fmt"
)

type Dim struct {
	W, H int
}

func NewDim(w, h int) *Dim   { return &Dim{W: w, H: h} }
func (d Dim) Size() int      { return d.W * d.H }
func (d Dim) String() string { return fmt.Sprintf("Size:%v W:%v H:%v", d.Size(), d.W, d.H) }
