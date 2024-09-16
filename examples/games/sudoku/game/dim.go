package game

import (
	"fmt"
)

type Dim struct {
	W, H int
}

func NewDim(w, h int) Dim       { return Dim{W: w, H: h} }
func (d Dim) Size() int         { return d.W * d.H }
func (d Dim) Eq(other Dim) bool { return d.W == other.W && d.H == other.H }
func (d Dim) String() string    { return fmt.Sprintf("Dim:%v(W:%v H:%v)", d.Size(), d.W, d.H) }

type DimsBySize []Dim

func (ds DimsBySize) Len() int           { return len(ds) }
func (ds DimsBySize) Less(i, j int) bool { return ds[i].Size() < ds[j].Size() }
func (ds DimsBySize) Swap(i, j int)      { ds[i], ds[j] = ds[j], ds[i] }
