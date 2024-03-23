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

func (n *Note) Reset() {
	n.excess = nil
	n.values = nil
	for i := 1; i <= n.dim*n.dim; i++ {
		n.values = append(n.values, i)
	}
}

func (n Note) IsFound(value int) bool {
	return ArrIsContain(n.values, value)
}

func (n Note) IsExcess(value int) bool {
	return ArrIsContain(n.excess, value)
}

func (n *Note) AddNote(value int) {
	if value == 0 {
		panic("Note found 0")
	}
	found := ArrIsContain(n.values, value)
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
	n.values = RemoveFromArr(n.values, value)
	n.excess = RemoveFromArr(n.excess, value)
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

func RemoveFromArr(arr []int, value int) []int {
	if !ArrIsContain(arr, value) {
		return arr
	}
	idx := ArrIsContainGetIdx(arr, value)
	copy(arr[idx:], arr[idx+1:])
	arr[len(arr)-1] = value
	arr = arr[:len(arr)-1]
	return arr
}

func ArrIsContain(arr []int, value int) bool {
	for _, v := range arr {
		if value == v {
			return true
		}
	}
	return false
}

func ArrIsContainGetIdx(arr []int, value int) int {
	for i, v := range arr {
		if value == v {
			return i
		}
	}
	return -1
}
