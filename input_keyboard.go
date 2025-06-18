package eui

import (
	"log"
	"slices"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type KeyboardData struct {
	keysP, keysR []ebiten.Key
	runes        []rune
}

func NewKeyboardData(p, r []ebiten.Key, rn []rune) KeyboardData {
	return KeyboardData{keysP: p, keysR: r, runes: rn}
}

func (k *KeyboardData) GetKeysPressed() []ebiten.Key     { return k.keysP }
func (k *KeyboardData) GetKeysReleased() []ebiten.Key    { return k.keysR }
func (k *KeyboardData) GetRunes() []rune                 { return k.runes }
func (k *KeyboardData) IsPressed(value ebiten.Key) bool  { return slices.Contains(k.keysP, value) }
func (k *KeyboardData) IsReleased(value ebiten.Key) bool { return slices.Contains(k.keysR, value) }

// Умею передать подписчикам события от клавиатуры. При нажатой клавише, более 250 мс символ дублируется.
type KeyboardInput struct {
	*Signal[Event]
	timer *Timer
}

// Пауза 250мс до следующего нажатия
func NewKeyboardInput(fn SlotFunc[Event]) *KeyboardInput {
	k := &KeyboardInput{
		Signal: NewSignal[Event](),
		timer:  NewTimer(250*time.Millisecond, nil),
	}
	k.Connect(fn)
	return k
}

func (s *KeyboardInput) SetKeypressDelay(value time.Duration) { s.timer.SetDuration(value) }

// Передать новое или повторное нажатие клавиши после истечения паузы, для символов([]rune) повтор не работает
func (s *KeyboardInput) update(int) {
	keysP := inpututil.AppendPressedKeys(nil)
	keysR := inpututil.AppendJustReleasedKeys(nil)
	runes := ebiten.AppendInputChars(nil)
	kd := NewKeyboardData(keysP, keysR, runes)
	if len(keysP) == 0 && len(keysR) == 0 {
		s.timer.Off()
		return
	}
	if len(keysP) > 0 {
		if !s.timer.IsOn() {
			s.Emit(NewEvent(EventKeyPressed, kd))
			s.timer.On()
			log.Println("KeyboardInput:Emit:pressed", kd, s.timer.String())
		}
		if s.timer.IsDone() {
			s.Emit(NewEvent(EventKeyPressed, kd))
			s.timer.Reset()
			log.Println("KeyboardInput:Timer:repeat", kd, s.timer.String())
		}
	}
	if len(keysR) > 0 {
		s.Emit(NewEvent(EventKeyReleased, kd))
		log.Println("KeyboardInput:Emit:released", kd, s.timer.String())
	}
}
