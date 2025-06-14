package eui

import (
	"sync"
	"sync/atomic"
)

type Slot func(data any)
type EqualFunc func(a, b any) bool

type Signal struct {
	mu        sync.RWMutex
	slots     map[int64]Slot
	lastVal   any
	equalFunc EqualFunc
	nextID    int64
}

func NewSignal(equal ...EqualFunc) *Signal {
	var eq EqualFunc
	if len(equal) > 0 {
		eq = equal[0]
	} else {
		eq = func(a, b any) bool { return a == b }
	}
	return &Signal{
		slots:     make(map[int64]Slot),
		equalFunc: eq,
	}
}

// Возвращает id для отписки
func (s *Signal) Connect(slot Slot) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := atomic.AddInt64(&s.nextID, 1)
	s.slots[id] = slot
	return id
}

// ConnectAndFire — сразу исполняет слот с текущим значением
func (s *Signal) ConnectAndFire(slot Slot, value any) int64 {
	id := s.Connect(slot)
	s.lastVal = value
	go slot(value)
	return id
}

func (s *Signal) Disconnect(id int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.slots, id)
}

func (s *Signal) Emit(value any) {
	s.mu.Lock()
	if s.equalFunc != nil && s.equalFunc(s.lastVal, value) {
		s.mu.Unlock()
		return
	}
	s.lastVal = value
	slots := make([]Slot, 0, len(s.slots))
	for _, slot := range s.slots {
		slots = append(slots, slot)
	}
	s.mu.Unlock()
	for _, slot := range slots {
		go slot(value)
	}
}

func (s *Signal) Value() any { s.mu.RLock(); defer s.mu.RUnlock(); return s.lastVal }
