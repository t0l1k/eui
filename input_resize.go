package eui

import "log"

type ResizeListener struct{ *Signal[Event] }

func NewResizeListener(fn SlotFunc[Event]) *ResizeListener {
	r := &ResizeListener{Signal: NewSignal(func(a, b Event) bool {
		if a != b {
			return false
		}
		aR := a.Value.(Rect)
		bR := b.Value.(Rect)
		log.Println("ResizeListener:Equal", aR, bR)
		return aR.Equal(bR)
	})}
	r.Connect(fn)
	return r
}
