package eui

import (
	"fmt"
	"log"
)

type EventType int

const (
	EventResize EventType = iota
)

func (e EventType) String() string {
	return [...]string{
		"Resize",
	}[e]
}

type Event struct {
	Type  EventType
	Value any
}

func NewEvent(t EventType, v any) Event {
	e := Event{Type: t, Value: v}
	log.Println("New:", e.String())
	return e
}

func (e Event) String() string { return fmt.Sprintf("Event:%v %v", e.Type.String(), e.Value) }

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
