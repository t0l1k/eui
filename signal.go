package eui

import (
	"sync"
	"sync/atomic"
)

type SlotFunc[T any] func(data T)
type EqualFunc[T any] func(a, b T) bool

type Signal[T any] struct {
	mu        sync.RWMutex
	slots     map[int64]SlotFunc[T]
	lastVal   T
	equalFunc EqualFunc[T]
	nextID    int64
}

func NewSignal[T any](equal ...EqualFunc[T]) *Signal[T] {
	var eq EqualFunc[T]
	if len(equal) > 0 {
		eq = equal[0]
	} else {
		eq = nil
	}
	return &Signal[T]{
		slots:     make(map[int64]SlotFunc[T]),
		equalFunc: eq,
	}
}

// Возвращает id для отписки
func (s *Signal[T]) Connect(slot SlotFunc[T]) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := atomic.AddInt64(&s.nextID, 1)
	s.slots[id] = slot
	return id
}

// ConnectAndFire — сразу исполняет слот с текущим значением
func (s *Signal[T]) ConnectAndFire(slot SlotFunc[T], value T) int64 {
	id := s.Connect(slot)
	s.lastVal = value
	go slot(value)
	return id
}

func (s *Signal[T]) Disconnect(id int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.slots, id)
}

func (s *Signal[T]) Emit(value T) {
	s.mu.Lock()
	if s.equalFunc != nil && s.equalFunc(s.lastVal, value) {
		s.mu.Unlock()
		return
	}
	s.lastVal = value
	slots := make([]SlotFunc[T], 0, len(s.slots))
	for _, slot := range s.slots {
		slots = append(slots, slot)
	}
	s.mu.Unlock()
	for _, slot := range slots {
		go slot(value)
	}
}

func (s *Signal[T]) Value() T { s.mu.RLock(); defer s.mu.RUnlock(); return s.lastVal }
