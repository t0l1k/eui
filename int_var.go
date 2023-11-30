package eui

import "strconv"

// Умею передать подписчикам об изменении переменной.
type IntVar struct {
	SubjectBase
}

func NewIntVar(p int) *IntVar {
	i := &IntVar{}
	i.SetValue(p)
	return i
}

func (s *IntVar) String() string {
	return strconv.Itoa(s.value.(int))
}
