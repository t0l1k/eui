package eui

type ResizeListener struct{ *Signal[Event] }

func NewResizeListener(fn SlotFunc[Event]) *ResizeListener {
	r := &ResizeListener{Signal: NewSignal(func(a, b Event) bool { return false })}
	r.Connect(fn)
	return r
}
