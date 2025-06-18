package eui

import (
	"sync/atomic"
)

// Signal реализует потокобезопасный паттерн "сигнал-слот" (наблюдатель) с поддержкой generics.
//
// Signal[T] позволяет подписчикам (слотам) получать уведомления об изменении значения типа T.
// Подписчики добавляются через метод Connect и могут быть удалены методом Disconnect.
// Для предотвращения лишних уведомлений можно передать функцию сравнения EqualFunc[T].
//
// Если функция сравнения не указана, уведомления будут отправляться всегда при вызове Emit.
//
// Пример использования:
//
//	// Создание сигнала для int с проверкой на равенство
//	sig := eui.NewSignal[int](func(a, b int) bool { return a == b })
//
//	// Подписка на изменения
//	id := sig.Connect(func(val int) { fmt.Println("Новое значение:", val) })
//
//	// Изменение значения
//	sig.Emit(42)
//
//	// Отписка
//	sig.Disconnect(id)
//
// Для сложных типов (например, срезов или структур) рекомендуется явно передавать функцию сравнения.
//

type SlotFunc[T any] func(data T)
type EqualFunc[T any] func(a, b T) bool

type Signal[T any] struct {
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
	id := atomic.AddInt64(&s.nextID, 1)
	s.slots[id] = slot
	return id
}

// ConnectAndFire — сразу исполняет слот с текущим значением
func (s *Signal[T]) ConnectAndFire(slot SlotFunc[T], value T) int64 {
	id := s.Connect(slot)
	s.lastVal = value
	slot(value)
	return id
}

func (s *Signal[T]) Disconnect(id int64) {
	delete(s.slots, id)
}

func (s *Signal[T]) Emit(value T) {
	if s.equalFunc != nil && s.equalFunc(s.lastVal, value) {
		return
	}
	s.lastVal = value
	slots := make([]SlotFunc[T], 0, len(s.slots))
	for _, slot := range s.slots {
		slots = append(slots, slot)
	}
	for _, slot := range slots {
		slot(value)
	}
	// log.Println("Signal:Emit:", value)
}

func (s *Signal[T]) Value() T { return s.lastVal }
