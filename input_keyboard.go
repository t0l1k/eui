package eui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type KeyboardData struct{ keys []ebiten.Key }

func (k *KeyboardData) GetKeys() []ebiten.Key { return k.keys }

// Умею передать подписчикам события от клавиатуры, пока только цифры. При нажатой клавише, более 250 мс символ дублируется.
type KeyboardInput struct {
	value    KeyboardData
	timer    *Timer
	listener []Inputer
}

// Пауза 250мс до следующего нажатия
func NewKeyboardInput() *KeyboardInput { return &KeyboardInput{timer: NewTimer(250)} }

func (s *KeyboardInput) Attach(o Inputer) { s.listener = append(s.listener, o) }

func (s *KeyboardInput) Detach(o Inputer) {
	for i, observer := range s.listener {
		if observer == o {
			s.listener = append(s.listener[:i], s.listener[i+1:]...)
			break
		}
	}
}

func (s *KeyboardInput) SetValue(keys []ebiten.Key) {
	s.value.keys = nil
	s.value.keys = append(s.value.keys, keys...)
	s.Notify()
}

func (s *KeyboardInput) Notify() {
	for _, observer := range s.listener {
		observer.UpdateInput(s.value)
	}
}

// Передать новое или повторное нажатие после истечения паузы
func (s *KeyboardInput) update(dt int) {
	keys := inpututil.AppendPressedKeys(s.value.keys[:0])
	if len(keys) == 0 {
		s.timer.Off()
		return
	}
	s.timer.Update(dt)
	if len(keys) > 0 {
		if !s.timer.IsOn() {
			s.SetValue(keys)
			s.timer.On()
		}
		if s.timer.IsDone() {
			s.SetValue(keys)
			s.timer.Reset()
		}
	}
}
