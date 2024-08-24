package utils

import (
	"fmt"
	"strconv"
)

type IntList []int

func NewIntList() IntList { return make(IntList, 0) }

func (l IntList) Add(value int) IntList { return append(l, value) }

func (l IntList) Remove(value int) (IntList, error) {
	if !l.IsContain(value) {
		return nil, fmt.Errorf("элемент %v отстутствует", value)
	}
	idx, err := l.Index(value)
	if err != nil {
		return nil, fmt.Errorf("ошибка при удалении %v", value)
	}
	copy(l[idx:], l[idx+1:])
	l[len(l)-1] = value
	l = l[:len(l)-1]
	return l, nil
}

func (l IntList) IsContain(value int) bool {
	for _, v := range l {
		if value == v {
			return true
		}
	}
	return false
}

func (l IntList) Equals(other []int) bool {
	if len(l) != len(other) {
		return false
	}
	for i, v := range l {
		if other[i] != v {
			return false
		}
	}
	return true
}

func (l IntList) Index(value int) (int, error) {
	for i, v := range l {
		if value == v {
			return i, nil
		}
	}
	return -1, fmt.Errorf("элемент %v отстутствует", value)
}

func (l IntList) Pop() IntList {
	l[len(l)-1] = -1
	l = l[:len(l)-1]
	return l
}

func (l IntList) Max() (max int) {
	for _, v := range l {
		if v > max {
			max = v
		}
	}
	return max
}

func (l IntList) Min() (min int) {
	for _, v := range l {
		if v < min {
			min = v
		}
	}
	return min
}

func (l IntList) String() (result string) {
	for _, v := range l {
		result += fmt.Sprintf("(%v)", strconv.Itoa(v))
	}
	return result
}
