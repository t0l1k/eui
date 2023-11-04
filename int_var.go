package eui

import "strconv"

// Умею передать подписчикам об изменении переменной.
type IntVar struct {
	value    int
	listener []Observer
}

func NewIntVar(p int) *IntVar {
	return &IntVar{
		value: p,
	}
}

func (s *IntVar) Attach(o Observer) {
	s.listener = append(s.listener, o)
}

func (s *IntVar) Detach(o Observer) {
	for i, observer := range s.listener {
		if observer == o {
			s.listener = append(s.listener[:i], s.listener[i+1:]...)
			break
		}
	}
}

func (s *IntVar) Notify() {
	for _, observer := range s.listener {
		observer.UpdateData(s.value)
	}
}

func (s *IntVar) Get() int {
	return s.value
}

func (s *IntVar) Set(value int) {
	s.value = value
	s.Notify()
}

func (s *IntVar) String() string {
	return strconv.Itoa(s.value)
}
