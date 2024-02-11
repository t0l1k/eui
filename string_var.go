package eui

// Умею передать подписчикам об изменении переменной
type StringVar struct {
	SubjectBase
}

func NewStringVar(value string) *StringVar {
	s := &StringVar{}
	s.SetValue(value)
	return s
}
