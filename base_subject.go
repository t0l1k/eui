package eui

import "fmt"

type SubjectBase struct {
	value    interface{}
	listener []Observerer
}

func NewSubject() *SubjectBase { return &SubjectBase{} }

func (s *SubjectBase) Attach(o Observerer) {
	s.listener = append(s.listener, o)
}

func (s *SubjectBase) Detach(o Observerer) {
	for i, observer := range s.listener {
		if observer == o {
			s.listener = append(s.listener[:i], s.listener[i+1:]...)
			break
		}
	}
}

func (s *SubjectBase) Notify() {
	for _, observer := range s.listener {
		observer.UpdateData(s.value)
	}
}

func (s *SubjectBase) Value() interface{} {
	return s.value
}

func (s *SubjectBase) SetValue(value interface{}) {
	s.value = value
	s.Notify()
}

func (s *SubjectBase) String() string {
	return fmt.Sprintf("%v", s.value)
}
