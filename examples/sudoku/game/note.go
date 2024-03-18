package game

import "fmt"

type Note struct {
	dim            int
	values, excess []int
}

func NewNote(dim int) *Note {
	n := &Note{dim: dim}
	for i := 1; i <= n.dim*n.dim; i++ {
		n.values = append(n.values, i)
	}
	return n
}

func (n Note) IsFound(value int) bool {
	found, _ := n.isContain(n.values, value)
	return found
}

func (n Note) IsExcess(value int) bool {
	found, _ := n.isContain(n.excess, value)
	return found
}

func (n *Note) AddNote(value int) {
	found, _ := n.isContain(n.values, value)
	if found {
		n.excess = append(n.excess, value)
	} else {
		n.values = append(n.values, value)
	}
}

func (n *Note) RemoveNote(value int) {
	n.values = n.removeFrom(n.values, value)
	n.excess = n.removeFrom(n.excess, value)
}

func (n Note) String() (result string) {
	if len(n.values) > 0 {
		size := n.dim * n.dim
		for i := 1; i <= size; i++ {
			if n.IsFound(i) && !n.IsExcess(i) {
				result += fmt.Sprintf("%v", i)
			} else if n.IsExcess(i) {
				result += fmt.Sprintf("%v!", i)
			} else {
				result += "."
			}
			if i%n.dim == 0 {
				result += "\n"
			}
		}
	}
	return result
}

func (n Note) removeFrom(arr []int, value int) []int {
	found, idx := n.isContain(arr, value)
	if !found {
		return arr
	}
	copy(arr[idx:], arr[idx+1:])
	arr[len(arr)-1] = 0
	arr = arr[:len(arr)-1]
	return arr
}

func (n Note) isContain(arr []int, value int) (bool, int) {
	for i, v := range arr {
		if value == v {
			return true, i
		}
	}
	return false, -1
}
