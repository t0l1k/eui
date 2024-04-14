package game

import (
	"fmt"

	"github.com/t0l1k/eui"
)

type Note struct {
	dim    int
	values []int
}

func NewNote(dim int) *Note {
	n := &Note{dim: dim}
	n.Reset()
	return n
}

func (n *Note) Reset() {
	n.values = nil
	for i := 1; i <= n.dim*n.dim; i++ {
		n.values = append(n.values, i)
	}
}

func (n *Note) GetNoteValues() []int  { return n.values }
func (n Note) IsFound(value int) bool { return eui.IntSliceContains(n.values, value) }

func (n *Note) UpdateNote(value int) bool {
	if n.IsFound(value) {
		n.values = eui.RemoveFromIntSliceValue(n.values, value)
		return true
	}
	return false
}

func (n Note) String() (result string) {
	if len(n.values) > 0 {
		size := n.dim * n.dim
		for i := 1; i <= size; i++ {
			if n.IsFound(i) {
				result += fmt.Sprintf("(%v)", i)
			}
		}
	}
	return result
}
