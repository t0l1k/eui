package eui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type buttonPressStatus int

const (
	buttonPressed  buttonPressStatus = 1
	buttonReleased buttonPressStatus = -1
)

type MouseData struct {
	position PointInt
	button   buttonPressStatus
}

// Умею передать подписчикам события от мыши это текущее местоположение на экране и нажатие кнопок левая средняя или правая
type MouseInput struct {
	value    MouseData
	listener []Inputer
}

func NewMouseInput(fn SlotFunc[Event]) *MouseInput { return &MouseInput{} }

func (s *MouseInput) Attach(o Inputer) {
	s.listener = append(s.listener, o)
}

func (s *MouseInput) Detach(o Inputer) {
	for i, observer := range s.listener {
		if observer == o {
			s.listener = append(s.listener[:i], s.listener[i+1:]...)
			break
		}
	}
}

func (s *MouseInput) Notify() {
	for _, observer := range s.listener {
		observer.UpdateInput(s.value)
	}
}

func (s *MouseInput) getPosition() (int, int) {
	return s.value.position.X, s.value.position.Y
}

func (s *MouseInput) getButton() buttonPressStatus {
	return s.value.button
}

func (s *MouseInput) set(value MouseData) {
	s.value = value
	s.Notify()
}

func (s *MouseInput) update(_ int) {
	x0, y0 := ebiten.CursorPosition()
	b0 := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) ||
		ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) ||
		ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle)
	x, y := s.getPosition()
	b := s.getButton()
	if !b0 && b != buttonReleased || !b0 && (x != x0 || y != y0) {
		// Изменилось местоположение курсора мыши кнопка отпущена
		s.set(MouseData{position: PointInt{x0, y0}, button: buttonReleased})
		// log.Println("Новые кординаты у мыши", x0, y0, b0, dt)
	} else if (b0 && b != buttonPressed) || (x != x0 || y != y0) {
		// Нажата кнопка мыши
		s.set(MouseData{position: PointInt{x0, y0}, button: buttonPressed})
		// log.Println("Нажата кнопка или новые кординаты", x0, y0, b0, dt)
	}
}
