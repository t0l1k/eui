package eui

// Умею передать подписчикам об изменении переменной
type StringVar struct {
	value    string
	listener []Observer
}

func NewStringVar(value string) *StringVar {
	return &StringVar{
		value: value,
	}
}

func (s *StringVar) Attach(o Observer) {
	s.listener = append(s.listener, o)
}

func (s *StringVar) Detach(o Observer) {
	for i, observer := range s.listener {
		if observer == o {
			s.listener = append(s.listener[:i], s.listener[i+1:]...)
			break
		}
	}
}

func (s *StringVar) Notify() {
	for _, observer := range s.listener {
		observer.UpdateData(s.value)
	}
}

func (s *StringVar) Get() string {
	return s.value
}

func (s *StringVar) Set(value string) {
	if s.value == value {
		return
	}
	s.value = value
	s.Notify()
}
