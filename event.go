package eui

import (
	"fmt"
)

type EventType int

const (
	EventResize EventType = iota
	EventTick
	EventKeyPressed
	EventKeyReleased
	EventMouseDown
	EventMouseUp
	EventMouseDrag
	EventMouseMovement
	EventMouseWheel
	EventMouseEnter
	EventMouseLeave
)

func (e EventType) String() string {
	return [...]string{
		"Resize",
		"Tick",
		"KeyPressed",
		"KeyReleased",
		"MouseDown",
		"MouseUp",
		"MouseDrag",
		"MouseMovement",
		"MouseWheel",
		"MouseEnter",
		"MouseLeave",
	}[e]
}

type Event struct {
	Type  EventType
	Value any
}

func NewEvent(t EventType, v any) Event {
	e := Event{Type: t, Value: v}
	if !(e.Type == EventTick || e.Type == EventMouseMovement) {
		// log.Println("New:", e.String())
	}
	return e
}

func (e Event) String() string { return fmt.Sprintf("Event:%v %v", e.Type.String(), e.Value) }
