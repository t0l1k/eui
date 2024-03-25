package game

import (
	"fmt"

	"github.com/t0l1k/eui"
)

type Note struct {
	dim                  int
	values, excess, user []int
}

func NewNote(dim int) *Note {
	n := &Note{dim: dim}
	for i := 1; i <= n.dim*n.dim; i++ {
		n.values = append(n.values, i)
	}
	return n
}

func (n *Note) Reset() {
	n.excess = nil
	n.values = nil
	n.user = nil
	for i := 1; i <= n.dim*n.dim; i++ {
		n.values = append(n.values, i)
	}
}

func (n *Note) GetNoteValues() ([]int, []int, []int) {
	return n.values, n.excess, n.user
}

func (n Note) IsFound(value int) bool {
	return eui.IntSliceContains(n.values, value)
}

func (n Note) IsExcess(value int) bool {
	return eui.IntSliceContains(n.excess, value)
}

func (n *Note) AddNote(value int) {
	if value == 0 {
		panic("Note found 0")
	}
	found := eui.IntSliceContains(n.values, value)
	if found {
		n.excess = append(n.excess, value)
	} else {
		n.values = append(n.values, value)
	}
}

func (n *Note) RemoveNote(value int) {
	if value == 0 {
		panic("Note found 0")
	}
	n.values = eui.RemoveFromIntSliceValue(n.values, value)
	n.excess = eui.RemoveFromIntSliceValue(n.excess, value)
}

func (n Note) String() (result string) {
	if len(n.values) > 0 {
		size := n.dim * n.dim
		for i := 1; i <= size; i++ {
			if n.IsFound(i) && !n.IsExcess(i) {
				result += fmt.Sprintf("(%v)", i)
			} else if n.IsExcess(i) {
				result += fmt.Sprintf("(%v!)", i)
			}
		}
	}
	return result
}
