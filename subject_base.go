package eui

type SubjectBase struct {
	value    interface{}
	listener []Observer
}

func (s *SubjectBase) Attach(o Observer) {
	s.listener = append(s.listener, o)
}

func (s *SubjectBase) Detach(o Observer) {
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
