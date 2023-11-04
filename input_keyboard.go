package eui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type KeyboardData struct {
	keys []ebiten.Key
}

// Умею передать подписчикам события от клавиатуры, пока только цифры. При нажатой клавише, более 250 мс символ дублируется.
type KeyboardInput struct {
	value    KeyboardData
	timer    *Timer
	listener []Input
}

// Пауза 250мс до следующего нажатия
func NewKeyboardInput() *KeyboardInput { return &KeyboardInput{timer: NewTimer(250)} }

func (s *KeyboardInput) Attach(o Input) { s.listener = append(s.listener, o) }

func (s *KeyboardInput) Detach(o Input) {
	for i, observer := range s.listener {
		if observer == o {
			s.listener = append(s.listener[:i], s.listener[i+1:]...)
			break
		}
	}
}

func (s *KeyboardInput) Notify() {
	for _, observer := range s.listener {
		observer.UpdateInput(s.value)
	}
}

// Передать новое или повторное нажатие после истечения паузы
func (s *KeyboardInput) update(dt int) {
	s.value.keys = inpututil.AppendPressedKeys(s.value.keys[:0])
	if len(s.value.keys) > 0 && !s.timer.IsOn() {
		s.timer.On()
		s.Notify()
	}
	if s.timer.IsOn() {
		s.timer.Update(dt)
	}
	if s.timer.IsDone() {
		s.Notify()
		s.timer.Reset()
	}
	if len(s.value.keys) == 0 {
		s.timer.Off()
	}
}
